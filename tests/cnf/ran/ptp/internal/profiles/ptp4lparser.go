package profiles

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	ptpv1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/ptp/v1"
)

type parsedPtp4lConf struct {
	sections    map[string]map[string]string
	interfaces  map[string]PtpClockType
	profileType PtpProfileType
}

// parsePtpProfile parses the ptp4l configuration file and returns a parsedPtp4lConf struct.
// It first checks Ptp4lOpts and determines if the client flag is set.
// If interface is set on the profile, it gets added to the ones parsed from the ptp4l configuration file.
func parsePtpProfile(profile ptpv1.PtpProfile) (*parsedPtp4lConf, error) {
	parsedConfig := &parsedPtp4lConf{
		sections: make(map[string]map[string]string),
	}
	clientFlag := hasClientFlag(profile.Ptp4lOpts)

	var err error
	if profile.Ptp4lConf != nil && *profile.Ptp4lConf != "" {
		parsedConfig.sections, err = getSectionsFromPtp4lConf(*profile.Ptp4lConf)
		if err != nil {
			return nil, fmt.Errorf("failed to get sections from ptp4lConf: %w", err)
		}
	}

	parsedConfig.interfaces = getInterfacesFromPtp4lSections(clientFlag, parsedConfig.sections)

	if profile.Interface != nil && *profile.Interface != "" {
		// If the interface is not set in the config file, it must be client only.
		if _, ok := parsedConfig.interfaces[*profile.Interface]; !ok {
			parsedConfig.interfaces[*profile.Interface] = ClockTypeClient
		}
	}

	parsedConfig.profileType, err = determineProfileType(parsedConfig.interfaces, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to determine profile type: %w", err)
	}

	return parsedConfig, nil
}

// getSectionsFromPtp4lConf parses the ptp4l configuration file and returns a map of sections and their key-value pairs.
func getSectionsFromPtp4lConf(ptp4lConf string) (map[string]map[string]string, error) {
	var (
		currentSectionName string
		currentSectionMap  map[string]string
	)

	sections := make(map[string]map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(ptp4lConf))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore empty lines and comments.
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Lines with text between brackets are considered section names.
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") && len(line) > 2 {
			if currentSectionName != "" {
				sections[currentSectionName] = currentSectionMap
			}

			currentSectionName = line[1 : len(line)-1]

			if _, ok := sections[currentSectionName]; !ok {
				sections[currentSectionName] = make(map[string]string)
			}

			currentSectionMap = sections[currentSectionName]

			continue
		}

		// If the first section has not been found yet, skip the line.
		if currentSectionName == "" {
			continue
		}

		// This is not a section name, so it should be a key-value pair, separated by a space.
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		key := fields[0]
		value := strings.Join(fields[1:], " ")
		currentSectionMap[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ptp4l configuration: %w", err)
	}

	// Add the last section if it exists.
	if currentSectionName != "" {
		sections[currentSectionName] = currentSectionMap
	}

	return sections, nil
}

// getInterfacesFromPtp4lSections extracts the interfaces and their clock types from the ptp4l configuration sections.
// The provided clientFlag indicates whether the clientOnly command line flag is set in ptp4lOpts.
func getInterfacesFromPtp4lSections(clientFlag bool, sections map[string]map[string]string) map[string]PtpClockType {
	interfaces := make(map[string]PtpClockType)

	if globalSection, ok := sections["global"]; ok && globalSection != nil {
		// slaveOnly is deprecated but still used and supported by ptp4l.
		if globalSection["clientOnly"] == "1" || globalSection["slaveOnly"] == "1" {
			clientFlag = true
		}
	}

	for sectionName, sectionValues := range sections {
		if sectionName == "global" || sectionName == "unicast_master_table" {
			continue
		}

		if clientFlag {
			interfaces[sectionName] = ClockTypeClient

			continue
		}

		// masterOnly is deprecated but still used and supported by ptp4l, similar to slaveOnly.
		if sectionValues["serverOnly"] == "1" || sectionValues["masterOnly"] == "1" {
			interfaces[sectionName] = ClockTypeServer

			continue
		}

		interfaces[sectionName] = ClockTypeClient
	}

	return interfaces
}

func determineProfileType(interfaces map[string]PtpClockType, profile ptpv1.PtpProfile) (PtpProfileType, error) {
	numInterfaces := len(interfaces)
	numClientInterfaces := 0
	numServerInterfaces := 0

	// To track the number of NICs in this profile, keep a set of NIC IDs. A NIC ID is the interface name without
	// the last character. Since we only use this for count, no need to append an "x" to the end.
	nicIDSet := make(map[string]struct{})

	for iface, clockType := range interfaces {
		// Although the interface should not be empty, check again to avoid a panic when indexing the string.
		if len(iface) == 0 {
			continue
		}

		nicIDSet[iface[:len(iface)-1]] = struct{}{}

		switch clockType {
		case ClockTypeClient:
			numClientInterfaces++
		case ClockTypeServer:
			numServerInterfaces++
		}
	}

	// If the profile has PtpSettings and haProfiles is set, return ProfileTypeHA.
	if profile.PtpSettings != nil {
		if haProfiles, ok := profile.PtpSettings["haProfiles"]; ok && haProfiles != "" {
			return ProfileTypeHA, nil
		}
	}

	switch {
	// If the profile has one interface and one client interface, return ProfileTypeOC.
	case numInterfaces == 1 && numClientInterfaces == 1:
		return ProfileTypeOC, nil
	// If the profile has two interfaces and two client interfaces, return ProfileTypeTwoPortOC.
	case numInterfaces == 2 && numClientInterfaces == 2:
		return ProfileTypeTwoPortOC, nil
	// If the profile has at least two interfaces and only one client interface, return ProfileTypeBC.
	case numInterfaces >= 2 && numClientInterfaces == 1:
		return ProfileTypeBC, nil
	// If the profile has one NIC and all interfaces are servers, return ProfileTypeGM.
	case len(nicIDSet) == 1 && numServerInterfaces == numInterfaces:
		return ProfileTypeGM, nil
	// If the profile has multiple NICs and all interfaces are servers, return ProfileTypeMultiNICGM.
	case len(nicIDSet) > 1 && numServerInterfaces == numInterfaces:
		return ProfileTypeMultiNICGM, nil
	// All other profile types are considered unsupported.
	default:
		return 0, fmt.Errorf("unable to determine PTP profile type based on defined rules")
	}
}

// hasClientFlag checks if the ptp4lOpts string contains any client-only flags.
func hasClientFlag(ptp4lOpts *string) bool {
	if ptp4lOpts == nil {
		return false
	}

	fields := strings.Fields(*ptp4lOpts)
	flags := []string{"-s", "--clientOnly", "--clientOnly=1", "--slaveOnly", "--slaveOnly=1"}

	for _, field := range fields {
		if slices.Contains(flags, field) {
			return true
		}
	}

	return false
}
