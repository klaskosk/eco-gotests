package profiles

import (
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Needs/user stories:
//
//  - I want to be able to get all the slave interfaces of a node.
//  - I want to be able to check if a node has X profiles of type Y.
//  - I want to be able to get all the profiles of type X on a node.
//  - I want to be able to get all interfaces of a profile with type/role X.

// PtpProfileType enumerates the supported types of profiles.
type PtpProfileType int

const (
	// ProfileTypeOC refers to a PTP profile with a single interface set to client only. It is an ordinary clock
	// profile.
	ProfileTypeOC PtpProfileType = iota
	// ProfileTypeTwoPortOC refers to a PTP profile with two interfaces set to client only. Only one of these
	// interfaces will be active at a time.
	ProfileTypeTwoPortOC
	// ProfileTypeBC refers to a PTP profile in a boundary clock configuration, i.e., one client interface and one
	// server interface.
	ProfileTypeBC
	// ProfileTypeHA refers to a PTP profile that does not correspond to individual interfaces but indicates other
	// profiles are in a highly available configuration.
	ProfileTypeHA
	// ProfileTypeGM refers to a PTP profile for one NIC with all interfaces set to server only.
	ProfileTypeGM
	// ProfileTypeMultiNICGM refers to a PTP profile for multiple NICs where all interfaces are set to server only.
	// SMA cables are used to synchronize the NICs so they can all act as grand masters.
	ProfileTypeMultiNICGM
)

// PtpClockType enumerates the roles of each interface. It is different from the roles in metrics, which include extra
// runtime values not represented in the profile. The zero value is a client and only serverOnly (or masterOnly) values
// of 1 indicate a server.
type PtpClockType int

const (
	// ClockTypeClient indicates an interface is acting as a follower of time signals. Formerly slave.
	ClockTypeClient PtpClockType = iota
	// ClockTypeServer indicates an interface is acting as a leader of time signals. Formerly master.
	ClockTypeServer
)

// ProfileReference contains the information needed to identify a profile on a cluster.
type ProfileReference struct {
	// ConfigReference is the reference to the PtpConfig object that contains the profile.
	ConfigReference runtimeclient.ObjectKey
	// ProfileIndex is the index of the profile in the PtpConfig object.
	ProfileIndex int
	// ProfileName is the name of the profile. It is not necessary to get the profile directly, but is used as a key
	// when recommending profiles to nodes.
	ProfileName string
}

// ProfileInfo contains information about a PTP profile. Since profiles can be readily retrieved from the cluster, it
// only contains information that must be parsed and a reference to the profile on the cluster.
type ProfileInfo struct {
	ProfileType PtpProfileType
	Reference   ProfileReference
	// Interfaces is a map of interface names to a struct holding more detailed information. Values should never be
	// nil.
	Interfaces map[string]*InterfaceInfo
}

// GetInterfacesByClockType returns a slice of InterfaceInfo pointers for each interface in the profile matching the
// provided clockType. Elements are guaranteed not to be nil.
func (profileInfo *ProfileInfo) GetInterfacesByClockType(clockType PtpClockType) []*InterfaceInfo {
	var interfaces []*InterfaceInfo

	for _, interfaceInfo := range profileInfo.Interfaces {
		if interfaceInfo.ClockType == clockType {
			interfaces = append(interfaces, interfaceInfo)
		}
	}

	return interfaces
}

// InterfaceInfo contains information about the PTP clock type of an interface. In the future, it may also contain
// information about which interface it is connected to.
type InterfaceInfo struct {
	Name      string
	ClockType PtpClockType
}

// ProfileCounts records the number of profiles of each type. It is provided as a map rather than a struct to allow
// indexing using the profile type.
type ProfileCounts map[PtpProfileType]uint

// NodeInfo contains all the PTP config-related information for a single node. Common operations are provided as methods
// on this type to avoid the need to aggregate and query nested data.
type NodeInfo struct {
	// Name is the name of the node resource this struct is associated to.
	Name string
	// Counts records the number of each profile type recommended to this node. It will never be nil when this
	// struct is returned from a function in this package.
	Counts ProfileCounts
	// Profiles contains a list of information structs corresponding to each profile that is recommended to this
	// node. Elements should never be nil.
	Profiles []*ProfileInfo
}

// GetInterfacesByClockType returns a slice of InterfaceInfo pointers for each interface across all profiles on this
// node matching the provided clockType. Elements are guaranteed not to be nil.
func (nodeInfo *NodeInfo) GetInterfacesByClockType(clockType PtpClockType) []*InterfaceInfo {
	var nodeInterfaces []*InterfaceInfo

	for _, profileInfo := range nodeInfo.Profiles {
		nodeInterfaces = append(nodeInterfaces, profileInfo.GetInterfacesByClockType(clockType)...)
	}

	return nodeInterfaces
}

// GetProfilesByType returns a slice of ProfileInfo pointers for each profile on this node matching the provided
// profileType. Elements are guaranteed not to be nil.
func (nodeInfo *NodeInfo) GetProfilesByType(profileType PtpProfileType) []*ProfileInfo {
	var nodeProfiles []*ProfileInfo

	for _, profileInfo := range nodeInfo.Profiles {
		if profileInfo.ProfileType == profileType {
			nodeProfiles = append(nodeProfiles, profileInfo)
		}
	}

	return nodeProfiles
}
