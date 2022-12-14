//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubSubscription) DeepCopyInto(out *GoogleCloudPubSubSubscription) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubSubscription.
func (in *GoogleCloudPubSubSubscription) DeepCopy() *GoogleCloudPubSubSubscription {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubSubscription)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GoogleCloudPubSubSubscription) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubSubscriptionList) DeepCopyInto(out *GoogleCloudPubSubSubscriptionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GoogleCloudPubSubSubscription, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubSubscriptionList.
func (in *GoogleCloudPubSubSubscriptionList) DeepCopy() *GoogleCloudPubSubSubscriptionList {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubSubscriptionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GoogleCloudPubSubSubscriptionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubSubscriptionSpec) DeepCopyInto(out *GoogleCloudPubSubSubscriptionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubSubscriptionSpec.
func (in *GoogleCloudPubSubSubscriptionSpec) DeepCopy() *GoogleCloudPubSubSubscriptionSpec {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubSubscriptionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubSubscriptionStatus) DeepCopyInto(out *GoogleCloudPubSubSubscriptionStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubSubscriptionStatus.
func (in *GoogleCloudPubSubSubscriptionStatus) DeepCopy() *GoogleCloudPubSubSubscriptionStatus {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubSubscriptionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubTopic) DeepCopyInto(out *GoogleCloudPubSubTopic) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubTopic.
func (in *GoogleCloudPubSubTopic) DeepCopy() *GoogleCloudPubSubTopic {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubTopic)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GoogleCloudPubSubTopic) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubTopicList) DeepCopyInto(out *GoogleCloudPubSubTopicList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GoogleCloudPubSubTopic, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubTopicList.
func (in *GoogleCloudPubSubTopicList) DeepCopy() *GoogleCloudPubSubTopicList {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubTopicList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GoogleCloudPubSubTopicList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubTopicSpec) DeepCopyInto(out *GoogleCloudPubSubTopicSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubTopicSpec.
func (in *GoogleCloudPubSubTopicSpec) DeepCopy() *GoogleCloudPubSubTopicSpec {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubTopicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleCloudPubSubTopicStatus) DeepCopyInto(out *GoogleCloudPubSubTopicStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleCloudPubSubTopicStatus.
func (in *GoogleCloudPubSubTopicStatus) DeepCopy() *GoogleCloudPubSubTopicStatus {
	if in == nil {
		return nil
	}
	out := new(GoogleCloudPubSubTopicStatus)
	in.DeepCopyInto(out)
	return out
}
