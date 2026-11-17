package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

type resourceTestDNVariant struct {
	Patterns  []string
	RNPrepend string
}

func testCheckResourceDN(address, rnFormat, fixedParent string, identifiers map[string]string, variants []resourceTestDNVariant) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		root := state.RootModule()
		if root == nil || root.Resources[address] == nil || root.Resources[address].Primary == nil {
			return fmt.Errorf("resource %q is missing from state", address)
		}
		attributes := root.Resources[address].Primary.Attributes
		rn := rnFormat
		for property, attribute := range identifiers {
			rn = strings.ReplaceAll(rn, fmt.Sprintf("{%s}", property), attributes[attribute])
		}
		parent := fixedParent
		if value := attributes["parent_dn"]; value != "" {
			parent = value
		}
		for _, variant := range variants {
			for _, pattern := range variant.Patterns {
				if regexp.MustCompile(pattern).MatchString(parent) {
					rn = fmt.Sprintf("%s/%s", variant.RNPrepend, rn)
					return resource.TestCheckResourceAttr(address, "id", fmt.Sprintf("%s/%s", parent, rn))(state)
				}
			}
		}
		return resource.TestCheckResourceAttr(address, "id", fmt.Sprintf("%s/%s", parent, rn))(state)
	}
}

func testSkipUnsupportedResourceChanges(t *testing.T, groups [][]string) (bool, error) {
	// Each group is an intersection of owner/property/dependency versions; any supported change permits the step.
	version := getAPICVersionTest(t)
	for _, group := range groups {
		supported := true
		for _, versions := range group {
			inside, err := CompareVersionsRange(version, versions, "inside")
			if err != nil {
				return false, err
			}
			supported = supported && inside
		}
		if supported {
			return false, nil
		}
	}
	t.Logf("Skipping resource test step: APIC %s does not support its changed attributes or children (%v)", version, groups)
	return true, nil
}

func testDeleteResourceExternally(t *testing.T, class, dn string) {
	t.Helper()
	if dn == "" {
		t.Fatal("cannot delete a resource externally before its ID is known")
	}
	var diagnostics diag.Diagnostics
	payload := GetDeleteJsonPayload(context.Background(), &diagnostics, class, dn)
	if diagnostics.HasError() {
		t.Fatal(diagnostics)
	}
	DoRestRequest(context.Background(), &diagnostics, getACIClientTest(t), fmt.Sprintf("api/mo/%s.json", dn), "POST", payload)
	if diagnostics.HasError() {
		t.Fatalf("delete %s externally: %v", dn, diagnostics)
	}
	if err := testResourcePresence(dn, false); err != nil {
		t.Fatal(err)
	}
}

func testResourcePresence(dn string, expected bool) error {
	var diagnostics diag.Diagnostics
	response := getManagedObject(context.Background(), &diagnostics, getACIClientTest(nil), dn, nil)
	if diagnostics.HasError() {
		return fmt.Errorf("check presence of %s: %v", dn, diagnostics)
	}
	if response == nil {
		return fmt.Errorf("check presence of %s: empty APIC response", dn)
	}
	children, err := response.S("imdata").Children()
	if err != nil {
		return fmt.Errorf("check presence of %s: %w", dn, err)
	}
	if (len(children) > 0) != expected {
		return fmt.Errorf("object %s: expected present=%t, actual present=%t", dn, expected, len(children) > 0)
	}
	return nil
}

func testCheckGeneratedResourceDestroy(nonDeletableTypes ...string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for address, instance := range state.RootModule().Resources {
			if strings.HasPrefix(address, "data.") || !strings.HasPrefix(instance.Type, "aci_") || instance.Primary == nil {
				continue
			}
			ignored := false
			for _, resourceType := range nonDeletableTypes {
				ignored = ignored || instance.Type == resourceType
			}
			if !ignored {
				if err := testResourcePresence(instance.Primary.ID, false); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func TestResourceAcceptanceDN(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		RN         string
		Parent     string
		Attributes map[string]string
		IDs        map[string]string
		Variants   []resourceTestDNVariant
	}
	testCases := []test.TestCase{
		{Name: "test_fixed_root", Input: input{RN: "tn-{name}", Parent: "uni", Attributes: map[string]string{"id": "uni/tn-name_1", "name": "name_1"}, IDs: map[string]string{"name": "name"}}},
		{Name: "test_fixed_infra_path", Input: input{RN: "infra/l2IfP-{name}", Parent: "uni", Attributes: map[string]string{"id": "uni/infra/l2IfP-name_1", "name": "name_1"}, IDs: map[string]string{"name": "name"}}},
		{Name: "test_multiple_identifiers", Input: input{RN: "rule-{direction}-{port}", Attributes: map[string]string{"id": "uni/tn-test/rule-in-80", "parent_dn": "uni/tn-test", "direction": "in", "port": "80"}, IDs: map[string]string{"direction": "direction", "port": "port"}}},
		{Name: "test_bracketed_dn_identifier", Input: input{RN: "rsdomAtt-[{tDn}]", Attributes: map[string]string{"id": "uni/tn-test/rsdomAtt-[uni/phys-test]", "parent_dn": "uni/tn-test", "target_dn": "uni/phys-test"}, IDs: map[string]string{"tDn": "target_dn"}}},
		{Name: "test_parent_wrapper", Input: input{RN: "keyring-{name}", Attributes: map[string]string{"id": "uni/tn-test/certstore/keyring-name_1", "parent_dn": "uni/tn-test", "name": "name_1"}, IDs: map[string]string{"name": "name"}, Variants: []resourceTestDNVariant{{Patterns: []string{`^uni/tn-[^/]+$`}, RNPrepend: "certstore"}}}},
		{Name: "test_wrong_dn_fails", Input: input{RN: "tn-{name}", Parent: "uni", Attributes: map[string]string{"id": "uni/tn-other", "name": "name_1"}, IDs: map[string]string{"name": "name"}}, Expected: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(input)
			state := &terraform.State{Modules: []*terraform.ModuleState{{Path: []string{"root"}, Resources: map[string]*terraform.ResourceState{
				"aci_example.test_1": {Type: "aci_example", Primary: &terraform.InstanceState{ID: input.Attributes["id"], Attributes: input.Attributes}},
			}}}}
			err := testCheckResourceDN("aci_example.test_1", input.RN, input.Parent, input.IDs, input.Variants)(state)
			if fail, _ := testCase.Expected.(bool); fail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
