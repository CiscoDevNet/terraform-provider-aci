package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciSpanningTreeInterfacePolicy_Basic(t *testing.T) {
	var spanning_tree_interface_policy models.SpanningTreeInterfacePolicy
	description := "Spanning Tree Interface Policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists("aci_stp_if_pol.foospanning_tree_interface_policy", &spanning_tree_interface_policy),
					testAccCheckAciSpanningTreeInterfacePolicyAttributes(description, &spanning_tree_interface_policy),
				),
			},
		},
	})
}

func TestAccAciSpanningTreeInterfacePolicy_update(t *testing.T) {
	var spanning_tree_interface_policy models.SpanningTreeInterfacePolicy
	description := "Spanning Tree Interface Policy created while acceptance testing"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpanningTreeInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists("aci_stp_if_pol.foospanning_tree_interface_policy", &spanning_tree_interface_policy),
					testAccCheckAciSpanningTreeInterfacePolicyAttributes(description, &spanning_tree_interface_policy),
				),
			},
			{
				Config: testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpanningTreeInterfacePolicyExists("aci_stp_if_pol.foospanning_tree_interface_policy", &spanning_tree_interface_policy),
					testAccCheckAciSpanningTreeInterfacePolicyAttributes(description, &spanning_tree_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciSpanningTreeInterfacePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_stp_if_pol" "foospanning_tree_interface_policy" {
		description = "%s"
		name  = "example"
		annotation  = "example"
		ctrl  = ["bpdu-guard"]
		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSpanningTreeInterfacePolicyExists(name string, spanning_tree_interface_policy *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spanning Tree Interface Policy dn was set")
		}
		client := testAccProvider.Meta().(*client.Client)
		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		spanning_tree_interface_policyFound := models.SpanningTreeInterfacePolicyFromContainer(cont)
		if spanning_tree_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spanning Tree Interface Policy %s not found", rs.Primary.ID)
		}
		*spanning_tree_interface_policy = *spanning_tree_interface_policyFound
		return nil
	}
}

func testAccCheckAciSpanningTreeInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_stp_if_pol" {
			cont, err := client.Get(rs.Primary.ID)
			spanning_tree_interface_policy := models.SpanningTreeInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spanning Tree Interface Policy %s Still exists", spanning_tree_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpanningTreeInterfacePolicyAttributes(description string, spanning_tree_interface_policy *models.SpanningTreeInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != spanning_tree_interface_policy.Description {
			return fmt.Errorf("Bad Spanning Tree Interface Policy Description %s", spanning_tree_interface_policy.Description)
		}
		if "example" != spanning_tree_interface_policy.Name {
			return fmt.Errorf("Bad Spanning Tree Interface Policy name %s", spanning_tree_interface_policy.Name)
		}
		if "example" != spanning_tree_interface_policy.Annotation {
			return fmt.Errorf("Bad Spanning Tree Interface Policy annotation %s", spanning_tree_interface_policy.Annotation)
		}
		if "bpdu-guard" != spanning_tree_interface_policy.Ctrl {
			return fmt.Errorf("Bad Spanning Tree Interface Policy ctrl %s", spanning_tree_interface_policy.Ctrl)
		}
		if "example" != spanning_tree_interface_policy.NameAlias {
			return fmt.Errorf("Bad Spanning Tree Interface Policy name_alias %s", spanning_tree_interface_policy.NameAlias)
		}
		return nil
	}
}
