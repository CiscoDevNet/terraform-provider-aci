package provider

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAciRestManaged_childPayloadsOnRemoval(t *testing.T) {
	stringMap := func(values map[string]string) types.Map {
		elements := map[string]attr.Value{}
		for k, v := range values {
			elements[k] = types.StringValue(v)
		}
		return types.MapValueMust(types.StringType, elements)
	}
	child := func(rn, className string, content, contentOnDestroy map[string]string) ChildAciRestManagedResourceModel {
		model := ChildAciRestManagedResourceModel{
			Rn:               types.StringValue(rn),
			ClassName:        types.StringValue(className),
			Content:          types.MapNull(types.StringType),
			ContentOnDestroy: types.MapNull(types.StringType),
		}
		if content != nil {
			model.Content = stringMap(content)
		}
		if contentOnDestroy != nil {
			model.ContentOnDestroy = stringMap(contentOnDestroy)
		}
		return model
	}
	// The payload builder only checks whether the planned child set is null, not its elements
	nonEmptyChildren := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("placeholder")})
	nullChildren := types.SetNull(types.StringType)

	tests := []struct {
		name       string
		annotation types.String
		children   types.Set
		plan       []ChildAciRestManagedResourceModel
		state      []ChildAciRestManagedResourceModel
		expected   string
	}{
		{
			name:     "removed child with content_on_destroy is modified and others are kept",
			children: nonEmptyChildren,
			plan:     []ChildAciRestManagedResourceModel{child("rn-keep", "classKeep", map[string]string{"a": "1"}, nil)},
			state: []ChildAciRestManagedResourceModel{
				child("rn-keep", "classKeep", map[string]string{"a": "1"}, nil),
				child("rshIfPol", "infraRsHIfPol", map[string]string{"tnFabricHIfPolName": "default"}, map[string]string{"tnFabricHIfPolName": ""}),
			},
			expected: `[{"classKeep":{"attributes":{"a":"1","rn":"rn-keep"}}},{"infraRsHIfPol":{"attributes":{"rn":"rshIfPol","tnFabricHIfPolName":""}}}]`,
		},
		{
			name:     "removed child without content_on_destroy is deleted",
			children: nonEmptyChildren,
			plan:     []ChildAciRestManagedResourceModel{child("rn-keep", "classKeep", map[string]string{"a": "1"}, nil)},
			state: []ChildAciRestManagedResourceModel{
				child("rn-keep", "classKeep", map[string]string{"a": "1"}, nil),
				child("rn-removed", "classRemoved", map[string]string{"b": "2"}, nil),
			},
			expected: `[{"classKeep":{"attributes":{"a":"1","rn":"rn-keep"}}},{"classRemoved":{"attributes":{"rn":"rn-removed","status":"deleted"}}}]`,
		},
		{
			name:     "last child removed with content_on_destroy is modified",
			children: nullChildren,
			state:    []ChildAciRestManagedResourceModel{child("rshIfPol", "infraRsHIfPol", map[string]string{"tnFabricHIfPolName": "default"}, map[string]string{"tnFabricHIfPolName": ""})},
			expected: `[{"infraRsHIfPol":{"attributes":{"rn":"rshIfPol","tnFabricHIfPolName":""}}}]`,
		},
		{
			name:     "last child removed without content_on_destroy is deleted",
			children: nullChildren,
			state:    []ChildAciRestManagedResourceModel{child("rn-removed", "classRemoved", map[string]string{"b": "2"}, nil)},
			expected: `[{"classRemoved":{"attributes":{"rn":"rn-removed","status":"deleted"}}}]`,
		},
		{
			name:     "edited child with the same rn and class is not cleaned up",
			children: nonEmptyChildren,
			plan:     []ChildAciRestManagedResourceModel{child("rshIfPol", "infraRsHIfPol", map[string]string{"tnFabricHIfPolName": "inherit"}, map[string]string{"tnFabricHIfPolName": ""})},
			state:    []ChildAciRestManagedResourceModel{child("rshIfPol", "infraRsHIfPol", map[string]string{"tnFabricHIfPolName": "default"}, map[string]string{"tnFabricHIfPolName": ""})},
			expected: `[{"infraRsHIfPol":{"attributes":{"rn":"rshIfPol","tnFabricHIfPolName":"inherit"}}}]`,
		},
		{
			name:       "annotation is added to the modified child unless the class does not support it",
			annotation: types.StringValue("orchestrator:terraform"),
			children:   nullChildren,
			state: []ChildAciRestManagedResourceModel{
				child("rshIfPol", "infraRsHIfPol", nil, map[string]string{"tnFabricHIfPolName": ""}),
				child("domain-x", "aaaADomainRef", nil, map[string]string{"descr": "cleanup"}),
			},
			expected: `[{"infraRsHIfPol":{"attributes":{"annotation":"orchestrator:terraform","rn":"rshIfPol","tnFabricHIfPolName":""}}},{"aaaADomainRef":{"attributes":{"descr":"cleanup","rn":"domain-x"}}}]`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data := &AciRestManagedResourceModel{
				Annotation: types.StringNull(),
				Child:      test.children,
			}
			if !test.annotation.IsNull() {
				data.Annotation = test.annotation
			}

			var diags diag.Diagnostics
			payloads := getAciRestManagedChildPayloads(context.Background(), &diags, data, test.plan, test.state)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}

			actual, err := json.Marshal(payloads)
			if err != nil {
				t.Fatalf("failed to marshal payloads: %s", err)
			}
			if string(actual) != test.expected {
				t.Errorf("unexpected child payloads\nexpected: %s\nactual:   %s", test.expected, actual)
			}
		})
	}
}
