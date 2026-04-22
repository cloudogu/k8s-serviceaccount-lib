package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConditionTypeReady                       = "Ready"
	ConditionReadyReasonAuthSecretNotFound   = "AuthSecretNotFound"
	ConditionReadyReasonConnectionFailed     = "ConnectionFailed"     // the endpoint is not reachable e.g., because of missing netpols
	ConditionReadyReasonInvalidConfiguration = "InvalidConfiguration" // label selector is invalid and no pods were found
)

// ServiceAccountProducerAuthSecret defines the reference of the secret that should be used to create the service account.
type ServiceAccountProducerAuthSecret struct {
	LocalSecretRef `json:",inline"`
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=^[-_.a-zA-Z0-9]+$
	// +kubebuilder:example=auth
	Key string `json:"key"`
}

// HTTPProducer defines necessary information to create a service account via HTTP.
type HTTPProducer struct {
	// Endpoint describes the url used to create service accounts.
	// This endpoint should implement the following verbs: PUT, POST, DELETE to create, update and delete service accounts.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:example="https://nexus:8081/api/cloudogu/service-account"
	Endpoint string `json:"endpoint"`
	// AuthSecret describes the reference of the secret that should be used to create the service account.
	// +required
	AuthSecret ServiceAccountProducerAuthSecret `json:"authSecret"`
	// Priority describes the priority of choosing a way to create the service account if multiple producers (e.g. http and exec) a configured.
	// +optional
	Priority int `json:"priority,omitzero"`
	// Params defines the parameters that should and can be used to create the service account.
	// +optional
	Params *ProducerParams `json:"params,omitempty"`
	// Return defines the structure of the return values from the producer.
	// It is used to write the actual service account data to the secret.
	// +required
	Return map[AttributeName]ProducerReturnDefinition `json:"return,omitempty"`
}

// ExecProducer defines necessary information to create a service account via k8s exec api.
type ExecProducer struct {
	// Command describes the command used to create service accounts.
	// +required
	// +kubebuilder:example=/create-sa.sh
	Command string `json:"command"`
	// Selector describes the label selector used to select the pods that should be used to create service accounts.
	// +required
	Selector metav1.LabelSelector `json:"selector"`
	// Priority describes the priority of choosing a way to create the service account if multiple producers a configured.
	// +optional
	Priority int `json:"priority,omitzero"`
	// Params defines the parameters that should and can be used to create the service account.
	// +optional
	Params *ProducerParams `json:"params,omitempty"`
	// Return defines the structure of the return values from the producer.
	// It is used to write the actual service account data to the secret.
	// +required
	Return map[AttributeName]ProducerReturnDefinition `json:"return,omitempty"`
}

// AttributeDefinition defines the structure of an attribute used to create service accounts.
type AttributeDefinition struct {
	// Description defines a human-readable description for the attribute.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Description string `json:"description,omitempty"`

	// Type defines the type of the attribute.
	// +required
	// +kubebuilder:validation:Enum=string;boolean;integer
	// +kubebuilder:default=string
	Type string `json:"type"`

	// Enum defines the possible values if AttributeDefinition.Type is string.
	// +optional
	// +kubebuilder:validation:MaxItems=20
	Enum []string `json:"enum,omitempty"`

	// Order defines the order of the attribute.
	// +optional
	Order int `json:"order,omitempty"`
}

// ProducerParams defines the parameters that should and can be used to create the service account.
type ProducerParams struct {
	// Attributes define the attributes that should be used to create the service account.
	// +kubebuilder:validation:MaxProperties=20
	// +kubebuilder:validation:MinProperties=1
	// +required
	Attributes map[AttributeName]AttributeDefinition `json:"attributes,omitempty"`

	// Required defines the attributes that are required to create the service account.
	// +optional
	Required []AttributeName `json:"required"`
}

// ProducerReturnDefinition defines the structure of the return values from the producer.
type ProducerReturnDefinition struct {
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Description string `json:"description"`
}

// AttributeName defines the name of an attribute.
type AttributeName string

// ServiceAccountProducerSpec defines the desired state of ServiceAccountProducer.
// +kubebuilder:validation:XValidation:rule="has(self.http) || has(self.exec)",message="At least one producer strategy (http or exec) must be defined"
type ServiceAccountProducerSpec struct {
	// Producer defines the name of the service account producing dogu or component.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?([.][a-z0-9]([-a-z0-9]*[a-z0-9])?)*$
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Producer is immutable"
	// +kubebuilder:example=prometheus
	Producer string `json:"producer"`
	// HTTP defines the necessary information to create a service account via HTTP.
	// +optional
	HTTP *HTTPProducer `json:"http,omitempty"`
	// Exec defines the necessary information to create a service account via k8s exec api.
	// +optional
	Exec *ExecProducer `json:"exec,omitempty"`
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
	// LastExecution defines the timestamp of the last execution.
	LastExecution metav1.Time `json:"lastExecution,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName=sapr
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type == 'Ready')].status",description="Whether the service account request is ready"
// +kubebuilder:printcolumn:name="Producer",type="string",JSONPath=".spec.producer"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.http.endpoint",priority=1
// +kubebuilder:printcolumn:name="Command",type="string",JSONPath=".spec.exec.command",priority=1
// +kubebuilder:printcolumn:name="Last Execution",type="date",JSONPath=".status.lastExecution"
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
