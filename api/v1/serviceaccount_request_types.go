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

// ServiceAccountRequestParams defines the param structure for the service account creation with unordered options as a mapping from map string to a slice from strings.
// We do not use just string to string because some parameters maybe require multiple values of the same parameter/option (e.g. --repo x --repo y --permissions read).
// Although some parameters are unnamed and positional, so we need the additional args field for this. It defines positional arguments without names.
type ServiceAccountRequestParams struct {
	// +kubebuilder:validation:MaxProperties=20
	// +optional
	Options map[string][]string `json:"options,omitempty"`

	// +kubebuilder:validation:MaxItems=50
	// +optional
	Args []string `json:"args,omitempty"`
}

// ServiceAccountRotation defines a plan for recreating the service account credentials. If ServiceAccountRotation.Enabled is false, nothing will happen.
// If ServiceAccountRotation.Enabled and a cron expression is set with ServiceAccountRotation.Rotation the account will be recreated defined with ServiceAccountRotation.Rotation.
// +kubebuilder:validation:XValidation:rule="self.enabled == true ? self.rotation.size() > 0 : true",message="rotation must be set if enabled is true"
type ServiceAccountRotation struct {
	// +kubebuilder:default=false
	Enabled bool `json:"enabled"`
	// +optional
	// +kubebuilder:validation:Pattern=`^(@(annually|yearly|monthly|weekly|daily|hourly)|(((\d+,)*\d+|(\d+(\/|-)\d+)|\*)\s?){5,6})$`
	Rotation string `json:"rotation,omitempty"`
}

// ServiceAccountRequestSpec defines the desired state of ServiceAccountRequest.
type ServiceAccountRequestSpec struct {
	// Optional defines if the request is optional and can be retried if the producer is currently not available.
	// +optional
	// +kubebuilder:default=false
	// +kubebuilder:example=false
	Optional bool `json:"optional,omitempty"`

	// Consumer defines the name of the requesting Dogu or Component.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Consumer is immutable"
	// +kubebuilder:example=grafana
	Consumer string `json:"consumer"`

	// ConsumerType defines the type of the requester. This is necessary for restarting the requester if the account is optional and will be created later.
	// +required
	// +kubebuilder:validation:Enum=Dogu;Component
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="ConsumerType is immutable"
	// +kubebuilder:example=Dogu
	ConsumerType ServiceAccountRequestConsumerType `json:"consumerType"`

	// Producer defines the producer who should process the service account request. The value has to correspond to the name of the producer resource.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Producer is immutable"
	// +kubebuilder:example=prometheus
	Producer string `json:"producer"`

	// SecretRef defines the secret where the generated service account will be written to.
	// If this field is empty, the service account will be written to a secret named like the name of the serviceaccount request resource.
	// The final secret reference will be written to the status object.
	// +optional
	SecretRef *LocalSecretRef `json:"secretRef,omitempty"`

	// Params defines the parameter which should be used when creating the service account.
	// +optional
	Params *ServiceAccountRequestParams `json:"params,omitempty"`

	// Rotation defines timings when the service account should be recreated.
	// +optional
	Rotation *ServiceAccountRotation `json:"rotation,omitempty"`
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

	// LastRotation defines the timestamp of the last rotation.
	LastRotation metav1.Time `json:"lastRotation,omitempty"`

	// SecretRef defines the actual secret where the service account was written to.
	SecretRef *LocalSecretRef `json:"secretRef,omitempty"`
}

// ServiceAccountRequest define the structure of a service account request.
// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=sare
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Optional",type="bool",JSONPath=".spec.optional",description="Optional nature of the resource"
// +kubebuilder:printcolumn:name="Consumer",type="string",JSONPath=".spec.consumer"
// +kubebuilder:printcolumn:name="ConsumerType",type="string",JSONPath=".spec.consumerType",description="Type of the requester"
// +kubebuilder:printcolumn:name="Producer",type="string",JSONPath=".spec.producer",description="Producer of the requester"
// +kubebuilder:printcolumn:name="Secret",type="string",JSONPath=".status.secretRef.name",description="Actual secret name"
// +kubebuilder:printcolumn:name="Rotation",type="boolean",JSONPath=".spec.rotation.enabled"
// +kubebuilder:printcolumn:name="Schedule",type="string",JSONPath=".spec.rotation.rotation"
// +kubebuilder:printcolumn:name="Last-Run",type="date",JSONPath=".status.lastRotation"
// +kubebuilder:printcolumn:name="Params",type="string",JSONPath=".spec.params",priority=1
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
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
