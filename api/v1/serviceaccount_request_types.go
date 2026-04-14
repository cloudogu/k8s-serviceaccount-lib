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

type ServiceAccountRequestConsumerType string

const (
	DoguConsumerType      ServiceAccountRequestConsumerType = "Dogu"
	ComponentConsumerType ServiceAccountRequestConsumerType = "Component"
)

// LocalSecretRef definiert eine Referenz auf ein Secret im selben Namespace.
type LocalSecretRef struct {
	// Name des Secrets.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	Name string `json:"name"`
}

type ServiceAccountRequestParams struct {
}

// ServiceAccountRequestSpec defines the desired state of ServiceAccountRequest.
type ServiceAccountRequestSpec struct {
	// Optional defines if the request is optional and can be retried if the producer is currently not available.
	// +optional
	// +kubebuilder:default=false
	// +kubebuilder:example=false
	Optional bool `json:"optional,omitempty"`

	// Consumer defines the name of the requesting Dogu or Component
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:example=grafana
	Consumer string `json:"consumer"`

	// ConsumerType defines the type of the requester. This is necessary for restarting the requester if the account is optional and will be created later.
	// +required
	// +kubebuilder:validation:Enum=Dogu;Component
	// +kubebuilder:example=Dogu
	ConsumerType string `json:"consumerType"`

	// Producer defines the producer who should process the service account request.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?$
	// +kubebuilder:example=k8s-prometheus
	Producer string `json:"producer"`

	// SecretRef defines the secret where the generated service account will be written to.
	// If this field is empty, the service account will be written to a secret named like the name of the serviceaccount request resource.
	// +optional
	SecretRef *LocalSecretRef `json:"secretRef,omitempty"`

	// Params defines the parameter which should be used when creating the service account.
	Params *ServiceAccountRequestParams `json:"params,omitempty"`
}

// ServiceAccountRequestStatus defines the observed state of ServiceAccountRequest.
type ServiceAccountRequestStatus struct {
	// Conditions represent the current state of the ServiceAccountRequest resource.
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
// +kubebuilder:resource:shortName=sare
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Optional",type="bool",JSONPath=".spec.optional",description="Optional nature of the resource"
// +kubebuilder:printcolumn:name="Consumer",type="string",JSONPath=".spec.consumer",description="Name of the requester"

// ServiceAccountRequest is the Schema for the serviceaccount request API.
type ServiceAccountRequest struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of ServiceAccountRequest.
	// +required
	Spec ServiceAccountRequestSpec `json:"spec"`

	// status defines the observed state of ServiceAccountRequest.
	// +optional
	Status ServiceAccountRequestStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// ServiceAccountRequestList contains a list of ServiceAccountRequest.
type ServiceAccountRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []ServiceAccountRequest `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServiceAccountRequest{}, &ServiceAccountRequestList{})
}
