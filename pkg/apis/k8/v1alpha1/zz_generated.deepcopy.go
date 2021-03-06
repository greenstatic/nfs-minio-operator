// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSMinio) DeepCopyInto(out *NFSMinio) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSMinio.
func (in *NFSMinio) DeepCopy() *NFSMinio {
	if in == nil {
		return nil
	}
	out := new(NFSMinio)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NFSMinio) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSMinioList) DeepCopyInto(out *NFSMinioList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NFSMinio, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSMinioList.
func (in *NFSMinioList) DeepCopy() *NFSMinioList {
	if in == nil {
		return nil
	}
	out := new(NFSMinioList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NFSMinioList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSMinioSpec) DeepCopyInto(out *NFSMinioSpec) {
	*out = *in
	out.NFS = in.NFS
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSMinioSpec.
func (in *NFSMinioSpec) DeepCopy() *NFSMinioSpec {
	if in == nil {
		return nil
	}
	out := new(NFSMinioSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSMinioSpecNFS) DeepCopyInto(out *NFSMinioSpecNFS) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSMinioSpecNFS.
func (in *NFSMinioSpecNFS) DeepCopy() *NFSMinioSpecNFS {
	if in == nil {
		return nil
	}
	out := new(NFSMinioSpecNFS)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NFSMinioStatus) DeepCopyInto(out *NFSMinioStatus) {
	*out = *in
	if in.SecretKeyHash != nil {
		in, out := &in.SecretKeyHash, &out.SecretKeyHash
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NFSMinioStatus.
func (in *NFSMinioStatus) DeepCopy() *NFSMinioStatus {
	if in == nil {
		return nil
	}
	out := new(NFSMinioStatus)
	in.DeepCopyInto(out)
	return out
}
