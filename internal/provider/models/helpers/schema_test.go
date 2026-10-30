package helpers

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	resourcedefaults "github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type parentDNAttributeExpected struct {
	Required       bool
	Optional       bool
	Computed       bool
	DefaultValue   string
	ValidatorCount int
}

type propertySchemaExpected struct {
	Type              PropertyType
	Required          bool
	Optional          bool
	Computed          bool
	PlanModifierCount int
	ValidatorCount    int
}

type resourcePropertySchemaInput struct {
	Property               PropertySchema
	IncludeRequiresReplace bool
}

func TestPropertySchemaResourceAttribute(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name: "test_top_level_replace_string",
			Input: resourcePropertySchemaInput{
				Property: PropertySchema{
					Type:            StringProperty,
					Required:        true,
					RequiresReplace: true,
					ValidValues:     []string{"one"},
				},
				IncludeRequiresReplace: true,
			},
			Expected: propertySchemaExpected{Type: StringProperty, Required: true, PlanModifierCount: 2, ValidatorCount: 1},
		},
		{
			Name: "test_child_replace_string",
			Input: resourcePropertySchemaInput{
				Property: PropertySchema{
					Type:            StringProperty,
					Required:        true,
					RequiresReplace: true,
				},
			},
			Expected: propertySchemaExpected{Type: StringProperty, Required: true, PlanModifierCount: 1},
		},
		{
			Name: "test_set",
			Input: resourcePropertySchemaInput{
				Property: PropertySchema{
					Type:           SetProperty,
					Optional:       true,
					Computed:       true,
					SetElementType: types.StringType,
					ValidValues:    []string{"one", "two"},
				},
			},
			Expected: propertySchemaExpected{Type: SetProperty, Optional: true, Computed: true, PlanModifierCount: 1, ValidatorCount: 2},
		},
		{
			Name: "test_read_only_object",
			Input: resourcePropertySchemaInput{
				Property: PropertySchema{
					Type:                 ObjectProperty,
					ReadOnly:             true,
					Computed:             true,
					ObjectAttributeTypes: map[string]attr.Type{},
				},
			},
			Expected: propertySchemaExpected{Type: ObjectProperty, Computed: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(resourcePropertySchemaInput)
			expected := testCase.Expected.(propertySchemaExpected)
			attribute := input.Property.ResourceAttribute(input.IncludeRequiresReplace)

			switch expected.Type {
			case SetProperty:
				actual, ok := attribute.(resourceschema.SetAttribute)
				if !ok {
					t.Fatalf("expected set attribute, got %T", attribute)
				}
				if actual.Required != expected.Required || actual.Optional != expected.Optional || actual.Computed != expected.Computed || len(actual.PlanModifiers) != expected.PlanModifierCount || len(actual.Validators) != expected.ValidatorCount {
					t.Fatalf("unexpected set attribute: %#v", actual)
				}
			case ObjectProperty:
				actual, ok := attribute.(resourceschema.ObjectAttribute)
				if !ok {
					t.Fatalf("expected object attribute, got %T", attribute)
				}
				if actual.Required != expected.Required || actual.Optional != expected.Optional || actual.Computed != expected.Computed || len(actual.PlanModifiers) != expected.PlanModifierCount {
					t.Fatalf("unexpected object attribute: %#v", actual)
				}
			default:
				actual, ok := attribute.(resourceschema.StringAttribute)
				if !ok {
					t.Fatalf("expected string attribute, got %T", attribute)
				}
				if actual.Required != expected.Required || actual.Optional != expected.Optional || actual.Computed != expected.Computed || len(actual.PlanModifiers) != expected.PlanModifierCount || len(actual.Validators) != expected.ValidatorCount {
					t.Fatalf("unexpected string attribute: %#v", actual)
				}
			}
		})
	}
}

func TestPropertySchemaDataSourceAttribute(t *testing.T) {
	property := PropertySchema{
		Type:        StringProperty,
		Sensitive:   true,
		ValidValues: []string{"one"},
		Description: "Example property.",
	}

	computed := property.DataSourceAttribute(false).(datasourceschema.StringAttribute)
	if !computed.Computed || computed.Required || !computed.Sensitive || len(computed.Validators) != 0 {
		t.Fatalf("unexpected computed data-source attribute: %#v", computed)
	}

	required := property.DataSourceAttribute(true).(datasourceschema.StringAttribute)
	if !required.Required || required.Computed || !required.Sensitive || len(required.Validators) != 1 {
		t.Fatalf("unexpected required data-source attribute: %#v", required)
	}
}

func childClassSchemaForTest() ClassSchema {
	return ClassSchema{
		Properties: map[string]PropertySchema{
			"key": {
				Type:            StringProperty,
				Required:        true,
				RequiresReplace: true,
			},
		},
		Children: map[string]ChildSchema{},
	}
}

type requiredPropertyContextInput struct {
	ReadOnlyRoot bool
}

type requiredPropertyContextExpected struct {
	Required       bool
	Optional       bool
	Computed       bool
	ValidatorCount int
}

func TestTopLevelRequiredPropertyContext(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_without_read_only_properties",
			Input:    requiredPropertyContextInput{},
			Expected: requiredPropertyContextExpected{Required: true, ValidatorCount: 1},
		},
		{
			Name:     "test_with_read_only_property",
			Input:    requiredPropertyContextInput{ReadOnlyRoot: true},
			Expected: requiredPropertyContextExpected{Optional: true, Computed: true, ValidatorCount: 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(requiredPropertyContextInput)
			expected := testCase.Expected.(requiredPropertyContextExpected)
			properties := map[string]PropertySchema{
				"name": {
					Type:            StringProperty,
					Required:        true,
					RequiresReplace: true,
					ValidValues:     []string{"example"},
				},
			}
			if input.ReadOnlyRoot {
				properties["read_only"] = PropertySchema{Type: StringProperty, Computed: true, ReadOnly: true}
			}
			definition := ClassSchema{Properties: properties, Children: map[string]ChildSchema{}}

			resourceSchema := definition.ResourceSchema("")
			if diagnostics := resourceSchema.ValidateImplementation(context.Background()); diagnostics.HasError() {
				t.Fatalf("invalid resource schema: %v", diagnostics)
			}
			attribute := resourceSchema.Attributes["name"].(resourceschema.StringAttribute)

			if attribute.Required != expected.Required || attribute.Optional != expected.Optional || attribute.Computed != expected.Computed {
				t.Fatalf("unexpected required property flags: required=%t optional=%t computed=%t", attribute.Required, attribute.Optional, attribute.Computed)
			}
			if len(attribute.Validators) != expected.ValidatorCount {
				t.Fatalf("validator count %d, expected %d", len(attribute.Validators), expected.ValidatorCount)
			}
			if len(attribute.PlanModifiers) != 2 {
				t.Fatalf("plan modifier count %d, expected 2", len(attribute.PlanModifiers))
			}
		})
	}
}

type singleChildContextInput struct {
	Required bool
}

func TestSingleChildRequiredPropertyContext(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_optional_single_child",
			Input:    singleChildContextInput{},
			Expected: requiredPropertyContextExpected{Optional: true, Computed: true, ValidatorCount: 1},
		},
		{
			Name:     "test_required_single_child",
			Input:    singleChildContextInput{Required: true},
			Expected: requiredPropertyContextExpected{Required: true, ValidatorCount: 1},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(singleChildContextInput)
			expected := testCase.Expected.(requiredPropertyContextExpected)
			definition := ClassSchema{
				Properties: map[string]PropertySchema{},
				Children: map[string]ChildSchema{
					"child": {
						Cardinality: SingleChild,
						Required:    input.Required,
						Schema: func() ClassSchema {
							return ClassSchema{
								Properties: map[string]PropertySchema{
									"name": {Type: StringProperty, Required: true, ValidValues: []string{"example"}},
								},
								Children: map[string]ChildSchema{},
							}
						},
					},
				},
			}

			resourceSchema := definition.ResourceSchema("")
			if diagnostics := resourceSchema.ValidateImplementation(context.Background()); diagnostics.HasError() {
				t.Fatalf("invalid resource schema: %v", diagnostics)
			}
			attribute := resourceSchema.Attributes["child"].(resourceschema.SingleNestedAttribute)
			name := attribute.Attributes["name"].(resourceschema.StringAttribute)

			expectedOptionalAndComputed := !input.Required
			if attribute.Required != input.Required || attribute.Optional != expectedOptionalAndComputed || attribute.Computed != expectedOptionalAndComputed {
				t.Fatalf("unexpected single child flags: required=%t optional=%t computed=%t", attribute.Required, attribute.Optional, attribute.Computed)
			}
			expectedObjectValidators := 1
			if input.Required {
				expectedObjectValidators = 0
			}
			if len(attribute.Validators) != expectedObjectValidators {
				t.Fatalf("object validator count %d, expected %d", len(attribute.Validators), expectedObjectValidators)
			}
			if len(attribute.PlanModifiers) != 1 {
				t.Fatalf("singleton child plan modifier count %d, expected 1", len(attribute.PlanModifiers))
			}
			if name.Required != expected.Required || name.Optional != expected.Optional || name.Computed != expected.Computed {
				t.Fatalf("unexpected child property flags: required=%t optional=%t computed=%t", name.Required, name.Optional, name.Computed)
			}
			if len(name.Validators) != expected.ValidatorCount {
				t.Fatalf("property validator count %d, expected %d", len(name.Validators), expected.ValidatorCount)
			}
			if len(name.PlanModifiers) != 0 {
				t.Fatalf("nested property plan modifier count %d, expected 0", len(name.PlanModifiers))
			}
		})
	}
}

func TestSetChildRequiredPropertyContext(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_without_read_only_root_property",
			Input:    requiredPropertyContextInput{},
			Expected: requiredPropertyContextExpected{Required: true, ValidatorCount: 1},
		},
		{
			Name:     "test_with_read_only_root_property",
			Input:    requiredPropertyContextInput{ReadOnlyRoot: true},
			Expected: requiredPropertyContextExpected{Optional: true, Computed: true, ValidatorCount: 2},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input := testCase.Input.(requiredPropertyContextInput)
			expected := testCase.Expected.(requiredPropertyContextExpected)
			properties := map[string]PropertySchema{}
			if input.ReadOnlyRoot {
				properties["read_only"] = PropertySchema{Type: StringProperty, Computed: true, ReadOnly: true}
			}
			definition := ClassSchema{
				Properties: properties,
				Children: map[string]ChildSchema{
					"children": {
						Cardinality: SetChild,
						Schema: func() ClassSchema {
							return ClassSchema{
								Properties: map[string]PropertySchema{
									"name": {Type: StringProperty, Required: true, ValidValues: []string{"example"}},
								},
								Children: map[string]ChildSchema{},
							}
						},
					},
				},
			}

			resourceSchema := definition.ResourceSchema("")
			if diagnostics := resourceSchema.ValidateImplementation(context.Background()); diagnostics.HasError() {
				t.Fatalf("invalid resource schema: %v", diagnostics)
			}
			attribute := resourceSchema.Attributes["children"].(resourceschema.SetNestedAttribute)
			name := attribute.NestedObject.Attributes["name"].(resourceschema.StringAttribute)
			if len(attribute.PlanModifiers) != 1 {
				t.Fatalf("set child plan modifier count %d, expected 1", len(attribute.PlanModifiers))
			}

			if name.Required != expected.Required || name.Optional != expected.Optional || name.Computed != expected.Computed {
				t.Fatalf("unexpected child property flags: required=%t optional=%t computed=%t", name.Required, name.Optional, name.Computed)
			}
			if len(name.Validators) != expected.ValidatorCount {
				t.Fatalf("property validator count %d, expected %d", len(name.Validators), expected.ValidatorCount)
			}
			if len(name.PlanModifiers) != 0 {
				t.Fatalf("nested property plan modifier count %d, expected 0", len(name.PlanModifiers))
			}
		})
	}
}

func TestReconciledChildOwnsDescendantStateHandling(t *testing.T) {
	grandChildSchema := func() ClassSchema {
		return ClassSchema{
			Properties: map[string]PropertySchema{
				"name": {Type: StringProperty, Required: true},
			},
			Children:    map[string]ChildSchema{},
			Identifiers: []string{"name"},
		}
	}
	childSchema := func() ClassSchema {
		return ClassSchema{
			Properties: map[string]PropertySchema{
				"name": {Type: StringProperty, Required: true},
			},
			Children: map[string]ChildSchema{
				"grand_children": {
					Cardinality: SetChild,
					Schema:      grandChildSchema,
				},
			},
			Identifiers: []string{"name"},
		}
	}
	definition := ClassSchema{
		Properties: map[string]PropertySchema{},
		Children: map[string]ChildSchema{
			"child": {
				Cardinality: SingleChild,
				Schema:      childSchema,
			},
		},
	}

	resourceSchema := definition.ResourceSchema("")
	child := resourceSchema.Attributes["child"].(resourceschema.SingleNestedAttribute)
	if len(child.PlanModifiers) != 1 {
		t.Fatalf("outer child plan modifier count %d, expected 1", len(child.PlanModifiers))
	}
	childName := child.Attributes["name"].(resourceschema.StringAttribute)
	if len(childName.PlanModifiers) != 0 {
		t.Fatalf("child property plan modifier count %d, expected 0", len(childName.PlanModifiers))
	}
	grandChildren := child.Attributes["grand_children"].(resourceschema.SetNestedAttribute)
	if len(grandChildren.PlanModifiers) != 0 {
		t.Fatalf("descendant child plan modifier count %d, expected 0", len(grandChildren.PlanModifiers))
	}
	grandChildName := grandChildren.NestedObject.Attributes["name"].(resourceschema.StringAttribute)
	if len(grandChildName.PlanModifiers) != 0 {
		t.Fatalf("grandchild property plan modifier count %d, expected 0", len(grandChildName.PlanModifiers))
	}
}

func TestClassSchema(t *testing.T) {
	definition := ClassSchema{
		ResourceDescription:   "Example resource.",
		DataSourceDescription: "Example data source.",
		Properties: map[string]PropertySchema{
			"annotation": {
				Type:     StringProperty,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:            StringProperty,
				Required:        true,
				RequiresReplace: true,
				ValidValues:     []string{"example"},
			},
		},
		Children: map[string]ChildSchema{
			"children": {
				Cardinality: SetChild,
				Schema:      childClassSchemaForTest,
			},
		},
		Identifiers: []string{"name"},
		ParentDN: &ParentDNSchema{
			AttributeName: "parent_dn",
			Config:        ParentDNAttributeConfig{DefaultValue: "uni"},
		},
		AnnotationAttribute: "annotation",
	}

	resourceSchema := definition.ResourceSchema("orchestrator:terraform")
	if diagnostics := resourceSchema.ValidateImplementation(context.Background()); diagnostics.HasError() {
		t.Fatalf("invalid resource schema: %v", diagnostics)
	}
	name := resourceSchema.Attributes["name"].(resourceschema.StringAttribute)
	if len(name.PlanModifiers) != 2 {
		t.Fatalf("top-level name plan modifier count %d, expected 2", len(name.PlanModifiers))
	}
	children := resourceSchema.Attributes["children"].(resourceschema.SetNestedAttribute)
	if len(children.PlanModifiers) != 1 {
		t.Fatalf("set child plan modifier count %d, expected 1", len(children.PlanModifiers))
	}
	childKey := children.NestedObject.Attributes["key"].(resourceschema.StringAttribute)
	if len(childKey.PlanModifiers) != 0 {
		t.Fatalf("child key plan modifier count %d, expected 0", len(childKey.PlanModifiers))
	}
	if _, ok := resourceSchema.Attributes["id"]; !ok {
		t.Fatal("resource schema does not contain id")
	}
	if parentDN := resourceSchema.Attributes["parent_dn"].(resourceschema.StringAttribute); !parentDN.Optional || !parentDN.Computed {
		t.Fatalf("unexpected resource parent_dn attribute: %#v", parentDN)
	}

	dataSourceSchema := definition.DataSourceSchema()
	if diagnostics := dataSourceSchema.ValidateImplementation(context.Background()); diagnostics.HasError() {
		t.Fatalf("invalid data-source schema: %v", diagnostics)
	}
	identifier := dataSourceSchema.Attributes["name"].(datasourceschema.StringAttribute)
	if !identifier.Required || identifier.Computed || len(identifier.Validators) != 1 {
		t.Fatalf("unexpected data-source identifier: %#v", identifier)
	}
	dataSourceChildren := dataSourceSchema.Attributes["children"].(datasourceschema.SetNestedAttribute)
	if !dataSourceChildren.Computed {
		t.Fatalf("unexpected data-source child: %#v", dataSourceChildren)
	}
}

func TestResourceParentDNAttribute(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_required_without_default",
			Input:    ParentDNAttributeConfig{},
			Expected: parentDNAttributeExpected{Required: true},
		},
		{
			Name:  "test_optional_computed_with_default_and_valid_values",
			Input: ParentDNAttributeConfig{DefaultValue: "uni/infra", ValidValues: []string{"uni/infra", "uni/fabric"}},
			Expected: parentDNAttributeExpected{
				Optional:       true,
				Computed:       true,
				DefaultValue:   "uni/infra",
				ValidatorCount: 1,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			attribute := ResourceParentDNAttribute(testCase.Input.(ParentDNAttributeConfig))
			expected := testCase.Expected.(parentDNAttributeExpected)

			if attribute.Required != expected.Required || attribute.Optional != expected.Optional || attribute.Computed != expected.Computed {
				t.Fatalf("unexpected resource parent_dn flags: required=%t optional=%t computed=%t", attribute.Required, attribute.Optional, attribute.Computed)
			}
			if len(attribute.Validators) != expected.ValidatorCount {
				t.Fatalf("resource parent_dn validator count %d, expected %d", len(attribute.Validators), expected.ValidatorCount)
			}
			if len(attribute.PlanModifiers) != 2 {
				t.Fatalf("resource parent_dn plan modifier count %d, expected 2", len(attribute.PlanModifiers))
			}

			actualDefault := ""
			if attribute.Default != nil {
				var response resourcedefaults.StringResponse
				attribute.Default.DefaultString(context.Background(), resourcedefaults.StringRequest{}, &response)
				actualDefault = response.PlanValue.ValueString()
			}
			if actualDefault != expected.DefaultValue {
				t.Fatalf("resource parent_dn default %q, expected %q", actualDefault, expected.DefaultValue)
			}
		})
	}
}

func TestDataSourceParentDNAttribute(t *testing.T) {
	testCases := []test.TestCase{
		{
			Name:     "test_required_without_default",
			Input:    ParentDNAttributeConfig{},
			Expected: parentDNAttributeExpected{Required: true},
		},
		{
			Name:  "test_optional_with_known_default_and_valid_values",
			Input: ParentDNAttributeConfig{DefaultValue: "uni/infra", ValidValues: []string{"uni/infra"}},
			Expected: parentDNAttributeExpected{
				Optional:       true,
				ValidatorCount: 1,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			attribute := DataSourceParentDNAttribute(testCase.Input.(ParentDNAttributeConfig))
			expected := testCase.Expected.(parentDNAttributeExpected)

			if attribute.Required != expected.Required || attribute.Optional != expected.Optional || attribute.Computed != expected.Computed {
				t.Fatalf("unexpected data-source parent_dn flags: required=%t optional=%t computed=%t", attribute.Required, attribute.Optional, attribute.Computed)
			}
			if len(attribute.Validators) != expected.ValidatorCount {
				t.Fatalf("data-source parent_dn validator count %d, expected %d", len(attribute.Validators), expected.ValidatorCount)
			}
		})
	}
}

func TestCommonSchemaAttributes(t *testing.T) {
	resourceID := ResourceIDAttribute()
	if !resourceID.Computed || len(resourceID.PlanModifiers) != 1 {
		t.Fatalf("unexpected resource ID attribute: %#v", resourceID)
	}

	dataSourceID := DataSourceIDAttribute()
	if !dataSourceID.Computed {
		t.Fatalf("unexpected data-source ID attribute: %#v", dataSourceID)
	}

	annotation := ResourceAnnotationAttribute("orchestrator:terraform", "Annotation description.")
	if !annotation.Optional || !annotation.Computed || len(annotation.PlanModifiers) != 1 {
		t.Fatalf("unexpected annotation attribute: %#v", annotation)
	}
	var defaultResponse resourcedefaults.StringResponse
	annotation.Default.DefaultString(context.Background(), resourcedefaults.StringRequest{}, &defaultResponse)
	if !defaultResponse.PlanValue.Equal(types.StringValue("orchestrator:terraform")) {
		t.Fatalf("annotation default %q, expected %q", defaultResponse.PlanValue.ValueString(), "orchestrator:terraform")
	}
}

func TestParentDNValidValues(t *testing.T) {
	attribute := ResourceParentDNAttribute(ParentDNAttributeConfig{ValidValues: []string{"uni/infra"}})
	request := validator.StringRequest{ConfigValue: types.StringValue("uni/fabric")}
	var response validator.StringResponse
	attribute.Validators[0].ValidateString(context.Background(), request, &response)
	if !response.Diagnostics.HasError() {
		t.Fatal("expected an unsupported parent_dn value to produce a diagnostic")
	}
}

func TestResourcePropertyPlanModifiers(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "test_string_update", Input: false, Expected: 1},
		{Name: "test_string_replace", Input: true, Expected: 2},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := len(ResourceStringPlanModifiers(testCase.Input.(bool)))
			if actual != testCase.Expected.(int) {
				t.Fatalf("string plan modifier count %d, expected %d", actual, testCase.Expected.(int))
			}
		})
	}

	setTestCases := []test.TestCase{
		{Name: "test_set_update", Input: false, Expected: 1},
		{Name: "test_set_replace", Input: true, Expected: 2},
	}
	for _, testCase := range setTestCases {
		t.Run(testCase.Name, func(t *testing.T) {
			actual := len(ResourceSetPlanModifiers(testCase.Input.(bool)))
			if actual != testCase.Expected.(int) {
				t.Fatalf("set plan modifier count %d, expected %d", actual, testCase.Expected.(int))
			}
		})
	}

	stringModifiers := ResourceStringPlanModifiers(true)
	if reflect.TypeOf(stringModifiers[len(stringModifiers)-1]) != reflect.TypeOf(stringplanmodifier.RequiresReplace()) {
		t.Fatal("expected string RequiresReplace to remain the final plan modifier")
	}
	setModifiers := ResourceSetPlanModifiers(true)
	if reflect.TypeOf(setModifiers[len(setModifiers)-1]) != reflect.TypeOf(setplanmodifier.RequiresReplace()) {
		t.Fatal("expected set RequiresReplace to remain the final plan modifier")
	}
}
