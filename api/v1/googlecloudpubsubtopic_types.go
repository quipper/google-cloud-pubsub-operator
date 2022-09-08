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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GoogleCloudPubSubTopicSpec defines the desired state of GoogleCloudPubSubTopic
type GoogleCloudPubSubTopicSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ID of project
	ProjectID string `json:"projectID,omitempty"`

	// ID of topic
	TopicID string `json:"topicID,omitempty"`
}

// GoogleCloudPubSubTopicStatus defines the observed state of GoogleCloudPubSubTopic
type GoogleCloudPubSubTopicStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GoogleCloudPubSubTopic is the Schema for the googlecloudpubsubtopics API
type GoogleCloudPubSubTopic struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleCloudPubSubTopicSpec   `json:"spec,omitempty"`
	Status GoogleCloudPubSubTopicStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoogleCloudPubSubTopicList contains a list of GoogleCloudPubSubTopic
type GoogleCloudPubSubTopicList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoogleCloudPubSubTopic `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleCloudPubSubTopic{}, &GoogleCloudPubSubTopicList{})
}
