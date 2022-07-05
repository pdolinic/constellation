package k8sapi

import (
	"path/filepath"

	"github.com/edgelesssys/constellation/bootstrapper/internal/kubernetes/k8sapi/resources"
	"github.com/edgelesssys/constellation/internal/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeletconf "k8s.io/kubelet/config/v1beta1"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
)

// Uses types defined here: https://kubernetes.io/docs/reference/config-api/kubeadm-config.v1beta3/
// Slimmed down to the fields we require

const (
	bindPort        = 6443
	auditLogDir     = "/var/log/kubernetes/audit/"
	auditLogFile    = "audit.log"
	auditPolicyPath = "/etc/kubernetes/audit-policy.yaml"
)

type CoreOSConfiguration struct{}

func (c *CoreOSConfiguration) InitConfiguration(externalCloudProvider bool) KubeadmInitYAML {
	var cloudProvider string
	if externalCloudProvider {
		cloudProvider = "external"
	}
	return KubeadmInitYAML{
		InitConfiguration: kubeadm.InitConfiguration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: kubeadm.SchemeGroupVersion.String(),
				Kind:       "InitConfiguration",
			},
			NodeRegistration: kubeadm.NodeRegistrationOptions{
				CRISocket: "/run/containerd/containerd.sock",
				KubeletExtraArgs: map[string]string{
					"cloud-provider": cloudProvider,
					"network-plugin": "cni",
				},
			},
			// AdvertiseAddress will be overwritten later
			LocalAPIEndpoint: kubeadm.APIEndpoint{
				BindPort: bindPort,
			},
		},
		ClusterConfiguration: kubeadm.ClusterConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ClusterConfiguration",
				APIVersion: kubeadm.SchemeGroupVersion.String(),
			},
			KubernetesVersion: constants.KubernetesVersion,
			// necessary to be able to access the kubeapi server through localhost
			APIServer: kubeadm.APIServer{
				ControlPlaneComponent: kubeadm.ControlPlaneComponent{
					ExtraArgs: map[string]string{
						"audit-policy-file":   auditPolicyPath,
						"audit-log-path":      filepath.Join(auditLogDir, auditLogFile), // CIS benchmark
						"audit-log-maxage":    "30",                                     // CIS benchmark - Default value of Rancher
						"audit-log-maxbackup": "10",                                     // CIS benchmark - Default value of Rancher
						"audit-log-maxsize":   "100",                                    // CIS benchmark - Default value of Rancher
						"profiling":           "false",                                  // CIS benchmark
						"tls-cipher-suites": "TLS_AES_128_GCM_SHA256,TLS_AES_256_GCM_SHA384,TLS_CHACHA20_POLY1305_SHA256," +
							"TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256," +
							"TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384," +
							"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256," +
							"TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256," +
							"TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305," +
							"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,TLS_RSA_WITH_3DES_EDE_CBC_SHA,TLS_RSA_WITH_AES_128_CBC_SHA," +
							"TLS_RSA_WITH_AES_128_GCM_SHA256,TLS_RSA_WITH_AES_256_CBC_SHA,TLS_RSA_WITH_AES_256_GCM_SHA384", // CIS benchmark
					},
					ExtraVolumes: []kubeadm.HostPathMount{
						{
							Name:      "audit-log",
							HostPath:  auditLogDir,
							MountPath: auditLogDir,
							ReadOnly:  false,
							PathType:  corev1.HostPathDirectoryOrCreate,
						},
						{
							Name:      "audit",
							HostPath:  auditPolicyPath,
							MountPath: auditPolicyPath,
							ReadOnly:  true,
							PathType:  corev1.HostPathFile,
						},
					},
				},
				CertSANs: []string{"127.0.0.1", "10.118.0.1"},
			},
			ControllerManager: kubeadm.ControlPlaneComponent{
				ExtraArgs: map[string]string{
					"flex-volume-plugin-dir":      "/opt/libexec/kubernetes/kubelet-plugins/volume/exec/",
					"cloud-provider":              cloudProvider,
					"configure-cloud-routes":      "false",
					"profiling":                   "false", // CIS benchmark
					"terminated-pod-gc-threshold": "1000",  // CIS benchmark - Default value of Rancher
				},
			},
			Scheduler: kubeadm.ControlPlaneComponent{
				ExtraArgs: map[string]string{
					"profiling": "false",
				},
			},
		},
		// warning: this config is applied to every node in the cluster!
		KubeletConfiguration: kubeletconf.KubeletConfiguration{
			ProtectKernelDefaults: true, // CIS benchmark
			TLSCipherSuites: []string{
				"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
				"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
				"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
				"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
				"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
				"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
				"TLS_RSA_WITH_AES_256_GCM_SHA384",
				"TLS_RSA_WITH_AES_128_GCM_SHA256",
			}, // CIS benchmark
			TypeMeta: metav1.TypeMeta{
				APIVersion: kubeletconf.SchemeGroupVersion.String(),
				Kind:       "KubeletConfiguration",
			},
			RegisterWithTaints: []corev1.Taint{
				{
					Key:    "node.cloudprovider.kubernetes.io/uninitialized",
					Value:  "true",
					Effect: corev1.TaintEffectPreferNoSchedule,
				},
				{
					Key:    "node.cilium.io/agent-not-ready",
					Value:  "true",
					Effect: corev1.TaintEffectPreferNoSchedule,
				},
			},
		},
	}
}

func (c *CoreOSConfiguration) JoinConfiguration(externalCloudProvider bool) KubeadmJoinYAML {
	var cloudProvider string
	if externalCloudProvider {
		cloudProvider = "external"
	}
	return KubeadmJoinYAML{
		JoinConfiguration: kubeadm.JoinConfiguration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: kubeadm.SchemeGroupVersion.String(),
				Kind:       "JoinConfiguration",
			},
			NodeRegistration: kubeadm.NodeRegistrationOptions{
				CRISocket: "/run/containerd/containerd.sock",
				KubeletExtraArgs: map[string]string{
					"cloud-provider": cloudProvider,
				},
			},
			Discovery: kubeadm.Discovery{
				BootstrapToken: &kubeadm.BootstrapTokenDiscovery{},
			},
		},
		KubeletConfiguration: kubeletconf.KubeletConfiguration{
			TypeMeta: metav1.TypeMeta{
				APIVersion: kubeletconf.SchemeGroupVersion.String(),
				Kind:       "KubeletConfiguration",
			},
		},
	}
}

type KubeadmJoinYAML struct {
	JoinConfiguration    kubeadm.JoinConfiguration
	KubeletConfiguration kubeletconf.KubeletConfiguration
}

func (k *KubeadmJoinYAML) SetNodeName(nodeName string) {
	k.JoinConfiguration.NodeRegistration.Name = nodeName
}

func (k *KubeadmJoinYAML) SetApiServerEndpoint(apiServerEndpoint string) {
	k.JoinConfiguration.Discovery.BootstrapToken.APIServerEndpoint = apiServerEndpoint
}

func (k *KubeadmJoinYAML) SetToken(token string) {
	k.JoinConfiguration.Discovery.BootstrapToken.Token = token
}

func (k *KubeadmJoinYAML) AppendDiscoveryTokenCaCertHash(discoveryTokenCaCertHash string) {
	k.JoinConfiguration.Discovery.BootstrapToken.CACertHashes = append(k.JoinConfiguration.Discovery.BootstrapToken.CACertHashes, discoveryTokenCaCertHash)
}

func (k *KubeadmJoinYAML) SetNodeIP(nodeIP string) {
	if k.JoinConfiguration.NodeRegistration.KubeletExtraArgs == nil {
		k.JoinConfiguration.NodeRegistration.KubeletExtraArgs = map[string]string{"node-ip": nodeIP}
	} else {
		k.JoinConfiguration.NodeRegistration.KubeletExtraArgs["node-ip"] = nodeIP
	}
}

func (k *KubeadmJoinYAML) SetProviderID(providerID string) {
	k.KubeletConfiguration.ProviderID = providerID
}

func (k *KubeadmJoinYAML) SetControlPlane(advertiseAddress string, certificateKey string) {
	k.JoinConfiguration.ControlPlane = &kubeadm.JoinControlPlane{
		LocalAPIEndpoint: kubeadm.APIEndpoint{
			AdvertiseAddress: advertiseAddress,
			BindPort:         6443,
		},
		CertificateKey: certificateKey,
	}
}

func (k *KubeadmJoinYAML) Marshal() ([]byte, error) {
	return resources.MarshalK8SResources(k)
}

func (k *KubeadmJoinYAML) Unmarshal(yamlData []byte) (KubeadmJoinYAML, error) {
	var tmp KubeadmJoinYAML
	return tmp, resources.UnmarshalK8SResources(yamlData, &tmp)
}

type KubeadmInitYAML struct {
	InitConfiguration    kubeadm.InitConfiguration
	ClusterConfiguration kubeadm.ClusterConfiguration
	KubeletConfiguration kubeletconf.KubeletConfiguration
}

func (k *KubeadmInitYAML) SetNodeName(nodeName string) {
	k.InitConfiguration.NodeRegistration.Name = nodeName
}

// SetCertSANs sets the SANs for the certificate.
func (k *KubeadmInitYAML) SetCertSANs(certSANs []string) {
	for _, certSAN := range certSANs {
		if certSAN == "" {
			continue
		}
		k.ClusterConfiguration.APIServer.CertSANs = append(k.ClusterConfiguration.APIServer.CertSANs, certSAN)
	}
}

func (k *KubeadmInitYAML) SetApiServerAdvertiseAddress(apiServerAdvertiseAddress string) {
	k.InitConfiguration.LocalAPIEndpoint.AdvertiseAddress = apiServerAdvertiseAddress
}

// SetControlPlaneEndpoint sets the control plane endpoint if controlPlaneEndpoint is not empty.
func (k *KubeadmInitYAML) SetControlPlaneEndpoint(controlPlaneEndpoint string) {
	if controlPlaneEndpoint != "" {
		k.ClusterConfiguration.ControlPlaneEndpoint = controlPlaneEndpoint
	}
}

func (k *KubeadmInitYAML) SetServiceCIDR(serviceCIDR string) {
	k.ClusterConfiguration.Networking.ServiceSubnet = serviceCIDR
}

func (k *KubeadmInitYAML) SetPodNetworkCIDR(podNetworkCIDR string) {
	k.ClusterConfiguration.Networking.PodSubnet = podNetworkCIDR
}

func (k *KubeadmInitYAML) SetServiceDNSDomain(serviceDNSDomain string) {
	k.ClusterConfiguration.Networking.DNSDomain = serviceDNSDomain
}

func (k *KubeadmInitYAML) SetNodeIP(nodeIP string) {
	if k.InitConfiguration.NodeRegistration.KubeletExtraArgs == nil {
		k.InitConfiguration.NodeRegistration.KubeletExtraArgs = map[string]string{"node-ip": nodeIP}
	} else {
		k.InitConfiguration.NodeRegistration.KubeletExtraArgs["node-ip"] = nodeIP
	}
}

func (k *KubeadmInitYAML) SetProviderID(providerID string) {
	if k.InitConfiguration.NodeRegistration.KubeletExtraArgs == nil {
		k.InitConfiguration.NodeRegistration.KubeletExtraArgs = map[string]string{"provider-id": providerID}
	} else {
		k.InitConfiguration.NodeRegistration.KubeletExtraArgs["provider-id"] = providerID
	}
}

func (k *KubeadmInitYAML) Marshal() ([]byte, error) {
	return resources.MarshalK8SResources(k)
}

func (k *KubeadmInitYAML) Unmarshal(yamlData []byte) (KubeadmInitYAML, error) {
	var tmp KubeadmInitYAML
	return tmp, resources.UnmarshalK8SResources(yamlData, &tmp)
}