// Copyright Istio Authors
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

package zipkin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	istioKube "istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/test/env"
	"istio.io/istio/pkg/test/framework/components/cluster"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/resource"
	testKube "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/scopes"
)

const (
	appName    = "zipkin"
	tracesAPI  = "/api/v2/traces?limit=%d&spanName=%s&annotationQuery=%s"
	zipkinPort = 9411

	remoteZipkinEntry = `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: tracing-gateway
  namespace: istio-system
spec:
  selector:
    istio: ingressgateway
  servers:
  - port:
      number: 80
      name: http-tracing
      protocol: HTTP
    hosts:
    - "tracing.{INGRESS_DOMAIN}"
  - port:
      number: 9411
      name: http-tracing-span
      protocol: HTTP
    hosts:
    - "tracing.{INGRESS_DOMAIN}"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: tracing-vs
  namespace: istio-system
spec:
  hosts:
  - "tracing.{INGRESS_DOMAIN}"
  gateways:
  - tracing-gateway
  http:
  - match:
    - port: 80
    route:
    - destination:
        host: tracing
        port:
          number: 80
  - match:
    - port: 9411
    route:
    - destination:
        host: tracing
        port:
          number: 9411
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: tracing
  namespace: istio-system
spec:
  host: tracing
  trafficPolicy:
    tls:
      mode: DISABLE
---`

	extServiceEntry = `
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  name: zipkin
spec:
  hosts:
  # must be of form name.namespace.global
  - zipkin.istio-system.global
  # Treat remote cluster services as part of the service mesh
  # as all clusters in the service mesh share the same root of trust.
  location: MESH_INTERNAL
  ports:
  - name: http-tracing-span
    number: 9411
    protocol: http
  resolution: DNS
  addresses:
  # the IP address to which zipkin.istio-system.global will resolve to
  # must be unique for each remote service, within a given cluster.
  # This address need not be routable. Traffic for this IP will be captured
  # by the sidecar and routed appropriately.
  - 240.0.0.2
  endpoints:
  # This is the routable address of the ingress gateway in cluster1 that
  # sits in front of zipkin service. Traffic from the sidecar will be
  # routed to this address.
  - address: {INGRESS_DOMAIN}
    ports:
      http-tracing-span: 15443 # Do not change this port value
`
)

var (
	_ Instance  = &kubeComponent{}
	_ io.Closer = &kubeComponent{}
)

type kubeComponent struct {
	id        resource.ID
	address   string
	forwarder istioKube.PortForwarder
	cluster   cluster.Cluster
}

func getZipkinYaml() (string, error) {
	yamlBytes, err := os.ReadFile(filepath.Join(env.IstioSrc, "samples/addons/extras/zipkin.yaml"))
	if err != nil {
		return "", err
	}
	yaml := string(yamlBytes)
	return yaml, nil
}

func installZipkin(ctx resource.Context, ns string) error {
	yaml, err := getZipkinYaml()
	if err != nil {
		return err
	}
	return ctx.ConfigKube().ApplyYAML(ns, yaml)
}

func installServiceEntry(ctx resource.Context, ns, ingressAddr string) error {
	// Setup remote access to zipkin in cluster
	yaml := strings.ReplaceAll(remoteZipkinEntry, "{INGRESS_DOMAIN}", ingressAddr)
	err := ctx.ConfigIstio().ApplyYAML(ns, yaml)
	if err != nil {
		return err
	}
	yaml = strings.ReplaceAll(extServiceEntry, "{INGRESS_DOMAIN}", ingressAddr)
	err = ctx.ConfigIstio().ApplyYAML(ns, yaml)
	if err != nil {
		return err
	}
	return nil
}

func newKube(ctx resource.Context, cfgIn Config) (Instance, error) {
	c := &kubeComponent{
		cluster: ctx.Clusters().GetOrDefault(cfgIn.Cluster),
	}
	c.id = ctx.TrackResource(c)

	// Find the zipkin pod and service, and start forwarding a local port.
	cfg, err := istio.DefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	if err := installZipkin(ctx, cfg.TelemetryNamespace); err != nil {
		return nil, err
	}

	fetchFn := testKube.NewSinglePodFetch(c.cluster, cfg.SystemNamespace, fmt.Sprintf("app=%s", appName))
	pods, err := testKube.WaitUntilPodsAreReady(fetchFn)
	if err != nil {
		return nil, err
	}
	pod := pods[0]

	forwarder, err := c.cluster.NewPortForwarder(pod.Name, pod.Namespace, "", 0, zipkinPort)
	if err != nil {
		return nil, err
	}

	if err := forwarder.Start(); err != nil {
		return nil, err
	}
	c.forwarder = forwarder
	scopes.Framework.Debugf("initialized zipkin port forwarder: %v", forwarder.Address())

	isIP := net.ParseIP(cfgIn.IngressAddr).String() != "<nil>"
	ingressDomain := cfgIn.IngressAddr
	if isIP {
		ingressDomain = fmt.Sprintf("%s.sslip.io", strings.ReplaceAll(cfgIn.IngressAddr, ":", "-"))
	}

	c.address = fmt.Sprintf("http://tracing.%s", ingressDomain)
	scopes.Framework.Debugf("Zipkin address: %s ", c.address)
	err = installServiceEntry(ctx, cfg.TelemetryNamespace, ingressDomain)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *kubeComponent) ID() resource.ID {
	return c.id
}

func (c *kubeComponent) QueryTraces(limit int, spanName, annotationQuery string) ([]Trace, error) {
	// Get 100 most recent traces
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	scopes.Framework.Debugf("make get call to zipkin api %v", c.address+fmt.Sprintf(tracesAPI, limit, spanName, annotationQuery))
	resp, err := client.Get(c.address + fmt.Sprintf(tracesAPI, limit, spanName, annotationQuery))
	if err != nil {
		scopes.Framework.Debugf("zipking err %v", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		scopes.Framework.Debugf("response err %v", resp.StatusCode)
		return nil, fmt.Errorf("zipkin api returns non-ok: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	traces, err := extractTraces(body)
	if err != nil {
		return nil, err
	}
	return traces, nil
}

// Close implements io.Closer.
func (c *kubeComponent) Close() error {
	if c.forwarder != nil {
		c.forwarder.Close()
	}
	return nil
}

func extractTraces(resp []byte) ([]Trace, error) {
	var traceObjs []interface{}
	if err := json.Unmarshal(resp, &traceObjs); err != nil {
		return []Trace{}, err
	}
	var ret []Trace
	for _, t := range traceObjs {
		spanObjs, ok := t.([]interface{})
		if !ok || len(spanObjs) == 0 {
			scopes.Framework.Debugf("cannot parse or cannot find spans in trace object %+v", t)
			continue
		}
		var spans []Span
		for _, obj := range spanObjs {
			newSpan := buildSpan(obj)
			spans = append(spans, newSpan)
		}
		for p := range spans {
			for c := range spans {
				if spans[c].ParentSpanID == spans[p].SpanID {
					spans[p].ChildSpans = append(spans[p].ChildSpans, &spans[c])
				}
			}
			// make order of child spans deterministic
			sort.Slice(spans[p].ChildSpans, func(i, j int) bool {
				return spans[p].ChildSpans[i].Name < spans[p].ChildSpans[j].Name
			})
		}
		ret = append(ret, Trace{Spans: spans})
	}
	if len(ret) > 0 {
		return ret, nil
	}
	return []Trace{}, errors.New("cannot find any traces")
}

func buildSpan(obj interface{}) Span {
	var s Span
	spanSpec := obj.(map[string]interface{})
	if spanID, ok := spanSpec["id"]; ok {
		s.SpanID = spanID.(string)
	}
	if parentSpanID, ok := spanSpec["parentId"]; ok {
		s.ParentSpanID = parentSpanID.(string)
	}
	if endpointObj, ok := spanSpec["localEndpoint"]; ok {
		if em, ok := endpointObj.(map[string]interface{}); ok {
			s.ServiceName = em["serviceName"].(string)
		}
	}
	if name, ok := spanSpec["name"]; ok {
		s.Name = name.(string)
	}
	return s
}
