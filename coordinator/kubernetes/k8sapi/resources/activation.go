package resources

import (
	"fmt"

	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/secrets"
	apps "k8s.io/api/apps/v1"
	k8s "k8s.io/api/core/v1"
	rbac "k8s.io/api/rbac/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const activationImage = "ghcr.io/edgelesssys/constellation/activation-service:latest"

type activationDaemonset struct {
	ClusterRole        rbac.ClusterRole
	ClusterRoleBinding rbac.ClusterRoleBinding
	ConfigMap          k8s.ConfigMap
	DaemonSet          apps.DaemonSet
	ServiceAccount     k8s.ServiceAccount
	Service            k8s.Service
}

// NewActivationDaemonset returns a daemonset for the activation service.
func NewActivationDaemonset(csp, measurementsJSON, idJSON string) *activationDaemonset {
	return &activationDaemonset{
		ClusterRole: rbac.ClusterRole{
			TypeMeta: meta.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "ClusterRole",
			},
			ObjectMeta: meta.ObjectMeta{
				Name: "activation-service",
				Labels: map[string]string{
					"k8s-app": "activation-service",
				},
			},
			Rules: []rbac.PolicyRule{
				{
					APIGroups: []string{""},
					Resources: []string{"secrets"},
					Verbs:     []string{"get", "list", "create"},
				},
			},
		},
		ClusterRoleBinding: rbac.ClusterRoleBinding{
			TypeMeta: meta.TypeMeta{
				APIVersion: "rbac.authorization.k8s.io/v1",
				Kind:       "ClusterRoleBinding",
			},
			ObjectMeta: meta.ObjectMeta{
				Name: "activation-service",
			},
			RoleRef: rbac.RoleRef{
				APIGroup: "rbac.authorization.k8s.io",
				Kind:     "ClusterRole",
				Name:     "activation-service",
			},
			Subjects: []rbac.Subject{
				{
					Kind:      "ServiceAccount",
					Name:      "activation-service",
					Namespace: "kube-system",
				},
			},
		},
		DaemonSet: apps.DaemonSet{
			TypeMeta: meta.TypeMeta{
				APIVersion: "apps/v1",
				Kind:       "DaemonSet",
			},
			ObjectMeta: meta.ObjectMeta{
				Name:      "activation-service",
				Namespace: "kube-system",
				Labels: map[string]string{
					"k8s-app":                       "activation-service",
					"component":                     "activation-service",
					"kubernetes.io/cluster-service": "true",
				},
			},
			Spec: apps.DaemonSetSpec{
				Selector: &meta.LabelSelector{
					MatchLabels: map[string]string{
						"k8s-app": "activation-service",
					},
				},
				Template: k8s.PodTemplateSpec{
					ObjectMeta: meta.ObjectMeta{
						Labels: map[string]string{
							"k8s-app": "activation-service",
						},
					},
					Spec: k8s.PodSpec{
						PriorityClassName:  "system-cluster-critical",
						ServiceAccountName: "activation-service",
						Tolerations: []k8s.Toleration{
							{
								Key:      "CriticalAddonsOnly",
								Operator: k8s.TolerationOpExists,
							},
							{
								Key:      "node-role.kubernetes.io/master",
								Operator: k8s.TolerationOpEqual,
								Value:    "true",
								Effect:   k8s.TaintEffectNoSchedule,
							},
							{
								Operator: k8s.TolerationOpExists,
								Effect:   k8s.TaintEffectNoExecute,
							},
							{
								Operator: k8s.TolerationOpExists,
								Effect:   k8s.TaintEffectNoSchedule,
							},
						},
						// Only run on control plane nodes
						NodeSelector: map[string]string{
							"node-role.kubernetes.io/master": "",
						},
						ImagePullSecrets: []k8s.LocalObjectReference{
							{
								Name: secrets.PullSecretName,
							},
						},
						Containers: []k8s.Container{
							{
								Name:  "activation-service",
								Image: activationImage,
								Ports: []k8s.ContainerPort{
									{
										ContainerPort: 9090,
										Name:          "tcp",
									},
								},
								SecurityContext: &k8s.SecurityContext{
									Privileged: func(b bool) *bool { return &b }(true),
								},
								Args: []string{
									fmt.Sprintf("--cloud-provider=%s", csp),
									fmt.Sprintf("--kms-endpoint=kms.kube-system:%d", constants.KMSPort),
									"--v=5",
								},
								VolumeMounts: []k8s.VolumeMount{
									{
										Name:      "config",
										ReadOnly:  true,
										MountPath: constants.ActivationBasePath,
									},
									{
										Name:      "kubeadm",
										ReadOnly:  true,
										MountPath: "/etc/kubernetes",
									},
								},
							},
						},
						Volumes: []k8s.Volume{
							{
								Name: "config",
								VolumeSource: k8s.VolumeSource{
									ConfigMap: &k8s.ConfigMapVolumeSource{
										LocalObjectReference: k8s.LocalObjectReference{
											Name: "activation-config",
										},
									},
								},
							},
							{
								Name: "kubeadm",
								VolumeSource: k8s.VolumeSource{
									HostPath: &k8s.HostPathVolumeSource{
										Path: "/etc/kubernetes",
									},
								},
							},
						},
					},
				},
			},
		},
		ServiceAccount: k8s.ServiceAccount{
			TypeMeta: meta.TypeMeta{
				APIVersion: "v1",
				Kind:       "ServiceAccount",
			},
			ObjectMeta: meta.ObjectMeta{
				Name:      "activation-service",
				Namespace: "kube-system",
			},
		},
		Service: k8s.Service{
			TypeMeta: meta.TypeMeta{
				APIVersion: "v1",
				Kind:       "Service",
			},
			ObjectMeta: meta.ObjectMeta{
				Name:      "activation-service",
				Namespace: "kube-system",
			},
			Spec: k8s.ServiceSpec{
				Type: k8s.ServiceTypeNodePort,
				Ports: []k8s.ServicePort{
					{
						Name:       "grpc",
						Protocol:   k8s.ProtocolTCP,
						Port:       constants.ActivationServicePort,
						TargetPort: intstr.IntOrString{IntVal: constants.ActivationServicePort},
						NodePort:   constants.ActivationServiceNodePort,
					},
				},
				Selector: map[string]string{
					"k8s-app": "activation-service",
				},
			},
		},
		ConfigMap: k8s.ConfigMap{
			TypeMeta: meta.TypeMeta{
				APIVersion: "v1",
				Kind:       "ConfigMap",
			},
			ObjectMeta: meta.ObjectMeta{
				Name:      "activation-config",
				Namespace: "kube-system",
			},
			Data: map[string]string{
				"measurements": measurementsJSON,
				"id":           idJSON,
			},
		},
	}
}

// Marshal the daemonset using the Kubernetes resource marshaller.
func (a *activationDaemonset) Marshal() ([]byte, error) {
	return MarshalK8SResources(a)
}