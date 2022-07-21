package versions

const (
	// Constellation images.
	// These images are built in a way that they support all versions currently listed in VersionConfigs.
	JoinImage          = "ghcr.io/edgelesssys/constellation/join-service:v1.3.2-0.20220718102802-8c25a227"
	AccessManagerImage = "ghcr.io/edgelesssys/constellation/access-manager:v1.3.2-0.20220714151638-d295be31"
	KmsImage           = "ghcr.io/edgelesssys/constellation/kmsserver:v1.3.2-0.20220714151638-d295be31"
	VerificationImage  = "ghcr.io/edgelesssys/constellation/verification-service:v1.3.2-0.20220714151638-d295be31"
	GcpGuestImage      = "ghcr.io/edgelesssys/gcp-guest-agent:latest"
)

// versionConfigs holds download URLs for all required kubernetes components for every supported version.
var VersionConfigs map[string]KubernetesVersion = map[string]KubernetesVersion{
	"1.23": {
		CNIPluginsURL:     "https://github.com/containernetworking/plugins/releases/download/v1.1.1/cni-plugins-linux-amd64-v1.1.1.tgz",
		CrictlURL:         "https://github.com/kubernetes-sigs/cri-tools/releases/download/v1.24.1/crictl-v1.24.1-linux-amd64.tar.gz",
		KubeletServiceURL: "https://raw.githubusercontent.com/kubernetes/release/v0.13.0/cmd/kubepkg/templates/latest/deb/kubelet/lib/systemd/system/kubelet.service",
		KubeadmConfURL:    "https://raw.githubusercontent.com/kubernetes/release/v0.13.0/cmd/kubepkg/templates/latest/deb/kubeadm/10-kubeadm.conf",
		KubeletURL:        "https://storage.googleapis.com/kubernetes-release/release/v1.23.6/bin/linux/amd64/kubelet",
		KubeadmURL:        "https://storage.googleapis.com/kubernetes-release/release/v1.23.6/bin/linux/amd64/kubeadm",
		KubectlURL:        "https://storage.googleapis.com/kubernetes-release/release/v1.23.6/bin/linux/amd64/kubectl",
		// CloudControllerManagerImageAWS is the CCM image used on AWS.
		CloudControllerManagerImageAWS: "us.gcr.io/k8s-artifacts-prod/provider-aws/cloud-controller-manager:v1.22.0-alpha.0",
		// CloudControllerManagerImageGCP is the CCM image used on GCP.
		// TODO: use newer "cloud-provider-gcp" from https://github.com/kubernetes/cloud-provider-gcp when newer releases are available.
		CloudControllerManagerImageGCP: "ghcr.io/edgelesssys/cloud-provider-gcp:sha-2f6a5b07fc2d37f24f8ff725132f87584d627d8f",
		// CloudControllerManagerImageAzure is the CCM image used on Azure.
		CloudControllerManagerImageAzure: "mcr.microsoft.com/oss/kubernetes/azure-cloud-controller-manager:v1.23.11",
		// CloudNodeManagerImageAzure is the cloud-node-manager image used on Azure.
		CloudNodeManagerImageAzure: "mcr.microsoft.com/oss/kubernetes/azure-cloud-node-manager:v1.23.11",
		// External service image. Depends on k8s version.
		ClusterAutoscalerImage: "k8s.gcr.io/autoscaling/cluster-autoscaler:v1.23.0",
	},
}

// KubernetesVersion bundles download URLs to all version-releated binaries necessary for installing/deploying a particular Kubernetes version.
type KubernetesVersion struct {
	CNIPluginsURL                    string
	CrictlURL                        string
	KubeletServiceURL                string
	KubeadmConfURL                   string
	KubeletURL                       string
	KubeadmURL                       string
	KubectlURL                       string
	CloudControllerManagerImageAWS   string
	CloudControllerManagerImageGCP   string
	CloudControllerManagerImageAzure string
	CloudNodeManagerImageAzure       string
	ClusterAutoscalerImage           string
}

// IsSupportedK8sVersion checks if a given Kubernetes version is supported by Constellation.
func IsSupportedK8sVersion(version string) bool {
	_, ok := VersionConfigs[version]
	return ok
}
