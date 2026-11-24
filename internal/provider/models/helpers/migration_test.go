package helpers

import (
	"context"
	"maps"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/container"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
	"github.com/CiscoDevNet/terraform-provider-aci/v2/internal/customtypes"
)

func TestLegacyRelationName(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "dn", Input: types.StringValue("uni/tn-common/monepg-policy-one"), Expected: types.StringValue("policy-one")},
		{Name: "name", Input: types.StringValue("policy-one"), Expected: types.StringValue("policy-one")},
		{Name: "empty", Input: types.StringValue(""), Expected: types.StringValue("")},
		{Name: "null", Input: types.StringNull(), Expected: types.StringNull()},
		{Name: "unknown", Input: types.StringUnknown(), Expected: types.StringUnknown()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			assert.Equal(t, testCase.Expected, LegacyRelationName(context.Background(), &diagnostics, testCase.Input.(types.String)))
			assert.False(t, diagnostics.HasError())
		})
	}
}

func TestLegacyChildObject(t *testing.T) {
	priorityType := customtypes.NewSemanticStringType(map[string]string{"unspecified": "level0"})
	attributeTypes := map[string]attr.Type{
		"name": types.StringType, "annotation": types.StringType, "priority": priorityType,
		"tags": types.SetType{ElemType: types.StringType}, "settings": types.ObjectType{AttrTypes: map[string]attr.Type{"enabled": types.BoolType}},
	}
	testCases := []test.TestCase{
		{Name: "known_value", Input: map[string]attr.Value{"name": types.StringValue("policy")}, Expected: map[string]attr.Value{"name": types.StringValue("policy")}},
		{Name: "null_value", Input: map[string]attr.Value{"name": types.StringNull()}, Expected: map[string]attr.Value{}},
		{Name: "unknown_value", Input: map[string]attr.Value{"name": types.StringUnknown()}, Expected: map[string]attr.Value{"name": types.StringUnknown()}},
		{Name: "multiple_values", Input: map[string]attr.Value{"name": types.StringValue("policy"), "annotation": types.StringValue("kept")}, Expected: map[string]attr.Value{"name": types.StringValue("policy"), "annotation": types.StringValue("kept")}},
		{Name: "no_values", Input: map[string]attr.Value(nil), Expected: map[string]attr.Value{}},
		{Name: "wrong_type", Input: map[string]attr.Value{"name": types.BoolValue(true)}},
		{Name: "unknown_attribute", Input: map[string]attr.Value{"other": types.StringValue("policy")}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(map[string]attr.Value)
			originalLength := len(values)
			var diagnostics diag.Diagnostics
			object := LegacyChildObject(context.Background(), &diagnostics, attributeTypes, values)
			assert.Len(t, values, originalLength, "construction must not add attributes to its input")
			if testCase.Expected == nil {
				assert.True(t, diagnostics.HasError())
				return
			}
			if !assert.False(t, diagnostics.HasError(), "%v", diagnostics) {
				return
			}
			expected := map[string]attr.Value{
				"name": types.StringNull(), "annotation": types.StringNull(), "priority": priorityType.Null(),
				"tags": types.SetNull(types.StringType), "settings": types.ObjectNull(map[string]attr.Type{"enabled": types.BoolType}),
			}
			maps.Copy(expected, testCase.Expected.(map[string]attr.Value))
			assert.Equal(t, expected, object.Attributes())
		})
	}
}

func TestLegacyChildAttributes(t *testing.T) {
	attributeTypes := map[string]attr.Type{"name": types.StringType, "annotation": types.StringType}
	prior := types.ObjectValueMust(attributeTypes, map[string]attr.Value{"name": types.StringValue("old"), "annotation": types.StringValue("retained")})
	testCases := []test.TestCase{
		{Name: "planned_values", Input: prior, Expected: types.StringValue("retained")},
		{Name: "unknown_uses_state", Input: types.ObjectUnknown(attributeTypes), Expected: types.StringValue("retained")},
		{Name: "unknown_attributes_use_state", Input: types.ObjectValueMust(attributeTypes, map[string]attr.Value{"name": types.StringUnknown(), "annotation": types.StringUnknown()}), Expected: types.StringValue("retained")},
		{Name: "null_uses_state", Input: types.ObjectNull(attributeTypes), Expected: types.StringValue("retained")},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			attributes := LegacyChildAttributes(context.Background(), &diagnostics, testCase.Input.(types.Object), &prior, attributeTypes)
			assert.Equal(t, testCase.Expected, attributes["annotation"])
			attributes["name"] = types.StringValue("changed")
			assert.Equal(t, types.StringValue("old"), prior.Attributes()["name"])
			assert.False(t, diagnostics.HasError())
		})
	}
	var diagnostics diag.Diagnostics
	attributes := LegacyChildAttributes(context.Background(), &diagnostics, types.ObjectNull(attributeTypes), nil, attributeTypes)
	assert.Equal(t, types.StringUnknown(), attributes["annotation"])
	empty := types.ObjectValueMust(attributeTypes, map[string]attr.Value{"name": types.StringNull(), "annotation": types.StringNull()})
	attributes = LegacyChildAttributes(context.Background(), &diagnostics, empty, &empty, attributeTypes)
	assert.Equal(t, types.StringUnknown(), attributes["annotation"])

	assert.False(t, diagnostics.HasError())
}

func TestLegacyChildRelationPlan(t *testing.T) {
	attributeTypes := map[string]attr.Type{"name": types.StringType, "annotation": types.StringType, "description": types.StringType}
	object := func(name, annotation, description types.String) types.Object {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{"name": name, "annotation": annotation, "description": description})
	}
	nullChild := types.ObjectNull(attributeTypes)
	prior := object(types.StringValue("previous"), types.StringValue("prior-annotation"), types.StringValue("prior-description"))
	planned := object(types.StringValue("planned"), types.StringValue("planned-annotation"), types.StringValue("planned-description"))
	annotation := types.StringValue("parent-annotation")
	legacy := types.StringValue("uni/tn-other/monepg-policy-one")
	type input struct {
		legacy     types.String
		planned    types.Object
		fallback   *types.Object
		annotation attr.Value
		named      bool
	}
	testCases := []test.TestCase{
		{
			Name:     "named_relation",
			Input:    input{legacy: legacy, planned: nullChild, annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), annotation, types.StringUnknown()),
		},
		{
			Name:     "explicit_target_dn",
			Input:    input{legacy: legacy, planned: nullChild, annotation: annotation},
			Expected: object(legacy, annotation, types.StringUnknown()),
		},
		{
			Name:     "unknown_relation",
			Input:    input{legacy: types.StringUnknown(), planned: nullChild, annotation: annotation, named: true},
			Expected: object(types.StringUnknown(), annotation, types.StringUnknown()),
		},
		{
			Name:     "null_relation",
			Input:    input{legacy: types.StringNull(), planned: nullChild, annotation: annotation, named: true},
			Expected: object(types.StringNull(), annotation, types.StringUnknown()),
		},
		{
			Name:     "retains_planned_values",
			Input:    input{legacy: legacy, planned: planned, fallback: &prior, annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), types.StringValue("planned-annotation"), types.StringValue("planned-description")),
		},
		{
			Name:     "retains_prior_values",
			Input:    input{legacy: legacy, planned: nullChild, fallback: &prior, annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), types.StringValue("prior-annotation"), types.StringValue("prior-description")),
		},
		{
			Name:     "fills_unknown_values_from_state",
			Input:    input{legacy: legacy, planned: object(types.StringUnknown(), types.StringUnknown(), types.StringUnknown()), fallback: &prior, annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), types.StringValue("prior-annotation"), types.StringValue("prior-description")),
		},
		{
			Name:     "inherits_unknown_annotation",
			Input:    input{legacy: legacy, planned: object(types.StringValue("planned"), types.StringUnknown(), types.StringValue("planned-description")), annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), annotation, types.StringValue("planned-description")),
		},
		{
			Name:     "retains_null_annotation",
			Input:    input{legacy: legacy, planned: object(types.StringValue("planned"), types.StringNull(), types.StringValue("planned-description")), annotation: annotation, named: true},
			Expected: object(types.StringValue("policy-one"), types.StringNull(), types.StringValue("planned-description")),
		},
		{
			Name:     "without_parent_annotation",
			Input:    input{legacy: legacy, planned: nullChild, named: true},
			Expected: object(types.StringValue("policy-one"), types.StringUnknown(), types.StringUnknown()),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(input)
			plannedBefore := values.planned.Attributes()
			var diagnostics diag.Diagnostics
			relation := LegacyChildRelation{AttributeTypes: attributeTypes, AttributeName: "name", Named: values.named}
			result := relation.Plan(context.Background(), &diagnostics, values.legacy, values.planned, values.fallback, values.annotation)
			if !assert.False(t, diagnostics.HasError(), "%v", diagnostics) {
				return
			}
			assert.True(t, result.Equal(testCase.Expected.(types.Object)), "%s", result)
			assert.Equal(t, plannedBefore, values.planned.Attributes(), "planning must not mutate the planned child")
			assert.Equal(t, types.StringValue("previous"), prior.Attributes()["name"], "planning must not mutate prior state")
		})
	}
}

func TestLegacyTargetDN(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "resolved_in_another_tenant", Input: `{"imdata":[{"fvTenant":{"attributes":{"dn":"uni/tn-example"},"children":[{"fvRsTenantMonPol":{"attributes":{"tnMonEPGPolName":"default","tDn":"uni/tn-common/monepg-default"}}}]}}]}`, Expected: types.StringValue("uni/tn-common/monepg-default")},
		{Name: "missing_relation", Input: `{"imdata":[{"fvTenant":{"attributes":{"dn":"uni/tn-example"}}}]}`, Expected: types.StringNull()},
		{Name: "missing_target", Input: `{"imdata":[{"fvTenant":{"children":[{"fvRsTenantMonPol":{"attributes":{}}}]}}]}`, Expected: types.StringNull()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			response, err := container.ParseJSON([]byte(testCase.Input.(string)))
			if !assert.NoError(t, err) {
				return
			}
			var diagnostics diag.Diagnostics
			assert.Equal(t, testCase.Expected, LegacyTargetDN(context.Background(), &diagnostics, response, "fvTenant", "fvRsTenantMonPol"))
			assert.False(t, diagnostics.HasError())
		})
	}
}

func TestLegacyAttributeSchema(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "deprecated_string", Input: LegacyAttributeSchema{Type: types.StringType, Replacement: "relation_to_monitoring_policy"}},
		{Name: "deprecated_ip_address", Input: LegacyAttributeSchema{Type: customtypes.IPAddressType{}, Replacement: "link_local_ipv6_address"}},
		{Name: "sensitive_string", Input: LegacyAttributeSchema{Type: types.StringType, Replacement: "password", Sensitive: true}},
		{Name: "deprecated_list", Input: LegacyAttributeSchema{Type: types.ListType{ElemType: types.StringType}, Replacement: "interface_controls"}},
		{Name: "frozen_set", Input: LegacyAttributeSchema{Type: types.SetType{ElemType: types.StringType}, Frozen: true}},
		{Name: "frozen_string", Input: LegacyAttributeSchema{Type: types.StringType, Frozen: true}},
		{Name: "deprecated_set", Input: LegacyAttributeSchema{Type: types.SetType{ElemType: types.StringType}, Replacement: "relations"}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			legacy := testCase.Input.(LegacyAttributeSchema)
			resource := legacy.ResourceAttribute()
			assert.True(t, resource.IsOptional())
			assert.True(t, resource.IsComputed())
			assert.Equal(t, legacy.Sensitive, resource.IsSensitive())
			assert.True(t, resource.GetType().Equal(legacy.Type))
			if attribute, ok := resource.(resourceschema.StringAttribute); ok {
				if legacy.Frozen {
					assert.Len(t, attribute.PlanModifiers, 1, "frozen strings must preserve prior state for unknown plans")
				} else {
					assert.Empty(t, attribute.PlanModifiers, "deprecated outputs must remain unknown during real changes")
				}
			}
			if attribute, ok := resource.(resourceschema.SetAttribute); ok {
				if legacy.Frozen {
					assert.Len(t, attribute.PlanModifiers, 1)
				} else {
					assert.Empty(t, attribute.PlanModifiers, "deprecated outputs must remain unknown during real changes")
				}
			}
			datasource := legacy.DataSourceAttribute()
			assert.False(t, datasource.IsOptional())
			assert.True(t, datasource.IsComputed())
			assert.Equal(t, legacy.Sensitive, datasource.IsSensitive())
			assert.True(t, datasource.GetType().Equal(legacy.Type))
		})
	}
}

func TestLegacyIPAddressValue(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "known", Input: types.StringValue("FE80::1"), Expected: customtypes.NewIPAddressValue("FE80::1")},
		{Name: "null", Input: types.StringNull(), Expected: customtypes.NewIPAddressNull()},
		{Name: "unknown", Input: types.StringUnknown(), Expected: customtypes.NewIPAddressUnknown()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			value := LegacyValue[customtypes.IPAddressValue](context.Background(), &diagnostics, testCase.Input.(types.String), customtypes.IPAddressType{})
			assert.Equal(t, testCase.Expected, value)
			assert.False(t, diagnostics.HasError(), diagnostics)
		})
	}
}

func TestLegacyValue(t *testing.T) {
	priority := customtypes.NewSemanticStringType(map[string]string{"unspecified": "level0"})
	testCases := []test.TestCase{
		{Name: "known", Input: types.StringValue("level1"), Expected: priority.Value("level1")},
		{Name: "null", Input: types.StringNull(), Expected: priority.Null()},
		{Name: "unknown", Input: types.StringUnknown(), Expected: priority.Unknown()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			converted := LegacyValue[customtypes.SemanticStringValue](context.Background(), &diagnostics, testCase.Input.(types.String), priority)
			assert.Equal(t, testCase.Expected, converted)
			assert.Equal(t, testCase.Input, LegacyValue[types.String](context.Background(), &diagnostics, converted, types.StringType))
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
		})
	}
}

func TestLegacyCollectionValues(t *testing.T) {
	type input struct {
		Value attr.Value
		Type  attr.Type
	}
	listType := types.ListType{ElemType: types.StringType}
	setType := types.SetType{ElemType: types.StringType}
	one, two := types.StringValue("one"), types.StringValue("two")
	testCases := []test.TestCase{
		{Name: "null_list", Input: input{types.ListNull(types.StringType), setType}, Expected: types.SetNull(types.StringType)},
		{Name: "unknown_list", Input: input{types.ListUnknown(types.StringType), setType}, Expected: types.SetUnknown(types.StringType)},
		{Name: "empty_list", Input: input{types.ListValueMust(types.StringType, []attr.Value{}), setType}, Expected: types.SetValueMust(types.StringType, []attr.Value{})},
		{Name: "duplicate_list_values", Input: input{types.ListValueMust(types.StringType, []attr.Value{one, two, one}), setType}, Expected: types.SetValueMust(types.StringType, []attr.Value{one, two})},
		{Name: "unknown_members_remain_distinct", Input: input{types.ListValueMust(types.StringType, []attr.Value{types.StringUnknown(), types.StringUnknown(), one, one}), setType}, Expected: types.SetValueMust(types.StringType, []attr.Value{types.StringUnknown(), types.StringUnknown(), one})},
		{Name: "null_set", Input: input{types.SetNull(types.StringType), listType}, Expected: types.ListNull(types.StringType)},
		{Name: "unknown_set", Input: input{types.SetUnknown(types.StringType), listType}, Expected: types.ListUnknown(types.StringType)},
		{Name: "known_set", Input: input{types.SetValueMust(types.StringType, []attr.Value{two, one}), listType}, Expected: types.ListValueMust(types.StringType, []attr.Value{two, one})},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(input)
			var diagnostics diag.Diagnostics
			result := LegacyValue[attr.Value](context.Background(), &diagnostics, values.Value, values.Type)
			assert.False(t, diagnostics.HasError(), diagnostics)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

func TestLegacyStringList(t *testing.T) {
	type input struct {
		Current types.Set
		Legacy  types.List
	}
	one, two := types.StringValue("one"), types.StringValue("two")
	current := types.SetValueMust(types.StringType, []attr.Value{two, one})
	null := types.ListNull(types.StringType)
	sorted := types.ListValueMust(types.StringType, []attr.Value{one, two})
	configured := types.ListValueMust(types.StringType, []attr.Value{two, one, one})
	testCases := []test.TestCase{
		{Name: "null", Input: input{types.SetNull(types.StringType), null}, Expected: null},
		{Name: "unknown", Input: input{types.SetUnknown(types.StringType), null}, Expected: types.ListUnknown(types.StringType)},
		{Name: "empty", Input: input{types.SetValueMust(types.StringType, []attr.Value{}), null}, Expected: types.ListValueMust(types.StringType, []attr.Value{})},
		{Name: "sorted_response", Input: input{current, null}, Expected: sorted},
		{Name: "preserve_configured_order_and_duplicates", Input: input{current, configured}, Expected: configured},
		{Name: "changed_response", Input: input{current, types.ListValueMust(types.StringType, []attr.Value{one})}, Expected: sorted},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			values := testCase.Input.(input)
			var diagnostics diag.Diagnostics
			result := LegacyStringList(context.Background(), &diagnostics, values.Current, values.Legacy)
			assert.False(t, diagnostics.HasError(), diagnostics)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

func TestLegacyChildFromSet(t *testing.T) {
	attributes := map[string]attr.Type{"name": types.StringType}
	element := types.ObjectType{AttrTypes: attributes}
	first := types.ObjectValueMust(attributes, map[string]attr.Value{"name": types.StringValue("first")})
	second := types.ObjectValueMust(attributes, map[string]attr.Value{"name": types.StringValue("second")})
	testCases := []test.TestCase{
		{Name: "null", Input: types.SetNull(element), Expected: types.ObjectNull(attributes)},
		{Name: "empty", Input: types.SetValueMust(element, []attr.Value{}), Expected: types.ObjectNull(attributes)},
		{Name: "unknown", Input: types.SetUnknown(element), Expected: types.ObjectUnknown(attributes)},
		{Name: "one", Input: types.SetValueMust(element, []attr.Value{first}), Expected: first},
		{Name: "multiple_preserve_first", Input: types.SetValueMust(element, []attr.Value{first, second}), Expected: first},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			result := LegacyChildFromSet(context.Background(), &diagnostics, testCase.Input.(types.Set), attributes)
			assert.False(t, diagnostics.HasError(), diagnostics)
			assert.Equal(t, testCase.Expected, result)
		})
	}
}

func TestLegacyRequiredPropertySchema(t *testing.T) {
	type input struct {
		Current types.String
		First   types.String
		Second  types.String
	}
	null, known := types.StringNull(), types.StringValue("range")
	testCases := []test.TestCase{
		{Name: "current", Input: input{known, null, null}, Expected: false},
		{Name: "first_historical_name", Input: input{null, known, null}, Expected: false},
		{Name: "second_historical_name", Input: input{null, null, known}, Expected: false},
		{Name: "missing", Input: input{null, null, null}, Expected: true},
		{Name: "conflicting_current", Input: input{known, known, null}, Expected: true},
		{Name: "conflicting_historical_names", Input: input{null, known, known}, Expected: true},
	}
	for _, withReadOnly := range []bool{false, true} {
		definition := ClassSchema{
			Properties: map[string]PropertySchema{
				"selector_type": {Type: StringProperty, Required: true, RequiresReplace: true, ValidValues: []string{"range", "ALL"}},
			},
			LegacyAttributes: map[string]LegacyAttributeSchema{
				"old_type":   {Type: types.StringType, Replacement: "selector_type"},
				"older_type": {Type: types.StringType, Replacement: "selector_type"},
			},
		}
		if withReadOnly {
			definition.Properties["read_only"] = PropertySchema{Type: StringProperty, ReadOnly: true, Computed: true}
		}
		schema := definition.ResourceSchema("")
		attribute := schema.Attributes["selector_type"].(resourceschema.StringAttribute)
		assert.False(t, attribute.Required)
		assert.True(t, attribute.Optional && attribute.Computed)
		assert.False(t, schema.ValidateImplementation(context.Background()).HasError())
		for _, testCase := range testCases {
			t.Run(testCase.Name, func(t *testing.T) {
				values := testCase.Input.(input)
				attributes := map[string]attr.Value{"id": null, "selector_type": values.Current, "old_type": values.First, "older_type": values.Second}
				typesMap := map[string]attr.Type{"id": types.StringType, "selector_type": types.StringType, "old_type": types.StringType, "older_type": types.StringType}
				if withReadOnly {
					attributes["read_only"], typesMap["read_only"] = null, types.StringType
				}
				object := types.ObjectValueMust(typesMap, attributes)
				raw, err := object.ToTerraformValue(context.Background())
				assert.NoError(t, err)
				request := validator.StringRequest{Config: tfsdk.Config{Schema: schema, Raw: raw}, ConfigValue: values.Current, Path: path.Root("selector_type"), PathExpression: path.MatchRoot("selector_type")}
				var response validator.StringResponse
				for _, check := range attribute.Validators {
					check.ValidateString(context.Background(), request, &response)
				}
				assert.Equal(t, testCase.Expected, response.Diagnostics.HasError(), response.Diagnostics)
			})
		}
	}
}

func TestLegacyParentDNSchema(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "modern_only", Input: map[string]LegacyAttributeSchema{}, Expected: true},
		{Name: "alias", Input: map[string]LegacyAttributeSchema{"application_profile_dn": {Type: types.StringType, Replacement: "parent_dn", RequiresReplace: true}}, Expected: false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			definition := ClassSchema{ParentDN: &ParentDNSchema{AttributeName: "parent_dn"}, LegacyAttributes: testCase.Input.(map[string]LegacyAttributeSchema)}
			resource := definition.ResourceSchema("")
			parent := resource.Attributes["parent_dn"].(resourceschema.StringAttribute)
			assert.Equal(t, testCase.Expected, parent.IsRequired())
			assert.Equal(t, !testCase.Expected.(bool), parent.IsOptional())
			assert.Equal(t, !testCase.Expected.(bool), parent.IsComputed())
			if !parent.IsRequired() {
				assert.Len(t, parent.Validators, 1, "either the current parent or its alias is required")
				alias := resource.Attributes["application_profile_dn"].(resourceschema.StringAttribute)
				assert.Len(t, alias.Validators, 1, "simultaneous current and old parent inputs conflict")
				assert.Len(t, alias.PlanModifiers, 1, "configured old parent changes still require replacement")
			}
			assert.True(t, definition.DataSourceSchema().Attributes["parent_dn"].IsRequired())
		})
	}
}

func TestLegacyChildSetUpgrade(t *testing.T) {
	attributeTypes := map[string]attr.Type{"name": types.StringType, "annotation": types.StringType}
	elementType := types.ObjectType{AttrTypes: attributeTypes}
	testCases := []test.TestCase{
		{Name: "named", Input: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("uni/tn-other/brc-policy-one")}), Expected: "policy-one"},
		{Name: "empty", Input: types.SetValueMust(types.StringType, []attr.Value{}), Expected: ""},
		{Name: "null", Input: types.SetNull(types.StringType), Expected: ""},
		{Name: "unknown", Input: types.SetUnknown(types.StringType), Expected: nil},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			mapping := LegacyChildSet{AttributeTypes: attributeTypes, AttributeName: "name", Named: true}
			result := mapping.Upgrade(context.Background(), &diagnostics, testCase.Input.(types.Set))
			assert.Equal(t, elementType, result.ElementType(context.Background()))
			if testCase.Expected == nil {
				assert.True(t, result.IsUnknown())
			} else if testCase.Expected == "" {
				assert.False(t, result.IsNull(), "the published SDKv2 upgrade produced an empty child set")
				assert.Empty(t, result.Elements())
			} else if assert.Len(t, result.Elements(), 1) {
				assert.Equal(t, types.StringValue(testCase.Expected.(string)), result.Elements()[0].(types.Object).Attributes()["name"])
				assert.True(t, result.Elements()[0].(types.Object).Attributes()["annotation"].IsNull())
			}
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
		})
	}
}

func TestLegacyChildSetPlan(t *testing.T) {
	attributeTypes := map[string]attr.Type{"target_dn": types.StringType, "annotation": types.StringType, "description": types.StringType}
	elementType := types.ObjectType{AttrTypes: attributeTypes}
	child := func(dn, annotation, description string) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{"target_dn": types.StringValue(dn), "annotation": types.StringValue(annotation), "description": types.StringValue(description)})
	}
	prior := types.SetValueMust(elementType, []attr.Value{child("first", "first-annotation", "first-description"), child("second", "second-annotation", "second-description")})
	testCases := []test.TestCase{
		{Name: "reordered_identifiers_retain_values", Input: []string{"second", "first"}, Expected: prior},
		{Name: "removed_identifier", Input: []string{"second"}, Expected: types.SetValueMust(elementType, []attr.Value{child("second", "second-annotation", "second-description")})},
		{Name: "changed_identifier", Input: []string{"new"}, Expected: nil},
		{Name: "cleared", Input: []string{}, Expected: types.SetValueMust(elementType, []attr.Value{})},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			values := []attr.Value{}
			for _, dn := range testCase.Input.([]string) {
				values = append(values, types.StringValue(dn))
			}
			legacy := types.SetValueMust(types.StringType, values)
			mapping := LegacyChildSet{AttributeTypes: attributeTypes, AttributeName: "target_dn"}
			result := mapping.Plan(context.Background(), &diagnostics, legacy, types.SetUnknown(elementType), &prior, types.StringValue("provider-annotation"))
			if testCase.Expected != nil {
				assert.True(t, result.Equal(testCase.Expected.(types.Set)), "%s", result)
			} else if assert.Len(t, result.Elements(), 1) {
				attributes := result.Elements()[0].(types.Object).Attributes()
				assert.Equal(t, types.StringValue("new"), attributes["target_dn"])
				assert.Equal(t, types.StringValue("provider-annotation"), attributes["annotation"])
				assert.True(t, attributes["description"].IsUnknown(), "new identifiers must not inherit an unrelated child's values")
			}
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
			assert.Len(t, prior.Elements(), 2, "planning must not mutate prior state")
		})
	}
}

func TestLegacyChildSetCombinedIdentifiers(t *testing.T) {
	attributeTypes := map[string]attr.Type{"policy": types.StringType, "filter": types.StringType, "annotation": types.StringType}
	legacyType := types.ObjectType{AttrTypes: map[string]attr.Type{"policy_name": types.StringType, "flt_type": types.StringType}}
	elementType := types.ObjectType{AttrTypes: attributeTypes}
	child := func(policy, filter, annotation string) attr.Value {
		return types.ObjectValueMust(attributeTypes, map[string]attr.Value{"policy": types.StringValue(policy), "filter": types.StringValue(filter), "annotation": types.StringValue(annotation)})
	}
	legacy := func(entries ...[2]string) types.Set {
		values := make([]attr.Value, 0, len(entries))
		for _, entry := range entries {
			values = append(values, types.ObjectValueMust(legacyType.AttrTypes, map[string]attr.Value{"policy_name": types.StringValue(entry[0]), "flt_type": types.StringValue(entry[1])}))
		}
		return types.SetValueMust(legacyType, values)
	}
	prior := types.SetValueMust(elementType, []attr.Value{child("first", "ce", "first-ce"), child("first", "ipv4", "first-ipv4"), child("second", "ce", "second-ce")})
	testCases := []test.TestCase{
		{Name: "same_policy_different_filters", Input: legacy([2]string{"first", "ipv4"}, [2]string{"first", "ce"}), Expected: types.SetValueMust(elementType, []attr.Value{child("first", "ipv4", "first-ipv4"), child("first", "ce", "first-ce")})},
		{Name: "same_filter_different_policies", Input: legacy([2]string{"second", "ce"}, [2]string{"first", "ce"}), Expected: types.SetValueMust(elementType, []attr.Value{child("second", "ce", "second-ce"), child("first", "ce", "first-ce")})},
		{Name: "new_combination", Input: legacy([2]string{"second", "ipv4"}), Expected: types.SetValueMust(elementType, []attr.Value{child("second", "ipv4", "parent-annotation")})},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			mapping := LegacyChildSet{
				AttributeTypes: attributeTypes, IdentifyingAttributes: []string{"policy", "filter"},
				LegacyType: legacyType, Attributes: map[string]string{"policy_name": "policy", "flt_type": "filter"},
			}
			result := mapping.Plan(context.Background(), &diagnostics, testCase.Input.(types.Set), types.SetUnknown(elementType), &prior, types.StringValue("parent-annotation"))
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
			assert.True(t, result.Equal(testCase.Expected.(types.Set)), "%s", result)
		})
	}
}

func TestLegacyChildSetNamedBlock(t *testing.T) {
	attributeTypes := map[string]attr.Type{"policy": types.StringType, "filter": types.StringType, "annotation": types.StringType}
	elementType := types.ObjectType{AttrTypes: attributeTypes}
	legacyType := types.ObjectType{AttrTypes: map[string]attr.Type{"policy_dn": types.StringType, "flt_type": types.StringType}}
	legacy := types.SetValueMust(legacyType, []attr.Value{types.ObjectValueMust(legacyType.AttrTypes, map[string]attr.Value{
		"policy_dn": types.StringValue("uni/tn-common/monitorpol-shared"), "flt_type": types.StringValue("ce"),
	})})
	prior := types.SetValueMust(elementType, []attr.Value{types.ObjectValueMust(attributeTypes, map[string]attr.Value{
		"policy": types.StringValue("shared"), "filter": types.StringValue("ce"), "annotation": types.StringValue("retained"),
	})})
	testCases := []test.TestCase{
		{Name: "plan", Expected: types.StringValue("retained")},
		{Name: "upgrade", Expected: types.StringNull()},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			mapping := LegacyChildSet{
				AttributeTypes: attributeTypes, AttributeName: "policy", Named: true, IdentifyingAttributes: []string{"policy", "filter"},
				LegacyType: legacyType, Attributes: map[string]string{"policy_dn": "policy", "flt_type": "filter"},
			}
			var result types.Set
			if testCase.Name == "plan" {
				result = mapping.Plan(context.Background(), &diagnostics, legacy, types.SetUnknown(elementType), &prior, types.StringValue("parent"))
			} else {
				result = mapping.Upgrade(context.Background(), &diagnostics, legacy)
			}
			if assert.False(t, diagnostics.HasError(), "%v", diagnostics) && assert.Len(t, result.Elements(), 1) {
				attributes := result.Elements()[0].(types.Object).Attributes()
				assert.Equal(t, types.StringValue("shared"), attributes["policy"])
				assert.Equal(t, types.StringValue("ce"), attributes["filter"])
				assert.Equal(t, testCase.Expected, attributes["annotation"])
			}
		})
	}
}

func TestLegacyChildSetFromResponse(t *testing.T) {
	legacyType := types.ObjectType{AttrTypes: map[string]attr.Type{"policy_dn": types.StringType, "flt_type": types.StringType}}
	testCases := []test.TestCase{
		{Name: "resolved_common_policy", Input: `{"imdata":[{"fvBD":{"children":[{"fvRsBDToNetflowMonitorPol":{"attributes":{"fltType":"ce","tnNetflowMonitorPolName":"shared","tDn":"uni/tn-common/monitorpol-shared"}}}]}}]}`, Expected: types.StringValue("uni/tn-common/monitorpol-shared")},
		{Name: "missing_target", Input: `{"imdata":[{"fvBD":{"children":[{"fvRsBDToNetflowMonitorPol":{"attributes":{"fltType":"ce","tnNetflowMonitorPolName":"shared"}}}]}}]}`, Expected: types.StringNull()},
		{Name: "malformed_target", Input: `{"imdata":[{"fvBD":{"children":[{"fvRsBDToNetflowMonitorPol":{"attributes":{"fltType":"ce","tDn":42}}}]}}]}`},
		{Name: "empty", Input: `{"imdata":[]}`, Expected: 0},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			response, err := container.ParseJSON([]byte(testCase.Input.(string)))
			if !assert.NoError(t, err) {
				return
			}
			var diagnostics diag.Diagnostics
			mapping := LegacyChildSet{LegacyType: legacyType, ResponseAttributes: map[string]string{"policy_dn": "tDn", "flt_type": "fltType"}}
			result := mapping.LegacyFromResponse(context.Background(), &diagnostics, response, "fvBD", "fvRsBDToNetflowMonitorPol")
			if testCase.Expected == nil {
				assert.True(t, diagnostics.HasError())
				return
			}
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
			if testCase.Expected == 0 {
				assert.Empty(t, result.Elements())
			} else if assert.Len(t, result.Elements(), 1) {
				attributes := result.Elements()[0].(types.Object).Attributes()
				assert.Equal(t, testCase.Expected, attributes["policy_dn"])
				assert.Equal(t, types.StringValue("ce"), attributes["flt_type"])
			}
		})
	}
}

func TestLegacyTargetDNs(t *testing.T) {
	testCases := []test.TestCase{
		{Name: "multiple", Input: `{"imdata":[{"fvAEPg":{"children":[{"fvRsCons":{"attributes":{"tDn":"uni/tn-common/brc-one"}}},{"fvRsCons":{"attributes":{"tDn":"uni/tn-other/brc-two"}}},{"fvRsProv":{"attributes":{"tDn":"ignored"}}}]}}]}`, Expected: 2},
		{Name: "missing", Input: `{"imdata":[{"fvAEPg":{"attributes":{}}}]}`, Expected: 0},
		{Name: "empty", Input: `{"imdata":[]}`, Expected: 0},
		{Name: "malformed_target", Input: `{"imdata":[{"fvAEPg":{"children":[{"fvRsCons":{"attributes":{"tDn":42}}}]}}]}`},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			response, err := container.ParseJSON([]byte(testCase.Input.(string)))
			if !assert.NoError(t, err) {
				return
			}
			var diagnostics diag.Diagnostics
			result := LegacyTargetDNs(context.Background(), &diagnostics, response, "fvAEPg", "fvRsCons")
			if testCase.Expected == nil {
				assert.True(t, diagnostics.HasError())
				return
			}
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
			assert.Len(t, result.Elements(), testCase.Expected.(int))
			if testCase.Expected == 2 {
				assert.True(t, result.Equal(types.SetValueMust(types.StringType, []attr.Value{types.StringValue("uni/tn-common/brc-one"), types.StringValue("uni/tn-other/brc-two")})))
			}
		})
	}
}

func TestLegacyChildBlock(t *testing.T) {
	ctx := context.Background()
	currentTypes := map[string]attr.Type{"target_dn": types.StringType, "description": types.StringType, "annotation": types.StringType}
	legacyType := types.ObjectType{AttrTypes: map[string]attr.Type{"node_dn": types.StringType, "description": types.StringType}}
	mapping := LegacyChildSet{AttributeTypes: currentTypes, AttributeName: "target_dn", LegacyType: legacyType, Attributes: map[string]string{"node_dn": "target_dn", "description": "description"}}
	legacy := types.SetValueMust(legacyType, []attr.Value{types.ObjectValueMust(legacyType.AttrTypes, map[string]attr.Value{"node_dn": types.StringValue("node-101"), "description": types.StringValue("configured")})})
	prior := types.SetValueMust(types.ObjectType{AttrTypes: currentTypes}, []attr.Value{types.ObjectValueMust(currentTypes, map[string]attr.Value{"target_dn": types.StringValue("node-101"), "description": types.StringValue("previous"), "annotation": types.StringValue("retained")})})
	testCases := []test.TestCase{
		{Name: "upgrade", Input: legacy, Expected: types.StringNull()},
		{Name: "plan_preserves_newer_fields", Input: legacy, Expected: types.StringValue("retained")},
		{Name: "unknown", Input: types.SetUnknown(legacyType)},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			result := mapping.Upgrade(ctx, &diagnostics, testCase.Input.(types.Set))
			if testCase.Name == "plan_preserves_newer_fields" {
				result = mapping.Plan(ctx, &diagnostics, legacy, types.SetUnknown(types.ObjectType{AttrTypes: currentTypes}), &prior, types.StringValue("default"))
			}
			if testCase.Expected == nil {
				assert.True(t, result.IsUnknown())
			} else if assert.Len(t, result.Elements(), 1) {
				attributes := result.Elements()[0].(types.Object).Attributes()
				assert.Equal(t, testCase.Expected, attributes["annotation"])
				assert.Equal(t, types.StringValue("configured"), attributes["description"])
				assert.True(t, legacy.Equal(mapping.Legacy(ctx, &diagnostics, result)), "readback must contain exactly the old block's members")
				unknown := mapping.Unknown(ctx, &diagnostics, result)
				if assert.Len(t, unknown.Elements(), 1) {
					for _, value := range unknown.Elements()[0].(types.Object).Attributes() {
						assert.True(t, value.IsUnknown())
					}
				}
			}
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
		})
	}
}

func TestLegacyBlockSchema(t *testing.T) {
	definition := ClassSchema{LegacyAttributes: map[string]LegacyAttributeSchema{"old_leaf": {
		Replacement: "new_leafs", Attributes: map[string]LegacyAttributeSchema{
			"node_dn":     {Type: types.StringType, Required: true, Replacement: "new_leafs.target_dn"},
			"description": {Type: types.StringType, Replacement: "new_leafs.description"},
		},
	}}}
	resource := definition.ResourceSchema("")
	datasource := definition.DataSourceSchema()
	assert.NotContains(t, resource.Attributes, "old_leaf", "legacy HCL block syntax must be preserved")
	assert.NotContains(t, datasource.Attributes, "old_leaf")
	resourceBlock := resource.Blocks["old_leaf"].(resourceschema.SetNestedBlock)
	assert.Len(t, resourceBlock.NestedObject.Attributes, 2)
	identifier := resourceBlock.NestedObject.Attributes["node_dn"].(resourceschema.StringAttribute)
	assert.True(t, identifier.IsOptional())
	assert.True(t, identifier.IsComputed())
	assert.Len(t, identifier.Validators, 2, "required members validate configuration and conflict with the new child")
	assert.False(t, resource.ValidateImplementation(context.Background()).HasError())
	assert.False(t, datasource.ValidateImplementation(context.Background()).HasError())
}

func TestLegacyBlockUnknown(t *testing.T) {
	elementType := types.ObjectType{AttrTypes: map[string]attr.Type{"node_dn": types.StringType}}
	child := types.ObjectValueMust(elementType.AttrTypes, map[string]attr.Value{"node_dn": types.StringValue("node-101")})
	testCases := []test.TestCase{
		{Name: "unknown_children_keep_absent_block_empty", Input: types.SetUnknown(elementType), Expected: 0},
		{Name: "null_children", Input: types.SetNull(elementType), Expected: 0},
		{Name: "empty_children", Input: types.SetValueMust(elementType, []attr.Value{}), Expected: 0},
		{Name: "known_child_count", Input: types.SetValueMust(elementType, []attr.Value{child}), Expected: 1},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			result := (LegacyChildSet{LegacyType: elementType}).Unknown(context.Background(), &diagnostics, testCase.Input.(types.Set))
			assert.False(t, diagnostics.HasError(), "%v", diagnostics)
			assert.False(t, result.IsUnknown(), "a non-computed legacy block must never become wholly unknown")
			assert.Len(t, result.Elements(), testCase.Expected.(int))
			for _, element := range result.Elements() {
				assert.True(t, element.(types.Object).Attributes()["node_dn"].IsUnknown())
			}
		})
	}
}
