package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type exampleScenarioInput struct {
	ValueType ValueTypeEnum
	Scenario  TestValueScenario
}

type exampleScenarioExpected struct {
	Included bool
	HCL      string
}

func TestExampleValueFromTestScenarioPreservesAbsence(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{
			Name:     "test_undefined_scalar_is_omitted",
			Input:    exampleScenarioInput{ValueType: String},
			Expected: exampleScenarioExpected{},
		},
		{
			Name: "test_excluded_scalar_is_omitted",
			Input: exampleScenarioInput{
				ValueType: String,
				Scenario:  newTestValueScenario([]TestValueEntry{{ConfigValue: "unused"}}),
			},
			Expected: exampleScenarioExpected{},
		},
		{
			Name: "test_empty_string_is_included",
			Input: exampleScenarioInput{
				ValueType: String,
				Scenario:  newTestValueScenario([]TestValueEntry{{ConfigInclude: true}}),
			},
			Expected: exampleScenarioExpected{Included: true, HCL: `""`},
		},
		{
			Name: "test_empty_set_is_included",
			Input: exampleScenarioInput{
				ValueType: Set,
				Scenario:  newTestValueScenario(nil),
			},
			Expected: exampleScenarioExpected{Included: true, HCL: "[]"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(exampleScenarioInput)
			expected := testCase.Expected.(exampleScenarioExpected)
			property := Property{PropertyName: "value", AttributeName: "value", ValueType: input.ValueType}
			value := property.exampleValueFromTestScenario(&DataStore{ctx: NewContext()}, testClassName("fvTenant"), input.Scenario)

			if !expected.Included {
				if value != nil {
					t.Fatalf("example value %q should be omitted", value.HCL())
				}
				return
			}
			if value == nil || value.HCL() != expected.HCL {
				t.Fatalf("example value %v, expected %q", value, expected.HCL)
			}
		})
	}
}

func TestExampleConfigOverridesGlobalAndDerivedValuesByScenario(t *testing.T) {
	t.Parallel()

	property := &Property{
		PropertyName:  "target",
		AttributeName: "target",
		Required:      true,
		ValueType:     String,
		TestValues: &TestValues{
			Default: newTestValueScenario([]TestValueEntry{{ConfigValue: "target_1", ConfigInclude: true}}),
			Create:  newTestValueScenario([]TestValueEntry{{ConfigValue: "target_1", ConfigInclude: true}}),
		},
		globalDefinition: GlobalMetaDefinition{PropertyExampleValueOverrides: map[string]string{"target": "global_value"}},
		propertyDefinition: PropertyDefinition{ExampleConfig: PropertyExampleConfigDefinition{
			Minimum: ExampleValueScenarioDefinition{Defined: true, Entries: []ExampleValueEntryDefinition{
				{Reference: &ExampleReferenceDefinition{ClassName: "fvAEPg"}},
			}},
			Full: ExampleValueScenarioDefinition{Defined: true, Entries: []ExampleValueEntryDefinition{
				{Value: "full_value"},
			}},
		}},
	}
	ds := &DataStore{
		ClassDefinitions: map[string]ClassDefinition{
			"fvAEPg": {ResourceName: "application_epg"},
		},
		ctx: NewContext(),
	}

	minimumValue := (&Class{}).propertyExampleValue(ds, property, ExampleContext{}, MinimumExample)
	if actual := minimumValue.HCL(); actual != "aci_application_epg.example.id" {
		t.Fatalf("minimum override %q", actual)
	}
	fullValue := (&Class{}).propertyExampleValue(ds, property, ExampleContext{}, FullExample)
	if actual := fullValue.HCL(); actual != `"full_value"` {
		t.Fatalf("explicit full override %q", actual)
	}

	property.propertyDefinition.ExampleConfig.Full = ExampleValueScenarioDefinition{}
	fullFallbackValue := (&Class{}).propertyExampleValue(ds, property, ExampleContext{}, FullExample)
	if actual := fullFallbackValue.HCL(); actual != "aci_application_epg.example.id" {
		t.Fatalf("full fallback to minimum override %q", actual)
	}

	property.propertyDefinition.ExampleConfig.Minimum = ExampleValueScenarioDefinition{}
	globalFallbackValue := (&Class{}).propertyExampleValue(ds, property, ExampleContext{}, FullExample)
	if actual := globalFallbackValue.HCL(); actual != `"global_value"` {
		t.Fatalf("global fallback %q", actual)
	}
}

func TestExampleConfigReferenceSupportsCompatibilityResourceName(t *testing.T) {
	t.Parallel()

	property := &Property{
		PropertyName:  "tDn",
		AttributeName: "target_dn",
		Required:      true,
		ValueType:     String,
		propertyDefinition: PropertyDefinition{ExampleConfig: PropertyExampleConfigDefinition{
			Minimum: ExampleValueScenarioDefinition{Defined: true, Entries: []ExampleValueEntryDefinition{{
				Reference: &ExampleReferenceDefinition{
					ClassName:    "l3extInstP",
					ResourceName: "external_network_instance_profile",
				},
			}}},
		}},
	}
	ds := &DataStore{ctx: NewContext()}

	value := (&Class{Name: testClassName("l3extRsLblToInstP")}).propertyExampleValue(ds, property, ExampleContext{}, MinimumExample)
	if actual := value.HCL(); actual != "aci_external_network_instance_profile.example.id" {
		t.Fatalf("compatibility reference %q", actual)
	}
}

func TestReferenceTestValueRequiresSemanticExampleData(t *testing.T) {
	t.Parallel()

	property := &Property{PropertyName: "scope", AttributeName: "scope", Required: true, ValueType: String}
	scenario := newTestValueScenario([]TestValueEntry{{
		ConfigValue:   "aci_bridge_domain.test.id",
		ConfigInclude: true,
		ValueType:     ReferenceValue,
	}})
	ds := &DataStore{ctx: NewContext()}

	property.exampleValueFromTestScenario(ds, testClassName("fvTrackMember"), scenario)
	if err := ds.ctx.Diagnostics.Error(); err == nil {
		t.Fatal("reference test value without semantic example data should produce a diagnostic")
	}
}
