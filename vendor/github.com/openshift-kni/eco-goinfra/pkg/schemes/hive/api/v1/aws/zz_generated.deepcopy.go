//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package aws

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AssumeRole) DeepCopyInto(out *AssumeRole) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AssumeRole.
func (in *AssumeRole) DeepCopy() *AssumeRole {
	if in == nil {
		return nil
	}
	out := new(AssumeRole)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EC2Metadata) DeepCopyInto(out *EC2Metadata) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EC2Metadata.
func (in *EC2Metadata) DeepCopy() *EC2Metadata {
	if in == nil {
		return nil
	}
	out := new(EC2Metadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EC2RootVolume) DeepCopyInto(out *EC2RootVolume) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EC2RootVolume.
func (in *EC2RootVolume) DeepCopy() *EC2RootVolume {
	if in == nil {
		return nil
	}
	out := new(EC2RootVolume)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MachinePoolPlatform) DeepCopyInto(out *MachinePoolPlatform) {
	*out = *in
	if in.Zones != nil {
		in, out := &in.Zones, &out.Zones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Subnets != nil {
		in, out := &in.Subnets, &out.Subnets
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.EC2RootVolume = in.EC2RootVolume
	if in.SpotMarketOptions != nil {
		in, out := &in.SpotMarketOptions, &out.SpotMarketOptions
		*out = new(SpotMarketOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.EC2Metadata != nil {
		in, out := &in.EC2Metadata, &out.EC2Metadata
		*out = new(EC2Metadata)
		**out = **in
	}
	if in.AdditionalSecurityGroupIDs != nil {
		in, out := &in.AdditionalSecurityGroupIDs, &out.AdditionalSecurityGroupIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.UserTags != nil {
		in, out := &in.UserTags, &out.UserTags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MachinePoolPlatform.
func (in *MachinePoolPlatform) DeepCopy() *MachinePoolPlatform {
	if in == nil {
		return nil
	}
	out := new(MachinePoolPlatform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Metadata) DeepCopyInto(out *Metadata) {
	*out = *in
	if in.HostedZoneRole != nil {
		in, out := &in.HostedZoneRole, &out.HostedZoneRole
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Metadata.
func (in *Metadata) DeepCopy() *Metadata {
	if in == nil {
		return nil
	}
	out := new(Metadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Platform) DeepCopyInto(out *Platform) {
	*out = *in
	out.CredentialsSecretRef = in.CredentialsSecretRef
	if in.CredentialsAssumeRole != nil {
		in, out := &in.CredentialsAssumeRole, &out.CredentialsAssumeRole
		*out = new(AssumeRole)
		**out = **in
	}
	if in.UserTags != nil {
		in, out := &in.UserTags, &out.UserTags
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.PrivateLink != nil {
		in, out := &in.PrivateLink, &out.PrivateLink
		*out = new(PrivateLinkAccess)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Platform.
func (in *Platform) DeepCopy() *Platform {
	if in == nil {
		return nil
	}
	out := new(Platform)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlatformStatus) DeepCopyInto(out *PlatformStatus) {
	*out = *in
	if in.PrivateLink != nil {
		in, out := &in.PrivateLink, &out.PrivateLink
		*out = new(PrivateLinkAccessStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlatformStatus.
func (in *PlatformStatus) DeepCopy() *PlatformStatus {
	if in == nil {
		return nil
	}
	out := new(PlatformStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateLinkAccess) DeepCopyInto(out *PrivateLinkAccess) {
	*out = *in
	if in.AdditionalAllowedPrincipals != nil {
		in, out := &in.AdditionalAllowedPrincipals, &out.AdditionalAllowedPrincipals
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateLinkAccess.
func (in *PrivateLinkAccess) DeepCopy() *PrivateLinkAccess {
	if in == nil {
		return nil
	}
	out := new(PrivateLinkAccess)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrivateLinkAccessStatus) DeepCopyInto(out *PrivateLinkAccessStatus) {
	*out = *in
	in.VPCEndpointService.DeepCopyInto(&out.VPCEndpointService)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrivateLinkAccessStatus.
func (in *PrivateLinkAccessStatus) DeepCopy() *PrivateLinkAccessStatus {
	if in == nil {
		return nil
	}
	out := new(PrivateLinkAccessStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SpotMarketOptions) DeepCopyInto(out *SpotMarketOptions) {
	*out = *in
	if in.MaxPrice != nil {
		in, out := &in.MaxPrice, &out.MaxPrice
		*out = new(string)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SpotMarketOptions.
func (in *SpotMarketOptions) DeepCopy() *SpotMarketOptions {
	if in == nil {
		return nil
	}
	out := new(SpotMarketOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VPCEndpointService) DeepCopyInto(out *VPCEndpointService) {
	*out = *in
	if in.DefaultAllowedPrincipal != nil {
		in, out := &in.DefaultAllowedPrincipal, &out.DefaultAllowedPrincipal
		*out = new(string)
		**out = **in
	}
	if in.AdditionalAllowedPrincipals != nil {
		in, out := &in.AdditionalAllowedPrincipals, &out.AdditionalAllowedPrincipals
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VPCEndpointService.
func (in *VPCEndpointService) DeepCopy() *VPCEndpointService {
	if in == nil {
		return nil
	}
	out := new(VPCEndpointService)
	in.DeepCopyInto(out)
	return out
}
