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

	severityFields := map[string]any{tsparams.AlarmDefinitionSeverityField: "critical"}

	valid := oranapi.AlarmDictionary{
		AlarmDictionaryId:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
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
		name       string
		dictionary oranapi.AlarmDictionary
		wantErr    bool
	}{
		{
			name:       "valid",
			dictionary: valid,
		},
		{
			name: "nil alarmDictionaryId",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDictionaryId = uuid.Nil

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty schema version",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDictionarySchemaVersion = ""

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty alarm dictionary version",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDictionaryVersion = ""

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty entityType",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.EntityType = ""

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty vendor",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.Vendor = ""

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty alarm definitions",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = nil

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty alarmName",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:             "",
					AlarmDescription:      "test alarm",
					AlarmAdditionalFields: &severityFields,
				}}

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty alarmDescription",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:             "TestAlarm",
					AlarmDescription:      "",
					AlarmAdditionalFields: &severityFields,
				}}

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "missing severity key",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:        "TestAlarm",
					AlarmDescription: "test alarm",
					AlarmAdditionalFields: new(map[string]any{
						"other": "value",
					}),
				}}

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "nil additional fields",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:        "TestAlarm",
					AlarmDescription: "test alarm",
				}}

				return dictionary
			}(),
			wantErr: true,
		},
		{
			name: "empty severity value allowed",
			dictionary: func() oranapi.AlarmDictionary {
				dictionary := valid
				dictionary.AlarmDefinition = []oranapi.AlarmDefinition{{
					AlarmName:        "TestAlarm",
					AlarmDescription: "test alarm",
					AlarmAdditionalFields: new(map[string]any{
						tsparams.AlarmDefinitionSeverityField: "",
					}),
				}}

				return dictionary
			}(),
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := VerifyAlarmDictionaryStructure(testCase.dictionary)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
