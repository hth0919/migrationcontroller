package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// MigrationSpec defines the desired state of Migration
type MigrationSpec struct {
	Purpose string `json:"purpose"`
	Namespace string `json:"namespace"`
	Node string `json:"node"`
	Podname string `json:"podname"`
	DestinationNode string `json:"destinationnode,omitempty"`
	Period int64 `json:"period,omitempty"`
}


type MigrationStatus struct {
	LastCheckpointCreate map[string]int64 `json:"lastcheckpointcreate,omitempty"`
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
