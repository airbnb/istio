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

package constants

const (
	// UnspecifiedIP constant for empty IP address
	UnspecifiedIP = "0.0.0.0"

	// AuthCertsPath is the path location for mTLS certificates
	AuthCertsPath = "/etc/certs/"

	// CertChainFilename is mTLS chain file
	CertChainFilename = "cert-chain.pem"

	// DefaultServerCertChain is the default path to the mTLS chain file
	DefaultCertChain = AuthCertsPath + CertChainFilename

	// KeyFilename is mTLS private key
	KeyFilename = "key.pem"

	// DefaultServerKey is the default path to the mTLS private key file
	DefaultKey = AuthCertsPath + KeyFilename

	// RootCertFilename is mTLS root cert
	RootCertFilename = "root-cert.pem"

	// DefaultRootCert is the default path to the mTLS root cert file
	DefaultRootCert = AuthCertsPath + RootCertFilename

	// ConfigPathDir config directory for storing envoy json config files.
	ConfigPathDir = "./etc/istio/proxy"

	// IstioDataDir is the directory to store binary data such as envoy core dump, profile, and downloaded Wasm modules.
	IstioDataDir = "/var/lib/istio/data"

	// BinaryPathFilename envoy binary location
	BinaryPathFilename = "/usr/local/bin/envoy"

	// ServiceClusterName service cluster name used in xDS calls
	ServiceClusterName = "istio-proxy"

	// IstioIngressGatewayName is the internal gateway name assigned to ingress
	IstioIngressGatewayName = "istio-autogenerated-k8s-ingress"

	KubernetesGatewayName = "istio-autogenerated-k8s-gateway"

	// IstioIngressNamespace is the namespace where Istio ingress controller is deployed
	IstioIngressNamespace = "istio-system"

	// DefaultKubernetesDomain the default service domain suffix for Kubernetes, if not overridden in config.
	DefaultKubernetesDomain = "cluster.local"

	// IstioLabel indicates that a workload is part of a named Istio system component.
	IstioLabel = "istio"

	// IstioIngressLabelValue is value for IstioLabel that identifies an ingress workload.
	// TODO we should derive this from IngressClass
	IstioIngressLabelValue = "ingressgateway"

	// IstioSystemNamespace is the namespace where Istio's components are deployed
	IstioSystemNamespace = "istio-system"

	// DefaultAuthenticationPolicyName is the name of the cluster-scoped authentication policy. Only
	// policy with this name in the cluster-scoped will be considered.
	DefaultAuthenticationPolicyName = "default"

	// IstioMeshGateway is the built in gateway for all sidecars
	IstioMeshGateway = "mesh"

	// The data name in the ConfigMap of each namespace storing the root cert of non-Kube CA.
	CACertNamespaceConfigMapDataName = "root-cert.pem"

	// PodInfoLabelsPath is the filepath that pod labels will be stored
	// This is typically set by the downward API
	PodInfoLabelsPath = "./etc/istio/pod/labels"

	// PodInfoAnnotationsPath is the filepath that pod annotations will be stored
	// This is typically set by the downward API
	PodInfoAnnotationsPath = "./etc/istio/pod/annotations"

	// DefaultServiceAccountName is the default service account to use for remote cluster access.
	DefaultServiceAccountName = "istio-reader-service-account"

	// DefaultConfigServiceAccountName is the default service account to use for external Istiod config cluster access.
	DefaultConfigServiceAccountName = "istiod"

	// KubeSystemNamespace is the system namespace where we place kubernetes system components.
	KubeSystemNamespace string = "kube-system"

	// KubePublicNamespace is the namespace where we place kubernetes public info (ConfigMaps).
	KubePublicNamespace string = "kube-public"

	// KubeNodeLeaseNamespace is the namespace for the lease objects associated with each kubernetes node.
	KubeNodeLeaseNamespace string = "kube-node-lease"

	// LocalPathStorageNamespace is the namespace for dynamically provisioning persistent local storage with
	// Kubernetes. Typically used with the Kind cluster: https://github.com/rancher/local-path-provisioner
	LocalPathStorageNamespace string = "local-path-storage"

	TestVMLabel = "istio.io/test-vm"

	TestVMVersionLabel = "istio.io/test-vm-version"

	// Label to skip config comparison.
	AlwaysPushLabel = "internal.istio.io/always-push"

	// TrustworthyJWTPath is the default 3P token to authenticate with third party services
	TrustworthyJWTPath = "./var/run/secrets/tokens/istio-token"

	// CertProviderIstiod uses istiod self signed DNS certificates for the control plane
	CertProviderIstiod = "istiod"
	// CertProviderKubernetes uses the Kubernetes CSR API to generate a DNS certificate for the control plane
	CertProviderKubernetes = "kubernetes"
	// CertProviderCustom uses the custom root certificate mounted in a well known location for the control plane
	CertProviderCustom = "custom"
	// CertProviderNone does not create any certificates for the control plane. It is assumed that some external
	// load balancer, such as an Istio Gateway, is terminating the TLS.
	CertProviderNone = "none"
)
