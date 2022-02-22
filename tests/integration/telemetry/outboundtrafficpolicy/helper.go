//go:build integ
// +build integ

// Copyright Istio Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package outboundtrafficpolicy

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"istio.io/istio/pkg/config/protocol"
	echoclient "istio.io/istio/pkg/test/echo"
	"istio.io/istio/pkg/test/echo/common"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/environment/kube"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/components/prometheus"
	"istio.io/istio/pkg/test/framework/resource"
	tmpl "istio.io/istio/pkg/test/util/tmpl"
	promtest "istio.io/istio/tests/integration/telemetry/stats/prometheus"
)

const (
	// This service entry exists to create conflicts on various ports
	// As defined below, the tcp-conflict and https-conflict ports are 9443 and 9091
	ServiceEntry = `
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: http
spec:
  hosts:
  - istio.io
  location: MESH_EXTERNAL
  ports:
  - name: http-for-https
    number: 9443
    protocol: HTTP
  - name: http-for-tcp
    number: 9091
    protocol: HTTP
  resolution: DNS
`
	SidecarScope = `
apiVersion: networking.istio.io/v1alpha3
kind: Sidecar
metadata:
  name: restrict-to-service-entry-namespace
spec:
  egress:
  - hosts:
    - "{{.ImportNamespace}}/*"
    - "istio-system/*"
  outboundTrafficPolicy:
    mode: "{{.TrafficPolicyMode}}"
`

	Gateway = `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: istio-egressgateway
spec:
  selector:
    istio: egressgateway
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "some-external-site.com"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: route-via-egressgateway
spec:
  hosts:
    - "some-external-site.com"
  gateways:
  - istio-egressgateway
  - mesh
  http:
    - match:
      - gateways:
        - mesh # from sidecars, route to egress gateway service
        port: 80
      route:
      - destination:
          host: istio-egressgateway.istio-system.svc.cluster.local
          port:
            number: 80
        weight: 100
    - match:
      - gateways:
        - istio-egressgateway
        port: 80
      route:
      - destination:
          host: some-external-site.com
      headers:
        request:
          add:
            handled-by-egress-gateway: "true"
---
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: ext-service-entry
spec:
  hosts:
  - "some-external-site.com"
  location: MESH_EXTERNAL
  endpoints:
  - address: destination.{{.AppNamespace}}.svc.cluster.local
    network: external
  ports:
  - number: 80
    name: http
  resolution: DNS
`
)

// TestCase represents what is being tested
type TestCase struct {
	Name     string
	PortName string
	HTTP2    bool
	Host     string
	Expected Expected
}

// Expected contains the metric and query to run against
// prometheus to validate that expected telemetry information was gathered;
// as well as the http response code
type Expected struct {
	Metric          string
	PromQueryFormat string
	ResponseCode    []string
	// Metadata includes headers and additional injected information such as Method, Proto, etc.
	// The test will validate the returned metadata includes all options specified here
	Metadata map[string]string
}

// TrafficPolicy is the mode of the outbound traffic policy to use
// when configuring the sidecar for the client
type TrafficPolicy string

const (
	AllowAny     TrafficPolicy = "ALLOW_ANY"
	RegistryOnly TrafficPolicy = "REGISTRY_ONLY"
)

// String implements fmt.Stringer
func (t TrafficPolicy) String() string {
	return string(t)
}

// We want to test "external" traffic. To do this without actually hitting an external endpoint,
// we can import only the service namespace, so the apps are not known
func createSidecarScope(t *testing.T, ctx resource.Context, tPolicy TrafficPolicy, appsNamespace namespace.Instance, serviceNamespace namespace.Instance) {
	b := tmpl.EvaluateOrFail(t, SidecarScope, map[string]string{"ImportNamespace": serviceNamespace.Name(), "TrafficPolicyMode": tPolicy.String()})
	if err := ctx.ConfigIstio().ApplyYAML(appsNamespace.Name(), b); err != nil {
		t.Errorf("failed to apply service entries: %v", err)
	}
}

func mustReadCert(t *testing.T, f string) string {
	b, err := os.ReadFile(path.Join(env.IstioSrc, "tests/testdata/certs", f))
	if err != nil {
		t.Fatalf("failed to read %v: %v", f, err)
	}
	return string(b)
}

// We want to test "external" traffic. To do this without actually hitting an external endpoint,
// we can import only the service namespace, so the apps are not known
func createGateway(t *testing.T, ctx resource.Context, appsNamespace namespace.Instance, serviceNamespace namespace.Instance) {
	b := tmpl.EvaluateOrFail(t, Gateway, map[string]string{"AppNamespace": appsNamespace.Name()})
	if err := ctx.ConfigIstio().ApplyYAML(serviceNamespace.Name(), b); err != nil {
		t.Fatalf("failed to apply gateway: %v. template: %v", err, b)
	}
}

// TODO support native environment for registry only/gateway. Blocked by #13177 because the listeners for native use static
// routes and this test relies on the dynamic routes sent through pilot to allow external traffic.

func RunExternalRequest(cases []*TestCase, prometheus prometheus.Instance, mode TrafficPolicy, t *testing.T) {
	// Testing of Blackhole and Passthrough clusters:
	// Setup of environment:
	// 1. client and destination are deployed to app-1-XXXX namespace
	// 2. client is restricted to talk to destination via Sidecar scope where outbound policy is set (ALLOW_ANY, REGISTRY_ONLY)
	//    and clients' egress can only be to service-2-XXXX/* and istio-system/*
	// 3. a namespace service-2-YYYY is created
	// 4. A gateway is put in service-2-YYYY where its host is set for some-external-site.com on port 80 and 443
	// 3. a VirtualService is also created in service-2-XXXX to:
	//    a) route requests for some-external-site.com to the istio-egressgateway
	//       * if the request on port 80, then it will add an http header `handled-by-egress-gateway`
	//    b) from the egressgateway it will forward the request to the destination pod deployed in the app-1-XXX
	//       namespace

	// Test cases:
	// 1. http case:
	//    client -------> Hits listener 0.0.0.0_80 cluster
	//    Metric is istio_requests_total i.e. HTTP
	//
	// 2. https case:
	//    client ----> Hits no listener -> 0.0.0.0_150001 -> ALLOW_ANY/REGISTRY_ONLY
	//    Metric is istio_tcp_connections_closed_total i.e. TCP
	//
	// 3. https conflict case:
	//    client ----> Hits listener 0.0.0.0_9443
	//    Metric is istio_tcp_connections_closed_total i.e. TCP
	//
	// 4. http_egress
	//    client ) ---HTTP request (Host: some-external-site.com----> Hits listener 0.0.0.0_80 ->
	//      VS Routing (add Egress Header) --> Egress Gateway --> destination
	//    Metric is istio_requests_total i.e. HTTP with destination as destination
	//
	// 5. TCP
	//    client ---TCP request at port 9090----> Matches no listener -> 0.0.0.0_150001 -> ALLOW_ANY/REGISTRY_ONLY
	//    Metric is istio_tcp_connections_closed_total i.e. TCP
	//
	// 5. TCP conflict
	//    client ---TCP request at port 9091 ----> Hits listener 0.0.0.0_9091 ->  ALLOW_ANY/REGISTRY_ONLY
	//    Metric is istio_tcp_connections_closed_total i.e. TCP
	//
	framework.
		NewTest(t).
		Run(func(ctx framework.TestContext) {
			client, dest := setupEcho(t, ctx, mode)

			for _, tc := range cases {
				t.Run(tc.Name, func(t *testing.T) {
					client.CallWithRetryOrFail(t, echo.CallOptions{
						Target:   dest,
						PortName: tc.PortName,
						Headers: map[string][]string{
							"Host": {tc.Host},
						},
						HTTP2: tc.HTTP2,
						Check: func(resp echoclient.Responses, err error) error {
							// the expected response from a blackhole test case will have err
							// set; use the length of the expected code to ignore this condition
							if err != nil && len(tc.Expected.ResponseCode) != 0 {
								return fmt.Errorf("request failed: %v", err)
							}
							codes := make([]string, 0, len(resp))
							for _, r := range resp {
								codes = append(codes, r.Code)
							}
							if !reflect.DeepEqual(codes, tc.Expected.ResponseCode) {
								return fmt.Errorf("got codes %q, expected %q", codes, tc.Expected.ResponseCode)
							}

							for _, r := range resp {
								for k, v := range tc.Expected.Metadata {
									if got := r.RawResponse[k]; got != v {
										return fmt.Errorf("expected metadata %v=%v, got %q", k, v, got)
									}
								}
							}
							return nil
						},
					})

					if tc.Expected.Metric != "" {
						promtest.ValidateMetric(t, ctx.Clusters().Default(), prometheus, tc.Expected.PromQueryFormat, tc.Expected.Metric, 1)
					}
				})
			}
		})
}

func setupEcho(t *testing.T, ctx resource.Context, mode TrafficPolicy) (echo.Instance, echo.Instance) {
	appsNamespace := namespace.NewOrFail(t, ctx, namespace.Config{
		Prefix: "app",
		Inject: true,
	})
	serviceNamespace := namespace.NewOrFail(t, ctx, namespace.Config{
		Prefix: "service",
		Inject: true,
	})

	// External traffic should work even if we have service entries on the same ports
	createSidecarScope(t, ctx, mode, appsNamespace, serviceNamespace)

	var client, dest echo.Instance
	echoboot.NewBuilder(ctx).
		With(&client, echo.Config{
			Service:   "client",
			Namespace: appsNamespace,
			Subsets:   []echo.SubsetConfig{{}},
		}).
		With(&dest, echo.Config{
			Service:   "destination",
			Namespace: appsNamespace,
			Subsets:   []echo.SubsetConfig{{Annotations: echo.NewAnnotations().SetBool(echo.SidecarInject, false)}},
			Ports: []echo.Port{
				{
					// Plain HTTP port, will match no listeners and fall through
					Name:         "http",
					Protocol:     protocol.HTTP,
					ServicePort:  80,
					InstancePort: 8080,
				},
				{
					// HTTPS port, will match no listeners and fall through
					Name:         "https",
					Protocol:     protocol.HTTPS,
					ServicePort:  443,
					InstancePort: 8443,
					TLS:          true,
				},
				{
					// HTTPS port, there will be an HTTP service defined on this port that will match
					Name:        "https-conflict",
					Protocol:    protocol.HTTPS,
					ServicePort: 9443,
					TLS:         true,
				},
				{
					// TCP port, will match no listeners and fall through
					Name:        "tcp",
					Protocol:    protocol.TCP,
					ServicePort: 9090,
				},
				{
					// TCP port, there will be an HTTP service defined on this port that will match
					Name:        "tcp-conflict",
					Protocol:    protocol.TCP,
					ServicePort: 9091,
				},
			},
			TLSSettings: &common.TLSSettings{
				// Echo has these test certs baked into the docker image
				ClientCert: mustReadCert(t, "cert.crt"),
				Key:        mustReadCert(t, "cert.key"),
			},
		}).BuildOrFail(t)

	if err := ctx.ConfigIstio().ApplyYAML(serviceNamespace.Name(), ServiceEntry); err != nil {
		t.Errorf("failed to apply service entries: %v", err)
	}

	if _, kube := ctx.Environment().(*kube.Environment); kube {
		createGateway(t, ctx, appsNamespace, serviceNamespace)
	}
	return client, dest
}
