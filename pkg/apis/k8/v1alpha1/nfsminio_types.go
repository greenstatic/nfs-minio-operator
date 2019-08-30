package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NFSMinioSpec defines the desired state of NFSMinio
// +k8s:openapi-gen=true
type NFSMinioSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Username string `json:"username"`
	NFS NFSMinioSpecNFS `json:"nfs"`
	Domain string `json:"domain"`
}

type NFSMinioSpecNFS struct {
	Server string `json:"server"`
	Path string `json:"path"`
	ReadOnly bool `json:"readOnly"`
}

// NFSMinioStatus defines the observed state of NFSMinio
// +k8s:openapi-gen=true
type NFSMinioStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	SecretKeyHash []byte `json:"secretKeyHash"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NFSMinio is the Schema for the nfsminios API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type NFSMinio struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NFSMinioSpec   `json:"spec,omitempty"`
	Status NFSMinioStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NFSMinioList contains a list of NFSMinio
type NFSMinioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NFSMinio `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NFSMinio{}, &NFSMinioList{})
}
