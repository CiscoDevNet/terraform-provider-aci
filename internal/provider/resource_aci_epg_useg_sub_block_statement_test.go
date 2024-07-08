// Code generated by "gen/generator.go"; DO NOT EDIT.
// In order to regenerate this file execute `go generate` from the repository root.
// More details can be found in the [README](https://github.com/CiscoDevNet/terraform-provider-aci/blob/master/README.md).

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceFvSCrtrnWithFvCrtrn(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvCrtrnAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_tag", ""),
				),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "false")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:      testConfigFvSCrtrnMinDependencyWithFvCrtrnAllowExisting,
				ExpectError: regexp.MustCompile("Object Already Exists"),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "true")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvCrtrnAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_tag", ""),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", ""),
				),
			},
			// Update with all config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnAllDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "annotation"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", "description"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "all"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", "name_alias"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", "owner_key"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", "owner_tag"),
				),
			},
			// Update with minimum config and verify config is unchanged
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
				),
			},
			// Update with empty strings config or default value
			{
				Config:             testConfigFvSCrtrnResetDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", ""),
				),
			},
			// Import testing
			{
				ResourceName:      "aci_epg_useg_sub_block_statement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children
			{
				Config:             testConfigFvSCrtrnChildrenDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.1.value", "value_2"),
				),
			},
			// Import testing with children
			{
				ResourceName:      "aci_epg_useg_sub_block_statement.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children removed from config
			{
				Config:             testConfigFvSCrtrnChildrenRemoveFromConfigDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.#", "2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.#", "2"),
				),
			},
			// Update with children first child removed
			{
				Config:             testConfigFvSCrtrnChildrenRemoveOneDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.0.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.#", "1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.0.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.#", "1"),
				),
			},
			// Update with all children removed
			{
				Config:             testConfigFvSCrtrnChildrenRemoveAllDependencyWithFvCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "annotations.#", "0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test", "tags.#", "0"),
				),
			},
		},
	})
}
func TestAccResourceFvSCrtrnWithFvSCrtrn(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvSCrtrnAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_tag", ""),
				),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "false")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:      testConfigFvSCrtrnMinDependencyWithFvSCrtrnAllowExisting,
				ExpectError: regexp.MustCompile("Object Already Exists"),
			},
		},
	})

	setEnvVariable(t, "ACI_ALLOW_EXISTING_ON_CREATE", "true")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvSCrtrnAllowExisting,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_0", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_2", "owner_tag", ""),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with minimum config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_tag", ""),
				),
			},
			// Update with all config and verify default APIC values
			{
				Config:             testConfigFvSCrtrnAllDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotation", "annotation"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "description", "description"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "match", "all"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name_alias", "name_alias"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_key", "owner_key"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_tag", "owner_tag"),
				),
			},
			// Update with minimum config and verify config is unchanged
			{
				Config:             testConfigFvSCrtrnMinDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name", "sub_criterion"),
				),
			},
			// Update with empty strings config or default value
			{
				Config:             testConfigFvSCrtrnResetDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_tag", ""),
				),
			},
			// Import testing
			{
				ResourceName:      "aci_epg_useg_sub_block_statement.test_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children
			{
				Config:             testConfigFvSCrtrnChildrenDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name", "sub_criterion"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "description", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "match", "any"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "name_alias", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_key", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "owner_tag", ""),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.1.value", "value_2"),
				),
			},
			// Import testing with children
			{
				ResourceName:      "aci_epg_useg_sub_block_statement.test_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update with children removed from config
			{
				Config:             testConfigFvSCrtrnChildrenRemoveFromConfigDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.#", "2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.key", "key_0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.value", "value_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.1.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.1.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.#", "2"),
				),
			},
			// Update with children first child removed
			{
				Config:             testConfigFvSCrtrnChildrenRemoveOneDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.0.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.#", "1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.key", "key_1"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.0.value", "value_2"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.#", "1"),
				),
			},
			// Update with all children removed
			{
				Config:             testConfigFvSCrtrnChildrenRemoveAllDependencyWithFvSCrtrn,
				ExpectNonEmptyPlan: false,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "annotations.#", "0"),
					resource.TestCheckResourceAttr("aci_epg_useg_sub_block_statement.test_1", "tags.#", "0"),
				),
			},
		},
	})
}

const testConfigFvSCrtrnMinDependencyWithFvCrtrnAllowExisting = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
}
resource "aci_epg_useg_sub_block_statement" "test_2" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  depends_on = [aci_epg_useg_sub_block_statement.test]
}
`

const testConfigFvSCrtrnMinDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
}
`

const testConfigFvSCrtrnAllDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  annotation = "annotation"
  description = "description"
  match = "all"
  name_alias = "name_alias"
  owner_key = "owner_key"
  owner_tag = "owner_tag"
}
`

const testConfigFvSCrtrnResetDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  annotation = "orchestrator:terraform"
  description = ""
  match = "any"
  name_alias = ""
  owner_key = ""
  owner_tag = ""
}
`
const testConfigFvSCrtrnChildrenDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  annotations = [
	{
	  key = "key_0"
	  value = "value_1"
	},
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
  tags = [
	{
	  key = "key_0"
	  value = "value_1"
	},
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
}
`

const testConfigFvSCrtrnChildrenRemoveFromConfigDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
}
`

const testConfigFvSCrtrnChildrenRemoveOneDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  annotations = [ 
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
  tags = [ 
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
}
`

const testConfigFvSCrtrnChildrenRemoveAllDependencyWithFvCrtrn = testConfigFvCrtrnMinDependencyWithFvAEPg + `
resource "aci_epg_useg_sub_block_statement" "test" {
  parent_dn = aci_epg_useg_block_statement.test.id
  name = "sub_criterion"
  annotations = []
  tags = []
}
`

const testConfigFvSCrtrnMinDependencyWithFvSCrtrnAllowExisting = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_0" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
}
resource "aci_epg_useg_sub_block_statement" "test_2" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  depends_on = [aci_epg_useg_sub_block_statement.test_0]
}
`

const testConfigFvSCrtrnMinDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
}
`

const testConfigFvSCrtrnAllDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  annotation = "annotation"
  description = "description"
  match = "all"
  name_alias = "name_alias"
  owner_key = "owner_key"
  owner_tag = "owner_tag"
}
`

const testConfigFvSCrtrnResetDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  annotation = "orchestrator:terraform"
  description = ""
  match = "any"
  name_alias = ""
  owner_key = ""
  owner_tag = ""
}
`
const testConfigFvSCrtrnChildrenDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  annotations = [
	{
	  key = "key_0"
	  value = "value_1"
	},
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
  tags = [
	{
	  key = "key_0"
	  value = "value_1"
	},
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
}
`

const testConfigFvSCrtrnChildrenRemoveFromConfigDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
}
`

const testConfigFvSCrtrnChildrenRemoveOneDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  annotations = [ 
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
  tags = [ 
	{
	  key = "key_1"
	  value = "value_2"
	},
  ]
}
`

const testConfigFvSCrtrnChildrenRemoveAllDependencyWithFvSCrtrn = testConfigFvSCrtrnMinDependencyWithFvCrtrn + `
resource "aci_epg_useg_sub_block_statement" "test_1" {
  parent_dn = aci_epg_useg_sub_block_statement.test.id
  name = "sub_criterion"
  annotations = []
  tags = []
}
`