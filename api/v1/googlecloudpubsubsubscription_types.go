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

// GoogleCloudPubSubSubscriptionSpec defines the desired state of GoogleCloudPubSubSubscription
type GoogleCloudPubSubSubscriptionSpec struct {
	// subscription ID
	SubscriptionID string `json:"subscriptionID,omitempty"`

	// project ID of subscription
	SubscriptionProjectID string `json:"subscriptionProjectID,omitempty"`

	// topic ID
	TopicID string `json:"topicID,omitempty"`

	// project ID of topic
	TopicProjectID string `json:"topicProjectID,omitempty"`
}

// GoogleCloudPubSubSubscriptionStatus defines the observed state of GoogleCloudPubSubSubscription
type GoogleCloudPubSubSubscriptionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GoogleCloudPubSubSubscription is the Schema for the googlecloudpubsubsubscriptions API
type GoogleCloudPubSubSubscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleCloudPubSubSubscriptionSpec   `json:"spec,omitempty"`
	Status GoogleCloudPubSubSubscriptionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GoogleCloudPubSubSubscriptionList contains a list of GoogleCloudPubSubSubscription
type GoogleCloudPubSubSubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoogleCloudPubSubSubscription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleCloudPubSubSubscription{}, &GoogleCloudPubSubSubscriptionList{})
}
