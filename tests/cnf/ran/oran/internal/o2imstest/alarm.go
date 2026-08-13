package o2imstest

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
)

// VerifyAlarmDictionaryStructure checks that an AlarmDictionary has the expected top-level fields and definitions.
func VerifyAlarmDictionaryStructure(dictionary oranapi.AlarmDictionary) error {
	var errs []error

	if dictionary.AlarmDictionaryId == uuid.Nil {
		errs = append(errs, fmt.Errorf("alarmDictionaryId: want non-nil UUID, got %s", dictionary.AlarmDictionaryId))
	}

	if dictionary.AlarmDictionarySchemaVersion == "" {
		errs = append(errs, fmt.Errorf("alarmDictionarySchemaVersion: want non-empty"))
	}

	if dictionary.AlarmDictionaryVersion == "" {
		errs = append(errs, fmt.Errorf("alarmDictionaryVersion: want non-empty"))
	}

	if dictionary.EntityType == "" {
		errs = append(errs, fmt.Errorf("entityType: want non-empty"))
	}

	if dictionary.Vendor == "" {
		errs = append(errs, fmt.Errorf("vendor: want non-empty"))
	}

	if len(dictionary.AlarmDefinition) == 0 {
		errs = append(errs, fmt.Errorf("alarmDefinition: want non-empty"))
	}

	for idx, definition := range dictionary.AlarmDefinition {
		if definition.AlarmName == "" {
			errs = append(errs, fmt.Errorf("alarmDefinition[%d].alarmName: want non-empty", idx))
		}

		if definition.AlarmDescription == "" {
			errs = append(errs, fmt.Errorf("alarmDefinition[%d].alarmDescription: want non-empty", idx))
		}

		// Severity lives in additionalFields; require the key but allow empty values when a Prometheus
		// rule omits the severity label.
		if definition.AlarmAdditionalFields == nil {
			errs = append(errs, fmt.Errorf(
				"alarmDefinition[%d].alarmAdditionalFields: want non-nil with %s key",
				idx, tsparams.AlarmDefinitionSeverityField))

			continue
		}

		if _, ok := (*definition.AlarmAdditionalFields)[tsparams.AlarmDefinitionSeverityField]; !ok {
			errs = append(errs, fmt.Errorf(
				"alarmDefinition[%d].alarmAdditionalFields.%s: want key present",
				idx, tsparams.AlarmDefinitionSeverityField))
		}
	}

	return errors.Join(errs...)
}
