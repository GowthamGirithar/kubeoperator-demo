/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TimeBasedScaler is the Schema for the timebasedscalers API
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
