package data

import (
	"testing"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestExampleValueHCL(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{Name: "test_string", Input: ExampleLiteral(`value "one"`), Expected: `"value \"one\""`},
		{Name: "test_empty_string", Input: ExampleLiteral(""), Expected: `""`},
		{Name: "test_multiline_string", Input: ExampleLiteral("line one\nline two"), Expected: `"line one\nline two"`},
		{Name: "test_expression", Input: ExampleExpression(`file("certificate.pem")`), Expected: `file("certificate.pem")`},
		{
			Name:     "test_resource_reference",
			Input:    newExampleReference(ResourceReference, "tenant", "id"),
			Expected: "aci_tenant.example.id",
		},
		{
			Name:     "test_datasource_reference",
			Input:    newExampleReference(DataSourceReference, "tenant", "id"),
			Expected: "data.aci_tenant.example.id",
		},
		{
			Name: "test_collection",
			Input: ExampleCollection{
				ExampleLiteral("one"),
				newExampleReference(ResourceReference, "tenant", "name"),
			},
			Expected: `["one", aci_tenant.example.name]`,
		},
		{Name: "test_empty_collection", Input: ExampleCollection{}, Expected: "[]"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testCase.Input.(ExampleValue).HCL()
			if actual != testCase.Expected.(string) {
				t.Fatalf("HCL value %q, expected %q", actual, testCase.Expected)
			}
		})
	}
}

func TestDataSourceExampleHCL(t *testing.T) {
	t.Parallel()

	testCases := []test.TestCase{
		{
			Name: "test_apic",
			Input: Example{
				Body: ExampleObject{Attributes: []ExampleAttribute{
					{Name: "parent_dn", Value: ExampleLiteral("uni/tn-example")},
					{Name: "name", Value: ExampleLiteral("example")},
				}},
			},
			Expected: "data \"aci_tenant\" \"example\" {\n  parent_dn = \"uni/tn-example\"\n  name      = \"example\"\n}",
		},
		{
			Name: "test_cloud",
			Input: Example{
				Platform: Cloud,
				Body: ExampleObject{Attributes: []ExampleAttribute{
					{Name: "name", Value: ExampleLiteral("example")},
				}},
			},
			Expected: "# This example is only applicable to Cisco Cloud Network Controller.\ndata \"aci_tenant\" \"example\" {\n  name = \"example\"\n}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			actual := testCase.Input.(Example).DataSourceHCL("tenant")
			if actual != testCase.Expected.(string) {
				t.Fatalf("data-source HCL %q, expected %q", actual, testCase.Expected)
			}
		})
	}
}
