package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciHSRPInterfacePolicy_Basic(t *testing.T) {
	var hsrp_interface_policy models.HSRPInterfacePolicy
	description := "hsrp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPInterfacePolicyExists("aci_hsrp_interface_policy.test", &hsrp_interface_policy),
					testAccCheckAciHSRPInterfacePolicyAttributes(description, &hsrp_interface_policy),
				),
			},
		},
	})
}

func TestAccAciHSRPInterfacePolicy_update(t *testing.T) {
	var hsrp_interface_policy models.HSRPInterfacePolicy
	description := "hsrp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPInterfacePolicyExists("aci_hsrp_interface_policy.test", &hsrp_interface_policy),
					testAccCheckAciHSRPInterfacePolicyAttributes(description, &hsrp_interface_policy),
				),
			},
			{
				Config: testAccCheckAciHSRPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPInterfacePolicyExists("aci_hsrp_interface_policy.test", &hsrp_interface_policy),
					testAccCheckAciHSRPInterfacePolicyAttributes(description, &hsrp_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciHSRPInterfacePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn    = "uni/tn-aaaaa"
		name         = "one"
		annotation   = "example"
		description  = "%s"
		ctrl         = "bia"
		delay        = "10"
		name_alias   = "example"
		reload_delay = "10"
	  }
	`, description)
}

func testAccCheckAciHSRPInterfacePolicyExists(name string, hsrp_interface_policy *models.HSRPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("HSRP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No HSRP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		hsrp_interface_policyFound := models.HSRPInterfacePolicyFromContainer(cont)
		if hsrp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("HSRP Interface Policy %s not found", rs.Primary.ID)
		}
		*hsrp_interface_policy = *hsrp_interface_policyFound
		return nil
	}
}

func testAccCheckAciHSRPInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_hsrp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			hsrp_interface_policy := models.HSRPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("HSRP Interface Policy %s Still exists", hsrp_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciHSRPInterfacePolicyAttributes(description string, hsrp_interface_policy *models.HSRPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != hsrp_interface_policy.Description {
			return fmt.Errorf("Bad hsrp_interface_policy Description %s", hsrp_interface_policy.Description)
		}

		if "one" != hsrp_interface_policy.Name {
			return fmt.Errorf("Bad hsrp_interface_policy name %s", hsrp_interface_policy.Name)
		}

		if "example" != hsrp_interface_policy.Annotation {
			return fmt.Errorf("Bad hsrp_interface_policy annotation %s", hsrp_interface_policy.Annotation)
		}

		if "bia" != hsrp_interface_policy.Ctrl {
			return fmt.Errorf("Bad hsrp_interface_policy ctrl %s", hsrp_interface_policy.Ctrl)
		}

		if "10" != hsrp_interface_policy.Delay {
			return fmt.Errorf("Bad hsrp_interface_policy delay %s", hsrp_interface_policy.Delay)
		}

		if "example" != hsrp_interface_policy.NameAlias {
			return fmt.Errorf("Bad hsrp_interface_policy name_alias %s", hsrp_interface_policy.NameAlias)
		}

		if "10" != hsrp_interface_policy.ReloadDelay {
			return fmt.Errorf("Bad hsrp_interface_policy reload_delay %s", hsrp_interface_policy.ReloadDelay)
		}

		return nil
	}
}
