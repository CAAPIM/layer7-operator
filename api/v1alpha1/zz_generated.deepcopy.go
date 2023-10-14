//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnrollmentBundle) DeepCopyInto(out *EnrollmentBundle) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnrollmentBundle.
func (in *EnrollmentBundle) DeepCopy() *EnrollmentBundle {
	if in == nil {
		return nil
	}
	out := new(EnrollmentBundle)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayProxy) DeepCopyInto(out *GatewayProxy) {
	*out = *in
	if in.Gateways != nil {
		in, out := &in.Gateways, &out.Gateways
		*out = make([]ProxyGateway, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayProxy.
func (in *GatewayProxy) DeepCopy() *GatewayProxy {
	if in == nil {
		return nil
	}
	out := new(GatewayProxy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7Api) DeepCopyInto(out *L7Api) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7Api.
func (in *L7Api) DeepCopy() *L7Api {
	if in == nil {
		return nil
	}
	out := new(L7Api)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *L7Api) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7ApiList) DeepCopyInto(out *L7ApiList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]L7Api, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7ApiList.
func (in *L7ApiList) DeepCopy() *L7ApiList {
	if in == nil {
		return nil
	}
	out := new(L7ApiList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *L7ApiList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7ApiSpec) DeepCopyInto(out *L7ApiSpec) {
	*out = *in
	if in.DeploymentTags != nil {
		in, out := &in.DeploymentTags, &out.DeploymentTags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7ApiSpec.
func (in *L7ApiSpec) DeepCopy() *L7ApiSpec {
	if in == nil {
		return nil
	}
	out := new(L7ApiSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7ApiStatus) DeepCopyInto(out *L7ApiStatus) {
	*out = *in
	if in.Gateways != nil {
		in, out := &in.Gateways, &out.Gateways
		*out = make([]LinkedGatewayStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7ApiStatus.
func (in *L7ApiStatus) DeepCopy() *L7ApiStatus {
	if in == nil {
		return nil
	}
	out := new(L7ApiStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7Portal) DeepCopyInto(out *L7Portal) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7Portal.
func (in *L7Portal) DeepCopy() *L7Portal {
	if in == nil {
		return nil
	}
	out := new(L7Portal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *L7Portal) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7PortalList) DeepCopyInto(out *L7PortalList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]L7Portal, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7PortalList.
func (in *L7PortalList) DeepCopy() *L7PortalList {
	if in == nil {
		return nil
	}
	out := new(L7PortalList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *L7PortalList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7PortalSpec) DeepCopyInto(out *L7PortalSpec) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DeploymentTags != nil {
		in, out := &in.DeploymentTags, &out.DeploymentTags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Auth = in.Auth
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7PortalSpec.
func (in *L7PortalSpec) DeepCopy() *L7PortalSpec {
	if in == nil {
		return nil
	}
	out := new(L7PortalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *L7PortalStatus) DeepCopyInto(out *L7PortalStatus) {
	*out = *in
	if in.GatewayProxies != nil {
		in, out := &in.GatewayProxies, &out.GatewayProxies
		*out = make([]GatewayProxy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.EnrollmentBundle = in.EnrollmentBundle
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new L7PortalStatus.
func (in *L7PortalStatus) DeepCopy() *L7PortalStatus {
	if in == nil {
		return nil
	}
	out := new(L7PortalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LinkedGatewayStatus) DeepCopyInto(out *LinkedGatewayStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LinkedGatewayStatus.
func (in *LinkedGatewayStatus) DeepCopy() *LinkedGatewayStatus {
	if in == nil {
		return nil
	}
	out := new(LinkedGatewayStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PortalAuth) DeepCopyInto(out *PortalAuth) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PortalAuth.
func (in *PortalAuth) DeepCopy() *PortalAuth {
	if in == nil {
		return nil
	}
	out := new(PortalAuth)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ProxyGateway) DeepCopyInto(out *ProxyGateway) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ProxyGateway.
func (in *ProxyGateway) DeepCopy() *ProxyGateway {
	if in == nil {
		return nil
	}
	out := new(ProxyGateway)
	in.DeepCopyInto(out)
	return out
}
