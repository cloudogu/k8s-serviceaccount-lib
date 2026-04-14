/*
Copyright 2026.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServiceAccountProducerSpec defines the desired state of ServiceAccountProducer.
type ServiceAccountProducerSpec struct {
}

// ServiceAccountProducerStatus defines the observed state of ServiceAccountProducer.
type ServiceAccountProducerStatus struct {
	// Conditions represent the current state of the ServiceAccountProducer resource.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

// TODO
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=sapr
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type == 'Ready')].status",description="Whether the service account request is ready"
// +kubebuilder:printcolumn:name="External-Ports",type="string",JSONPath=".status.allocatedPorts",description="Allocated external ports",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of the resource"

// ServiceAccountProducer is the Schema for the serviceaccount request API.
type ServiceAccountProducer struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ServiceAccountProducer.
	// +required
	Spec ServiceAccountProducerSpec `json:"spec"`

	// status defines the observed state of ServiceAccountProducer.
	// +optional
	Status ServiceAccountProducerStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ServiceAccountProducerList contains a list of ServiceAccountProducer.
type ServiceAccountProducerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ServiceAccountProducer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceAccountProducer{}, &ServiceAccountProducerList{})
}
