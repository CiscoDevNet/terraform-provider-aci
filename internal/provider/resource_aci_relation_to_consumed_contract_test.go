// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceFvRsConsWithFvAEPg(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvAEPgAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "priority", "unspecified"),
				),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "false")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:      testConfigFvRsConsMinDependencyWithFvAEPgAllowExisting,
				ExpectError: regexp.MustCompile("Object Already Exists"),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "true")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvAEPgAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "priority", "unspecified"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
				),
			},
			// Update with all config and verify default APIC values
			{
				Config:             testConfigFvRsConsAllDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "annotation"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "level1"),
				),
			},
			// Update with minimum config and verify config is unchanged
			{
				Config:             testConfigFvRsConsMinDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
				),
			},
			// Update with empty strings config or default value
			{
				Config:             testConfigFvRsConsResetDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
				),
			},
			// Import testing
			{
				ResourceName:      "aci_relation_to_consumed_contract.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children
			{
				Config:             testConfigFvRsConsChildrenDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.value", "test_value"),
				),
			},
			// Refresh State before import testing to ensure that the state is up to date
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: false,
			},
			// Import testing with children
			{
				ResourceName:      "aci_relation_to_consumed_contract.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children removed from config
			{
				Config:             testConfigFvRsConsChildrenRemoveFromConfigDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "2"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "2"),
				),
			},
			// Update with children first child removed
			{
				Config:             testConfigFvRsConsChildrenRemoveOneDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "1"),
				),
			},
			// Update with all children removed
			{
				Config:             testConfigFvRsConsChildrenRemoveAllDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "0"),
				),
			},
			// Update with minimum config and custom type semantic equivalent values
			{
				Config:             testConfigFvRsConsCustomTypeDependencyWithFvAEPg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "1"),
				),
			},
		},
		CheckDestroy: testCheckResourceDestroy,
	})
}
func TestAccResourceFvRsConsWithFvESg(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvESgAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "priority", "unspecified"),
				),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "false")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:      testConfigFvRsConsMinDependencyWithFvESgAllowExisting,
				ExpectError: regexp.MustCompile("Object Already Exists"),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "true")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvESgAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.allow_test_2", "priority", "unspecified"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, "both", "1.0(1e)-") },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvRsConsMinDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
				),
			},
			// Update with all config and verify default APIC values
			{
				Config:             testConfigFvRsConsAllDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "annotation"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "level1"),
				),
			},
			// Update with minimum config and verify config is unchanged
			{
				Config:             testConfigFvRsConsMinDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
				),
			},
			// Update with empty strings config or default value
			{
				Config:             testConfigFvRsConsResetDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
				),
			},
			// Import testing
			{
				ResourceName:      "aci_relation_to_consumed_contract.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children
			{
				Config:             testConfigFvRsConsChildrenDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "unspecified"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.value", "test_value"),
				),
			},
			// Refresh State before import testing to ensure that the state is up to date
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: false,
			},
			// Import testing with children
			{
				ResourceName:      "aci_relation_to_consumed_contract.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children removed from config
			{
				Config:             testConfigFvRsConsChildrenRemoveFromConfigDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "2"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.1.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "2"),
				),
			},
			// Update with children first child removed
			{
				Config:             testConfigFvRsConsChildrenRemoveOneDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.0.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.0.value", "test_value"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "1"),
				),
			},
			// Update with all children removed
			{
				Config:             testConfigFvRsConsChildrenRemoveAllDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "annotations.#", "0"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "tags.#", "0"),
				),
			},
			// Update with minimum config and custom type semantic equivalent values
			{
				Config:             testConfigFvRsConsCustomTypeDependencyWithFvESg,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "contract_name", "test_tn_vz_br_cp_name"),
					resource.TestCheckResourceAttr("aci_relation_to_consumed_contract.test", "priority", "1"),
				),
			},
		},
		CheckDestroy: testCheckResourceDestroy,
	})
}

const testDependencyConfigFvRsCons = `
resource "aci_contract" "test_contract_0" {
  tenant_dn = aci_tenant.test.id
  name = "contract_name_0"
}
resource "aci_contract" "test_contract_1" {
  tenant_dn = aci_tenant.test.id
  name = "contract_name_1"
}
`

const testConfigFvRsConsMinDependencyWithFvAEPgAllowExisting = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "allow_test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
resource "aci_relation_to_consumed_contract" "allow_test_2" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  depends_on = [aci_relation_to_consumed_contract.allow_test]
}
`

const testConfigFvRsConsMinDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
`

const testConfigFvRsConsAllDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotation = "annotation"
  priority = "level1"
}
`

const testConfigFvRsConsResetDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotation = "orchestrator:terraform"
  priority = "unspecified"
}
`
const testConfigFvRsConsChildrenDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = [
    {
      key = "key_0"
      value = "value_1"
    },
    {
      key = "key_1"
      value = "test_value"
    },
  ]
  tags = [
    {
      key = "key_0"
      value = "value_1"
    },
    {
      key = "key_1"
      value = "test_value"
    },
  ]
}
`

const testConfigFvRsConsChildrenRemoveFromConfigDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
`

const testConfigFvRsConsChildrenRemoveOneDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = [ 
	{
	  key = "key_1"
	  value = "test_value"
	},
  ]
  tags = [ 
	{
	  key = "key_1"
	  value = "test_value"
	},
  ]
}
`

const testConfigFvRsConsChildrenRemoveAllDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = []
  tags = []
}
`

const testConfigFvRsConsCustomTypeDependencyWithFvAEPg = testDependencyConfigFvRsCons + testConfigFvAEPgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_application_epg.test.id
  contract_name = "test_tn_vz_br_cp_name"
  priority = "1"
}
`

const testConfigFvRsConsMinDependencyWithFvESgAllowExisting = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "allow_test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
resource "aci_relation_to_consumed_contract" "allow_test_2" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  depends_on = [aci_relation_to_consumed_contract.allow_test]
}
`

const testConfigFvRsConsMinDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
`

const testConfigFvRsConsAllDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotation = "annotation"
  priority = "level1"
}
`

const testConfigFvRsConsResetDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotation = "orchestrator:terraform"
  priority = "unspecified"
}
`
const testConfigFvRsConsChildrenDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = [
    {
      key = "key_0"
      value = "value_1"
    },
    {
      key = "key_1"
      value = "test_value"
    },
  ]
  tags = [
    {
      key = "key_0"
      value = "value_1"
    },
    {
      key = "key_1"
      value = "test_value"
    },
  ]
}
`

const testConfigFvRsConsChildrenRemoveFromConfigDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
}
`

const testConfigFvRsConsChildrenRemoveOneDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = [ 
	{
	  key = "key_1"
	  value = "test_value"
	},
  ]
  tags = [ 
	{
	  key = "key_1"
	  value = "test_value"
	},
  ]
}
`

const testConfigFvRsConsChildrenRemoveAllDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  annotations = []
  tags = []
}
`

const testConfigFvRsConsCustomTypeDependencyWithFvESg = testDependencyConfigFvRsCons + testConfigFvESgMinDependencyWithFvAp + `
resource "aci_relation_to_consumed_contract" "test" {
  parent_dn = aci_endpoint_security_group.test.id
  contract_name = "test_tn_vz_br_cp_name"
  priority = "1"
}
`
