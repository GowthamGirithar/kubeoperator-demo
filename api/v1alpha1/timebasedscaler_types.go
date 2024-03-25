package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespacedName for the deployment details
type NamespacedName struct {
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// TimeBasedScalerSpec defines the desired state of TimeBasedScaler
type TimeBasedScalerSpec struct {
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	StartHour int32 `json:"start_hour,omitempty"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	EndHour      int32            `json:"end_hour,omitempty"`
	ReplicaCount int32            `json:"replica_count,omitempty"`
	Deployments  []NamespacedName `json:"deployments,omitempty"`
}

// TimeBasedScalerStatus defines the observed state of TimeBasedScaler
type TimeBasedScalerStatus struct {
	Status string `json:"status,omitempty"`
}

// TimeBasedScaler is the Schema for the timebasedscalers API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type TimeBasedScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TimeBasedScalerSpec   `json:"spec,omitempty"`
	Status TimeBasedScalerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TimeBasedScalerList contains a list of TimeBasedScaler
type TimeBasedScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TimeBasedScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TimeBasedScaler{}, &TimeBasedScalerList{})
}
