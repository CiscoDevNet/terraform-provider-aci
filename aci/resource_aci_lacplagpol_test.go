package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciLACPPolicy_Basic(t *testing.T) {
	var lacp_policy models.LACPPolicy
	description := "lacp_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLACPPolicyConfig_basic(description, "off"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists("aci_lacp_policy.foolacp_policy", &lacp_policy),
					testAccCheckAciLACPPolicyAttributes(description, "off", &lacp_policy),
				),
			},
			{
				ResourceName:      "aci_lacp_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciLACPPolicy_update(t *testing.T) {
	var lacp_policy models.LACPPolicy
	description := "lacp_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLACPPolicyConfig_basic(description, "off"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists("aci_lacp_policy.foolacp_policy", &lacp_policy),
					testAccCheckAciLACPPolicyAttributes(description, "off", &lacp_policy),
				),
			},
			{
				Config: testAccCheckAciLACPPolicyConfig_basic(description, "active"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLACPPolicyExists("aci_lacp_policy.foolacp_policy", &lacp_policy),
					testAccCheckAciLACPPolicyAttributes(description, "active", &lacp_policy),
				),
			},
		},
	})
}

func testAccCheckAciLACPPolicyConfig_basic(description, mode string) string {
	return fmt.Sprintf(`

	resource "aci_lacp_policy" "foolacp_policy" {
		description = "%s"
		name        = "demo_lacp_pol"
		annotation  = "tag_lacp"
		ctrl        = ["susp-individual"]
		max_links   = "16"
		min_links   = "1"
		mode        = "%s"
		name_alias  = "alias_lacp"
	}
	  
	`, description, mode)
}

func testAccCheckAciLACPPolicyExists(name string, lacp_policy *models.LACPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LACP Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LACP Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		lacp_policyFound := models.LACPPolicyFromContainer(cont)
		if lacp_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LACP Policy %s not found", rs.Primary.ID)
		}
		*lacp_policy = *lacp_policyFound
		return nil
	}
}

func testAccCheckAciLACPPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_lacp_policy" {
			cont, err := client.Get(rs.Primary.ID)
			lacp_policy := models.LACPPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LACP Policy %s Still exists", lacp_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLACPPolicyAttributes(description, mode string, lacp_policy *models.LACPPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != lacp_policy.Description {
			return fmt.Errorf("Bad lacp_policy Description %s", lacp_policy.Description)
		}

		if "demo_lacp_pol" != lacp_policy.Name {
			return fmt.Errorf("Bad lacp_policy name %s", lacp_policy.Name)
		}

		if "tag_lacp" != lacp_policy.Annotation {
			return fmt.Errorf("Bad lacp_policy annotation %s", lacp_policy.Annotation)
		}

		if "16" != lacp_policy.MaxLinks {
			return fmt.Errorf("Bad lacp_policy max_links %s", lacp_policy.MaxLinks)
		}

		if "1" != lacp_policy.MinLinks {
			return fmt.Errorf("Bad lacp_policy min_links %s", lacp_policy.MinLinks)
		}

		if mode != lacp_policy.Mode {
			return fmt.Errorf("Bad lacp_policy mode %s", lacp_policy.Mode)
		}

		if "example" != lacp_policy.NameAlias {
			return fmt.Errorf("Bad lacp_policy name_alias %s", lacp_policy.NameAlias)
		}

		return nil
	}
}
