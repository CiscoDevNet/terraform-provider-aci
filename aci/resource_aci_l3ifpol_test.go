package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3InterfacePolicy_Basic(t *testing.T) {
	var l3_interface_policy models.L3InterfacePolicy
	description := "l3_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3InterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists("aci_l3_interface_policy.test", &l3_interface_policy),
					testAccCheckAciL3InterfacePolicyAttributes(description, &l3_interface_policy),
				),
			},
		},
	})
}

func TestAccAciL3InterfacePolicy_update(t *testing.T) {
	var l3_interface_policy models.L3InterfacePolicy
	description := "l3_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3InterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists("aci_l3_interface_policy.test", &l3_interface_policy),
					testAccCheckAciL3InterfacePolicyAttributes(description, &l3_interface_policy),
				),
			},
			{
				Config: testAccCheckAciL3InterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3InterfacePolicyExists("aci_l3_interface_policy.test", &l3_interface_policy),
					testAccCheckAciL3InterfacePolicyAttributes(description, &l3_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciL3InterfacePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3_interface_policy" "test" {
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		bfd_isis = "disabled"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciL3InterfacePolicyExists(name string, l3_interface_policy *models.L3InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_interface_policyFound := models.L3InterfacePolicyFromContainer(cont)
		if l3_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Interface Policy %s not found", rs.Primary.ID)
		}
		*l3_interface_policy = *l3_interface_policyFound
		return nil
	}
}

func testAccCheckAciL3InterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3_interface_policy := models.L3InterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Interface Policy %s Still exists", l3_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3InterfacePolicyAttributes(description string, l3_interface_policy *models.L3InterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3_interface_policy.Description {
			return fmt.Errorf("Bad l3_interface_policy Description %s", l3_interface_policy.Description)
		}

		if "example" != l3_interface_policy.Name {
			return fmt.Errorf("Bad l3_interface_policy name %s", l3_interface_policy.Name)
		}

		if "example" != l3_interface_policy.Annotation {
			return fmt.Errorf("Bad l3_interface_policy annotation %s", l3_interface_policy.Annotation)
		}

		if "disabled" != l3_interface_policy.BfdIsis {
			return fmt.Errorf("Bad l3_interface_policy bfd_isis %s", l3_interface_policy.BfdIsis)
		}

		if "example" != l3_interface_policy.NameAlias {
			return fmt.Errorf("Bad l3_interface_policy name_alias %s", l3_interface_policy.NameAlias)
		}

		return nil
	}
}
