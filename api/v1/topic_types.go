/*
Copyright 2022.

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

// TopicSpec defines the desired state of Topic
type TopicSpec struct {
	// ID of project
	//+kubebuilder:validation:XValidation:message="Immutable field",rule="self == oldSelf"
	ProjectID string `json:"projectID,omitempty"`

	// ID of topic
	//+kubebuilder:validation:XValidation:message="Immutable field",rule="self == oldSelf"
	TopicID string `json:"topicID,omitempty"`

	// Authorative IAM Binding for topic
	//+kubebuilder:validation:Optional
	Bindings []IamBinding `json:"bindings,omitempty"`
}

// TopicStatus defines the observed state of Topic
type TopicStatus struct {
	Phase TopicStatusPhase `json:"phase,omitempty"`

	// Message is an error message of Cloud Pub/Sub.
	// Available only if Phase is Error.
	Message string `json:"message,omitempty"`
}

type TopicStatusPhase string

const (
	TopicStatusPhaseActive TopicStatusPhase = "Active"
	TopicStatusPhaseError  TopicStatusPhase = "Error"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Topic is the Schema for the topics API
type Topic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TopicSpec   `json:"spec,omitempty"`
	Status TopicStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TopicList contains a list of Topic
type TopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Topic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Topic{}, &TopicList{})
}
