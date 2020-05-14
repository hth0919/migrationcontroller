package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MigrationSpec defines the desired state of Migration
type MigrationSpec struct {
	Purpose string `json:"purpose"`
	Namespace string `json:"namespace"`
	Node string `json:"node"`
	Podname string `json:"pod"`
	DestinationNode string `json:"destinationnode,omitempty"`
	Period int64 `json:"period,omitempty"`
	Pod struct{
		Type metav1.TypeMeta `json:"type,omitempty"`
		Object metav1.ObjectMeta `json:"object,omitempty"`
		PodSpec v1.PodSpec `json:"podspec,omitempty"`
	}

}


type MigrationStatus struct {
	LastCheckpointCreate map[string]int64 `json:"lastcheckpointcreate,omitempty"`
	PodStatus v1.PodStatus `json:"podstatus,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Migration is the Schema for the migrations API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=migrations,scope=Namespaced
type Migration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MigrationSpec   `json:"spec"`
	Status MigrationStatus `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigrationList contains a list of Migration
type MigrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Migration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Migration{}, &MigrationList{})
}
