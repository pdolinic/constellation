
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NodeImageSpec defines the desired state of NodeImage
type NodeImageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of NodeImage. Edit nodeimage_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// NodeImageStatus defines the observed state of NodeImage
type NodeImageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// NodeImage is the Schema for the nodeimages API
type NodeImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeImageSpec   `json:"spec,omitempty"`
	Status NodeImageStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeImageList contains a list of NodeImage
type NodeImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeImage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeImage{}, &NodeImageList{})
}