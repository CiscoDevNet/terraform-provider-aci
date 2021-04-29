package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCDPInterfacePolicy_Basic(t *testing.T) {
	var cdp_interface_policy models.CDPInterfacePolicy
	description := "cdp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCDPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists("aci_cdp_interface_policy.foocdp_interface_policy", &cdp_interface_policy),
					testAccCheckAciCDPInterfacePolicyAttributes(description, &cdp_interface_policy),
				),
			},
			{
				ResourceName:      "aci_cdp_interface_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCDPInterfacePolicy_update(t *testing.T) {
	var cdp_interface_policy models.CDPInterfacePolicy
	description := "cdp_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCDPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists("aci_cdp_interface_policy.foocdp_interface_policy", &cdp_interface_policy),
					testAccCheckAciCDPInterfacePolicyAttributes(description, &cdp_interface_policy),
				),
			},
			{
				Config: testAccCheckAciCDPInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCDPInterfacePolicyExists("aci_cdp_interface_policy.foocdp_interface_policy", &cdp_interface_policy),
					testAccCheckAciCDPInterfacePolicyAttributes(description, &cdp_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciCDPInterfacePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_cdp_interface_policy" "foocdp_interface_policy" {
		description = "%s"
		
		name  = "example"
		  admin_st  = "enabled"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciCDPInterfacePolicyExists(name string, cdp_interface_policy *models.CDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("CDP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No CDP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cdp_interface_policyFound := models.CDPInterfacePolicyFromContainer(cont)
		if cdp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("CDP Interface Policy %s not found", rs.Primary.ID)
		}
		*cdp_interface_policy = *cdp_interface_policyFound
		return nil
	}
}

func testAccCheckAciCDPInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cdp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			cdp_interface_policy := models.CDPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("CDP Interface Policy %s Still exists", cdp_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCDPInterfacePolicyAttributes(description string, cdp_interface_policy *models.CDPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cdp_interface_policy.Description {
			return fmt.Errorf("Bad cdp_interface_policy Description %s", cdp_interface_policy.Description)
		}

		if "example" != cdp_interface_policy.Name {
			return fmt.Errorf("Bad cdp_interface_policy name %s", cdp_interface_policy.Name)
		}

		if "enabled" != cdp_interface_policy.AdminSt {
			return fmt.Errorf("Bad cdp_interface_policy admin_st %s", cdp_interface_policy.AdminSt)
		}

		if "example" != cdp_interface_policy.Annotation {
			return fmt.Errorf("Bad cdp_interface_policy annotation %s", cdp_interface_policy.Annotation)
		}

		if "example" != cdp_interface_policy.NameAlias {
			return fmt.Errorf("Bad cdp_interface_policy name_alias %s", cdp_interface_policy.NameAlias)
		}

		return nil
	}
}
