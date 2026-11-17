package acceptancetest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/CiscoDevNet/terraform-provider-aci/v2/gen/utils/test"
)

func TestResourceStateID(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	type input struct {
		Previous string
		Current  string
		Changed  bool
	}
	testCases := []test.TestCase{
		{Name: "test_first_apply", Input: input{Current: "uni/tn-first"}},
		{Name: "test_update_retains_id", Input: input{Previous: "uni/tn-first", Current: "uni/tn-first"}},
		{Name: "test_replacement_changes_id", Input: input{Previous: "uni/tn-first", Current: "uni/tn-second", Changed: true}},
		{Name: "test_update_cannot_replace", Input: input{Previous: "uni/tn-first", Current: "uni/tn-second"}, Expected: "expected change: false"},
		{Name: "test_replacement_must_change_id", Input: input{Previous: "uni/tn-first", Current: "uni/tn-first", Changed: true}, Expected: "expected change: true"},
		{Name: "test_empty_id_fails", Input: input{}, Expected: "expected to be set"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()
			input := testCase.Input.(input)
			identity := ResourceState{ID: input.Previous}
			state := &terraform.State{Modules: []*terraform.ModuleState{{Path: []string{"root"}, Resources: map[string]*terraform.ResourceState{
				"aci_tenant.test_1": {Type: "aci_tenant", Primary: &terraform.InstanceState{ID: input.Current, Attributes: map[string]string{"id": input.Current}}},
			}}}}
			err := identity.CheckID("aci_tenant.test_1", input.Changed)(state)
			if expected, _ := testCase.Expected.(string); expected != "" {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), expected)
				}
				assert.Equal(t, input.Previous, identity.ID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, input.Current, identity.ID)
			}
		})
	}
}

func TestResourceStateRemoved(t *testing.T) {
	t.Parallel()
	test.InitializeTest(t)
	identity := ResourceState{ID: "uni/tn-example"}
	state := &terraform.State{Modules: []*terraform.ModuleState{{Path: []string{"root"}, Resources: map[string]*terraform.ResourceState{}}}}
	assert.NoError(t, identity.CheckRemoved("aci_tenant.test_1")(state))
	state.RootModule().Resources["aci_tenant.test_1"] = &terraform.ResourceState{}
	assert.ErrorContains(t, identity.CheckRemoved("aci_tenant.test_1")(state), "remains in Terraform state")
	assert.Equal(t, "uni/tn-example", identity.ID, "remote deletion checks still need the last applied ID")
}
