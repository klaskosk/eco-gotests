package o2imstest

import (
	"testing"

	"github.com/google/uuid"
	oranapi "github.com/rh-ecosystem-edge/eco-goinfra/pkg/oran/api"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/oran/internal/tsparams"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen // Table-driven cases stay inline for readability.
func TestVerifyAlarmDictionaryStructure(t *testing.T) {
	t.Parallel()

	dictionaryID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	severityFields := map[string]any{tsparams.AlarmDefinitionSeverityField: "critical"}

	valid := oranapi.AlarmDictionary{
		AlarmDictionaryId:            dictionaryID,
		AlarmDictionarySchemaVersion: "1.0.0",
		AlarmDictionaryVersion:       "1.0.0",
		EntityType:                   "NodeClusterType",
		Vendor:                       "redhat",
		AlarmDefinition: []oranapi.AlarmDefinition{{
			AlarmName:             "TestAlarm",
			AlarmDescription:      "test alarm",
			AlarmAdditionalFields: &severityFields,
		}},
	}

	tests := []struct {
		name    string
		mutate  func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary
		wantErr bool
	}{
		{
			name:   "valid",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary { return dictionary },
		},
		{
			name: "nil alarmDictionaryId",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				dictionary.AlarmDictionaryId = uuid.Nil

				return dictionary
			},
			wantErr: true,
		},
		{
			name: "empty schema version",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				dictionary.AlarmDictionarySchemaVersion = ""

				return dictionary
			},
			wantErr: true,
		},
		{
			name: "empty alarm definitions",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				dictionary.AlarmDefinition = nil

				return dictionary
			},
			wantErr: true,
		},
		{
			name: "missing severity key",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				fields := map[string]any{"other": "value"}
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:             "TestAlarm",
					AlarmDescription:      "test alarm",
					AlarmAdditionalFields: &fields,
				}}

				return dictionary
			},
			wantErr: true,
		},
		{
			name: "nil additional fields",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:        "TestAlarm",
					AlarmDescription: "test alarm",
				}}

				return dictionary
			},
			wantErr: true,
		},
		{
			name: "empty severity value allowed",
			mutate: func(dictionary oranapi.AlarmDictionary) oranapi.AlarmDictionary {
				fields := map[string]any{tsparams.AlarmDefinitionSeverityField: ""}
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:             "TestAlarm",
					AlarmDescription:      "test alarm",
					AlarmAdditionalFields: &fields,
				}}

				return dictionary
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyAlarmDictionaryStructure(testCase.mutate(valid))
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
