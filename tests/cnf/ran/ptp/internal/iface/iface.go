// Package iface provides types and utilities for working with network interface names in the PTP test suite. Additional
// helpers are provided too for working with interfaces on running nodes.
package iface

// Name represents a network interface name. It provides a number of methods for working with interface names and is the
// canonical way to manipulate interface names in this test suite. The zero value is never valid.
type Name string

// GetNIC returns the NIC name associated with the interface. The NIC name is the interface name without the last
// character and an "x" appended. If the interface name already ends with "x", meaning it is already a NIC name, the
// NICName returned will be the same as the interface name.
//
// Names shorter than 2 characters are invalid and this method will return the zero value.
//
// Since this method is an identity operation for NIC names, it will return the special NIC names [ClockRealtime] and
// [Master] unchanged.
func (iface Name) GetNIC() NICName {
	if len(iface) < 2 {
		return ""
	}

	if NICName(iface) == ClockRealtime || NICName(iface) == Master {
		return NICName(iface)
	}

	return NICName(iface[:len(iface)-1] + "x")
}

// NICName represents a network interface name. It can be derived from [Name] using the [GetNIC] method. The zero value
// is never valid.
type NICName string

const (
	// ClockRealtime is the name of the NIC representing the realtime clock. It is not actually a NIC but appears as
	// one in some PTP metrics.
	ClockRealtime NICName = "CLOCK_REALTIME"
	// Master is the name of the NIC representing the master clock. It is not actually a NIC but appears as one in
	// some PTP metrics.
	Master NICName = "master"
)

// EnsureNIC verifies that the NIC name is actually a NIC name. If it is instead an interface name, it will be converted
// to a NIC name. If the NIC name is invalid, the zero value will be returned.
func (nic NICName) EnsureNIC() NICName {
	if nic[len(nic)-1] == 'x' {
		return nic
	}

	// Since [GetNIC] handles the case of [ClockRealtime] and [Master], we can just call it here rather than
	// duplicating the special case.
	return Name(nic).GetNIC()
}

// GroupInterfacesByNIC takes a slice of interface names and and returns a map of NIC names to slices of interface
// names. Invalid interface names are ignored, so not all inputs will be necessarily present in the output. The returned
// map is guaranteed to not be nil.
func GroupInterfacesByNIC(ifaces []Name) map[NICName][]Name {
	nicMap := make(map[NICName][]Name)

	for _, iface := range ifaces {
		nic := iface.GetNIC()
		if nic == "" {
			continue
		}

		nicMap[nic] = append(nicMap[nic], iface)
	}

	return nicMap
}
