//go:build integ
// +build integ

//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package util

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/hashicorp/go-multierror"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"istio.io/istio/pkg/config/protocol"
	"istio.io/istio/pkg/test"
	"istio.io/istio/pkg/test/echo/client"
	"istio.io/istio/pkg/test/framework"
	"istio.io/istio/pkg/test/framework/components/cluster"
	"istio.io/istio/pkg/test/framework/components/echo"
	"istio.io/istio/pkg/test/framework/components/echo/echoboot"
	"istio.io/istio/pkg/test/framework/components/echo/echotest"
	"istio.io/istio/pkg/test/framework/components/istio"
	"istio.io/istio/pkg/test/framework/components/istio/ingress"
	"istio.io/istio/pkg/test/framework/components/namespace"
	"istio.io/istio/pkg/test/framework/resource"
	"istio.io/istio/pkg/test/util/retry"
)

const (
	// The ID/name for the certificate chain in kubernetes tls secret.
	tlsScrtCert = "tls.crt"
	// The ID/name for the k8sKey in kubernetes tls secret.
	tlsScrtKey = "tls.key"
	// The ID/name for the CA certificate in kubernetes tls secret
	tlsScrtCaCert = "ca.crt"
	// The ID/name for the certificate chain in kubernetes generic secret.
	genericScrtCert = "cert"
	// The ID/name for the private key in kubernetes generic secret.
	genericScrtKey = "key"
	// The ID/name for the CA certificate in kubernetes generic secret.
	genericScrtCaCert = "cacert"
	ASvc              = "a"
	VMSvc             = "vm"
)

type EchoDeployments struct {
	ServerNs namespace.Instance
	All      echo.Instances
}

type IngressCredential struct {
	PrivateKey  string
	Certificate string
	CaCert      string
}

var IngressCredentialA = IngressCredential{
	PrivateKey:  TLSServerKeyA,
	Certificate: TLSServerCertA,
	CaCert:      CaCertA,
}

var IngressCredentialServerKeyCertA = IngressCredential{
	PrivateKey:  TLSServerKeyA,
	Certificate: TLSServerCertA,
}

var IngressCredentialCaCertA = IngressCredential{
	CaCert: CaCertA,
}

var IngressCredentialB = IngressCredential{
	PrivateKey:  TLSServerKeyB,
	Certificate: TLSServerCertB,
	CaCert:      CaCertB,
}

var IngressCredentialServerKeyCertB = IngressCredential{
	PrivateKey:  TLSServerKeyB,
	Certificate: TLSServerCertB,
}

// IngressKubeSecretYAML will generate a credential for a gateway
func IngressKubeSecretYAML(name, namespace string, ingressType CallType, ingressCred IngressCredential) string {
	// Create Kubernetes secret for ingress gateway
	secret := createSecret(ingressType, name, namespace, ingressCred, true)
	by, err := yaml.Marshal(secret)
	if err != nil {
		panic(err)
	}
	return string(by) + "\n---\n"
}

// CreateIngressKubeSecret reads credential names from credNames and key/cert from ingressCred,
// and creates K8s secrets for ingress gateway.
// nolint: interfacer
func CreateIngressKubeSecret(t framework.TestContext, credName string,
	ingressType CallType, ingressCred IngressCredential, isCompoundAndNotGeneric bool, clusters ...cluster.Cluster) {
	t.Helper()

	// Get namespace for ingress gateway pod.
	istioCfg := istio.DefaultConfigOrFail(t, t)
	systemNS := namespace.ClaimOrFail(t, t, istioCfg.SystemNamespace)
	CreateIngressKubeSecretInNamespace(t, credName, ingressType, ingressCred, isCompoundAndNotGeneric, systemNS.Name(), clusters...)
}

// CreateIngressKubeSecretInNamespace  reads credential names from credNames and key/cert from ingressCred,
// and creates K8s secrets for ingress gateway in the given namespace.
func CreateIngressKubeSecretInNamespace(t framework.TestContext, credName string,
	ingressType CallType, ingressCred IngressCredential, isCompoundAndNotGeneric bool, ns string, clusters ...cluster.Cluster) {
	t.Helper()

	t.ConditionalCleanup(func() {
		deleteKubeSecret(t, credName)
	})

	// Create Kubernetes secret for ingress gateway
	wg := multierror.Group{}
	if len(clusters) == 0 {
		clusters = t.Clusters()
	}
	for _, cluster := range clusters {
		cluster := cluster
		wg.Go(func() error {
			secret := createSecret(ingressType, credName, ns, ingressCred, isCompoundAndNotGeneric)
			_, err := cluster.CoreV1().Secrets(ns).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				if errors.IsAlreadyExists(err) {
					if _, err := cluster.CoreV1().Secrets(ns).Update(context.TODO(), secret, metav1.UpdateOptions{}); err != nil {
						return fmt.Errorf("failed to update secret (error: %s)", err)
					}
				} else {
					return fmt.Errorf("failed to update secret (error: %s)", err)
				}
			}
			// Check if Kubernetes secret is ready
			return retry.UntilSuccess(func() error {
				_, err := cluster.CoreV1().Secrets(ns).Get(context.TODO(), credName, metav1.GetOptions{})
				if err != nil {
					return fmt.Errorf("secret %v not found: %v", credName, err)
				}
				return nil
			}, retry.Timeout(time.Second*5))
		})
	}
	if err := wg.Wait().ErrorOrNil(); err != nil {
		t.Fatal(err)
	}
}

// deleteKubeSecret deletes a secret
// nolint: interfacer
func deleteKubeSecret(ctx framework.TestContext, credName string) {
	// Get namespace for ingress gateway pod.
	istioCfg := istio.DefaultConfigOrFail(ctx, ctx)
	systemNS := namespace.ClaimOrFail(ctx, ctx, istioCfg.SystemNamespace)

	// Create Kubernetes secret for ingress gateway
	cluster := ctx.Clusters().Default()
	var immediate int64
	err := cluster.CoreV1().Secrets(systemNS.Name()).Delete(context.TODO(), credName,
		metav1.DeleteOptions{GracePeriodSeconds: &immediate})
	if err != nil && !errors.IsNotFound(err) {
		ctx.Fatalf("Failed to delete secret (error: %s)", err)
	}
}

// createSecret creates a kubernetes secret which stores private key, server certificate for TLS ingress gateway.
// For mTLS ingress gateway, createSecret adds ca certificate into the secret object.

func createSecret(ingressType CallType, cn, ns string, ic IngressCredential, isCompoundAndNotGeneric bool) *v1.Secret {
	if ingressType == Mtls {
		if isCompoundAndNotGeneric {
			return &v1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      cn,
					Namespace: ns,
				},
				Data: map[string][]byte{
					tlsScrtCert:   []byte(ic.Certificate),
					tlsScrtKey:    []byte(ic.PrivateKey),
					tlsScrtCaCert: []byte(ic.CaCert),
				},
			}
		}
		return &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      cn,
				Namespace: ns,
			},
			Data: map[string][]byte{
				genericScrtCert:   []byte(ic.Certificate),
				genericScrtKey:    []byte(ic.PrivateKey),
				genericScrtCaCert: []byte(ic.CaCert),
			},
		}
	}
	data := map[string][]byte{}
	if ic.Certificate != "" {
		data[tlsScrtCert] = []byte(ic.Certificate)
	}
	if ic.PrivateKey != "" {
		data[tlsScrtKey] = []byte(ic.PrivateKey)
	}
	if ic.CaCert != "" {
		data[tlsScrtCaCert] = []byte(ic.CaCert)
	}
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cn,
			Namespace: ns,
		},
		Data: data,
	}
}

// CallType defines ingress gateway type
type CallType int

const (
	TLS CallType = iota
	Mtls
)

type ExpectedResponse struct {
	ResponseCode                 int
	SkipErrorMessageVerification bool
	ErrorMessage                 string
}

type TLSContext struct {
	// CaCert is inline base64 encoded root certificate that authenticates server certificate provided
	// by ingress gateway.
	CaCert string
	// PrivateKey is inline base64 encoded private key for test client.
	PrivateKey string
	// Cert is inline base64 encoded certificate for test client.
	Cert string
}

// SendRequestOrFail makes HTTPS request to ingress gateway to visit product page
func SendRequestOrFail(ctx framework.TestContext, ing ingress.Instance, host string, path string,
	callType CallType, tlsCtx TLSContext, exRsp ExpectedResponse) {
	doSendRequestsOrFail(ctx, ing, host, path, callType, tlsCtx, exRsp, false /* useHTTP3 */)
}

func SendQUICRequestsOrFail(ctx framework.TestContext, ing ingress.Instance, host string, path string,
	callType CallType, tlsCtx TLSContext, exRsp ExpectedResponse) {
	doSendRequestsOrFail(ctx, ing, host, path, callType, tlsCtx, exRsp, true /* useHTTP3 */)
}

func doSendRequestsOrFail(ctx framework.TestContext, ing ingress.Instance, host string, path string,
	callType CallType, tlsCtx TLSContext, exRsp ExpectedResponse, useHTTP3 bool) {
	ctx.Helper()
	opts := echo.CallOptions{
		Timeout: time.Second,
		Port: &echo.Port{
			Protocol: protocol.HTTPS,
		},
		Path: fmt.Sprintf("/%s", path),
		Headers: map[string][]string{
			"Host": {host},
		},
		HTTP3:  useHTTP3,
		CaCert: tlsCtx.CaCert,
		Validator: echo.And(
			echo.ValidatorFunc(
				func(resp client.ParsedResponses, err error) error {
					// Check that the error message is expected.
					if err != nil {
						// If expected error message is empty, but we got some error
						// message then it should be treated as error when error message
						// verification is not skipped. Error message verification is skipped
						// when the error message is non-deterministic.
						if !exRsp.SkipErrorMessageVerification && len(exRsp.ErrorMessage) == 0 {
							return fmt.Errorf("unexpected error: %w", err)
						}
						if !exRsp.SkipErrorMessageVerification && !strings.Contains(err.Error(), exRsp.ErrorMessage) {
							return fmt.Errorf("expected response error message %s but got %w",
								exRsp.ErrorMessage, err)
						}
						return nil
					}

					return resp.CheckCode(strconv.Itoa(exRsp.ResponseCode))
				})),
	}

	if callType == Mtls {
		opts.Key = tlsCtx.PrivateKey
		opts.Cert = tlsCtx.Cert
	}

	// Certs occasionally take quite a while to become active in Envoy, so retry for a long time (2min)
	ing.CallWithRetryOrFail(ctx, opts, retry.Timeout(time.Minute*2))
}

// RotateSecrets deletes kubernetes secrets by name in credNames and creates same secrets using key/cert
// from ingressCred.
func RotateSecrets(ctx framework.TestContext, credName string, // nolint:interfacer
	ingressType CallType, ingressCred IngressCredential, isCompoundAndNotGeneric bool) {
	ctx.Helper()
	cluster := ctx.Clusters().Default()
	ist := istio.GetOrFail(ctx, ctx)
	systemNS := namespace.ClaimOrFail(ctx, ctx, ist.Settings().SystemNamespace)
	scrt, err := cluster.CoreV1().Secrets(systemNS.Name()).Get(context.TODO(), credName, metav1.GetOptions{})
	if err != nil {
		ctx.Errorf("Failed to get secret %s:%s (error: %s)", systemNS.Name(), credName, err)
	}
	scrt = updateSecret(ingressType, scrt, ingressCred, isCompoundAndNotGeneric)
	if _, err = cluster.CoreV1().Secrets(systemNS.Name()).Update(context.TODO(), scrt, metav1.UpdateOptions{}); err != nil {
		ctx.Errorf("Failed to update secret %s:%s (error: %s)", scrt.Namespace, scrt.Name, err)
	}
	// Check if Kubernetes secret is ready
	retry.UntilSuccessOrFail(ctx, func() error {
		_, err := cluster.CoreV1().Secrets(systemNS.Name()).Get(context.TODO(), credName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("secret %v not found: %v", credName, err)
		}
		return nil
	}, retry.Timeout(time.Second*5))
}

// createSecret creates a kubernetes secret which stores private key, server certificate for TLS ingress gateway.
// For mTLS ingress gateway, createSecret adds ca certificate into the secret object.
func updateSecret(ingressType CallType, scrt *v1.Secret, ic IngressCredential, isCompoundAndNotGeneric bool) *v1.Secret {
	if ingressType == Mtls {
		if isCompoundAndNotGeneric {
			scrt.Data[tlsScrtCert] = []byte(ic.Certificate)
			scrt.Data[tlsScrtKey] = []byte(ic.PrivateKey)
			scrt.Data[tlsScrtCaCert] = []byte(ic.CaCert)
		} else {
			scrt.Data[genericScrtCert] = []byte(ic.Certificate)
			scrt.Data[genericScrtKey] = []byte(ic.PrivateKey)
			scrt.Data[genericScrtCaCert] = []byte(ic.CaCert)
		}
	} else {
		scrt.Data[tlsScrtCert] = []byte(ic.Certificate)
		scrt.Data[tlsScrtKey] = []byte(ic.PrivateKey)
	}
	return scrt
}

func EchoConfig(service string, ns namespace.Instance, buildVM bool) echo.Config {
	return echo.Config{
		Service:   service,
		Namespace: ns,
		Ports: []echo.Port{
			{
				Name:     "http",
				Protocol: protocol.HTTP,
				// We use a port > 1024 to not require root
				InstancePort: 8090,
			},
		},
		DeployAsVM: buildVM,
	}
}

func SetupTest(ctx resource.Context, apps *EchoDeployments) error {
	var err error
	apps.ServerNs, err = namespace.New(ctx, namespace.Config{
		Prefix: "ingress",
		Inject: true,
	})
	if err != nil {
		return err
	}
	buildVM := !ctx.Settings().SkipVM
	echos, err := echoboot.NewBuilder(ctx).
		WithClusters(ctx.Clusters()...).
		WithConfig(EchoConfig(ASvc, apps.ServerNs, false)).
		WithConfig(EchoConfig(VMSvc, apps.ServerNs, buildVM)).Build()
	if err != nil {
		return err
	}
	apps.All = echos
	return nil
}

type TestConfig struct {
	Mode           string
	CredentialName string
	Host           string
	ServiceName    string
}

const vsTemplate = `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: {{.CredentialName}}
spec:
  hosts:
  - "{{.Host}}"
  gateways:
  - {{.CredentialName}}
  http:
  - match:
    - uri:
        exact: /{{.CredentialName}}
    route:
    - destination:
        host: {{.ServiceName}}
        port:
          number: 80
`

const gwTemplate = `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: {{.CredentialName}}
spec:
  selector:
    istio: ingressgateway # use istio default ingress gateway
  servers:
  - port:
      number: 443
      name: https
      protocol: HTTPS
    tls:
      mode: {{.Mode}}
      credentialName: "{{.CredentialName}}"
    hosts:
    - "{{.Host}}"
`

func runTemplate(t test.Failer, tmpl string, params interface{}) string {
	tm, err := template.New("").Parse(tmpl)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	var buf bytes.Buffer
	if err := tm.Execute(&buf, params); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

func SetupConfig(ctx framework.TestContext, ns namespace.Instance, config ...TestConfig) func() {
	var apply []string
	for _, c := range config {
		apply = append(apply, runTemplate(ctx, vsTemplate, c), runTemplate(ctx, gwTemplate, c))
	}
	ctx.ConfigIstio().ApplyYAMLOrFail(ctx, ns.Name(), apply...)
	return func() {
		ctx.ConfigIstio().DeleteYAMLOrFail(ctx, ns.Name(), apply...)
	}
}

// RunTestMultiMtlsGateways deploys multiple mTLS gateways with SDS enabled, and creates kubernetes secret that stores
// private key, server certificate and CA certificate for each mTLS gateway. Verifies that all gateways are able to terminate
// mTLS connections successfully.
func RunTestMultiMtlsGateways(ctx framework.TestContext, inst istio.Instance, apps *EchoDeployments) { // nolint:interfacer
	var credNames []string
	var tests []TestConfig
	echotest.New(ctx, apps.All).
		SetupForDestination(func(ctx framework.TestContext, dst echo.Instances) error {
			for i := 1; i < 6; i++ {
				cred := fmt.Sprintf("runtestmultimtlsgateways-%d", i)
				tests = append(tests, TestConfig{
					Mode:           "MUTUAL",
					CredentialName: cred,
					Host:           fmt.Sprintf("runtestmultimtlsgateways%d.example.com", i),
					ServiceName:    dst[0].Config().Service,
				})
				credNames = append(credNames, cred)
			}
			SetupConfig(ctx, apps.ServerNs, tests...)
			return nil
		}).
		To(echotest.SingleSimplePodServiceAndAllSpecial()).
		RunFromClusters(func(ctx framework.TestContext, src cluster.Cluster, dest echo.Instances) {
			for _, cn := range credNames {
				CreateIngressKubeSecret(ctx, cn, Mtls, IngressCredentialA, false)
			}

			ing := inst.IngressFor(src)
			if ing == nil {
				ctx.Skip()
			}
			tlsContext := TLSContext{
				CaCert:     CaCertA,
				PrivateKey: TLSClientKeyA,
				Cert:       TLSClientCertA,
			}
			callType := Mtls

			for _, h := range tests {
				ctx.NewSubTest(h.Host).Run(func(t framework.TestContext) {
					SendRequestOrFail(t, ing, h.Host, h.CredentialName, callType, tlsContext,
						ExpectedResponse{ResponseCode: 200, ErrorMessage: ""})
				})
			}
		})
}

// RunTestMultiTLSGateways deploys multiple TLS gateways with SDS enabled, and creates kubernetes secret that stores
// private key and server certificate for each TLS gateway. Verifies that all gateways are able to terminate
// SSL connections successfully.
func RunTestMultiTLSGateways(ctx framework.TestContext, inst istio.Instance, apps *EchoDeployments) { // nolint:interfacer
	var credNames []string
	var tests []TestConfig
	echotest.New(ctx, apps.All).
		SetupForDestination(func(ctx framework.TestContext, dst echo.Instances) error {
			for i := 1; i < 6; i++ {
				cred := fmt.Sprintf("runtestmultitlsgateways-%d", i)
				tests = append(tests, TestConfig{
					Mode:           "SIMPLE",
					CredentialName: cred,
					Host:           fmt.Sprintf("runtestmultitlsgateways%d.example.com", i),
					ServiceName:    dst[0].Config().Service,
				})
				credNames = append(credNames, cred)
			}
			SetupConfig(ctx, apps.ServerNs, tests...)
			return nil
		}).
		To(echotest.SingleSimplePodServiceAndAllSpecial()).
		RunFromClusters(func(ctx framework.TestContext, src cluster.Cluster, dest echo.Instances) {
			for _, cn := range credNames {
				CreateIngressKubeSecret(ctx, cn, TLS, IngressCredentialA, false)
			}

			ing := inst.IngressFor(src)
			if ing == nil {
				ctx.Skip()
			}
			tlsContext := TLSContext{
				CaCert: CaCertA,
			}
			callType := TLS

			for _, h := range tests {
				ctx.NewSubTest(h.Host).Run(func(t framework.TestContext) {
					SendRequestOrFail(ctx, ing, h.Host, h.CredentialName, callType, tlsContext,
						ExpectedResponse{ResponseCode: 200, ErrorMessage: ""})
				})
			}
		})
}

// RunTestMultiQUICGateways deploys multiple TLS/mTLS gateways with SDS enabled, and creates kubernetes secret that stores
// private key and server certificate for each TLS/mTLS gateway. Verifies that all gateways are able to terminate
// QUIC connections successfully.
func RunTestMultiQUICGateways(ctx framework.TestContext, inst istio.Instance, callType CallType, apps *EchoDeployments) {
	var credNames []string
	var tests []TestConfig
	echotest.New(ctx, apps.All).
		SetupForDestination(func(ctx framework.TestContext, dst echo.Instances) error {
			for i := 1; i < 6; i++ {
				cred := fmt.Sprintf("runtestmultitlsgateways-%d", i)
				mode := "SIMPLE"
				if callType == Mtls {
					mode = "MUTUAL"
				}
				tests = append(tests, TestConfig{
					Mode:           mode,
					CredentialName: cred,
					Host:           fmt.Sprintf("runtestmultitlsgateways%d.example.com", i),
					ServiceName:    dst[0].Config().Service,
				})
				credNames = append(credNames, cred)
			}
			SetupConfig(ctx, apps.ServerNs, tests...)
			return nil
		}).
		To(echotest.SingleSimplePodServiceAndAllSpecial()).
		RunFromClusters(func(ctx framework.TestContext, src cluster.Cluster, dest echo.Instances) {
			for _, cn := range credNames {
				CreateIngressKubeSecret(ctx, cn, TLS, IngressCredentialA, false)
			}

			ing := inst.IngressFor(src)
			if ing == nil {
				ctx.Skip()
			}
			tlsContext := TLSContext{
				CaCert: CaCertA,
			}
			if callType == Mtls {
				tlsContext = TLSContext{
					CaCert:     CaCertA,
					PrivateKey: TLSClientKeyA,
					Cert:       TLSClientCertA,
				}
			}

			for _, h := range tests {
				ctx.NewSubTest(h.Host).Run(func(t framework.TestContext) {
					SendQUICRequestsOrFail(ctx, ing, h.Host, h.CredentialName, callType, tlsContext,
						ExpectedResponse{ResponseCode: 200, ErrorMessage: ""})
				})
			}
		})
}
