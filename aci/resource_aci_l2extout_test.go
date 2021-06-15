package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL2Outside_Basic(t *testing.T) {
	var l2_outside models.L2Outside
	description := "l2_outside created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2OutsideConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2OutsideExists("aci_l2_outside.fool2_outside", &l2_outside),
					testAccCheckAciL2OutsideAttributes(description, &l2_outside),
				),
			},
		},
	})
}

func TestAccAciL2Outside_update(t *testing.T) {
	var l2_outside models.L2Outside
	description := "l2_outside created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2OutsideConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2OutsideExists("aci_l2_outside.fool2_outside", &l2_outside),
					testAccCheckAciL2OutsideAttributes(description, &l2_outside),
				),
			},
			{
				Config: testAccCheckAciL2OutsideConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2OutsideExists("aci_l2_outside.fool2_outside", &l2_outside),
					testAccCheckAciL2OutsideAttributes(description, &l2_outside),
				),
			},
		},
	})
}

func testAccCheckAciL2OutsideConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name        = "tenant_1"
		description = "This tenant is created by terraform"
	}
	
	resource "aci_l2_outside" "fool2_outside" {
		tenant_dn  = aci_tenant.foo_tenant.id
		description = "%s"
		name  = "l2_outside_1"
  		annotation  = "example"
  		name_alias  = "example"
  		target_dscp = "AF11"
	}
	`, description)
}

func testAccCheckAciL2OutsideExists(name string, l2_outside *models.L2Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2 Outside %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2 Outside dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_outsideFound := models.L2OutsideFromContainer(cont)
		if l2_outsideFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2 Outside %s not found", rs.Primary.ID)
		}
		*l2_outside = *l2_outsideFound
		return nil
	}
}

func testAccCheckAciL2OutsideDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l2_outside" {
			cont, err := client.Get(rs.Primary.ID)
			l2_outside := models.L2OutsideFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2 Outside %s Still exists", l2_outside.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL2OutsideAttributes(description string, l2_outside *models.L2Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l2_outside.Description {
			return fmt.Errorf("Bad l2_outside Description %s", l2_outside.Description)
		}

		if "l2_outside_1" != l2_outside.Name {
			return fmt.Errorf("Bad l2_outside name %s", l2_outside.Name)
		}

		if "example" != l2_outside.Annotation {
			return fmt.Errorf("Bad l2_outside annotation %s", l2_outside.Annotation)
		}

		if "example" != l2_outside.NameAlias {
			return fmt.Errorf("Bad l2_outside name_alias %s", l2_outside.NameAlias)
		}

		if "AF11" != l2_outside.TargetDscp {
			return fmt.Errorf("Bad l2_outside target_dscp %s", l2_outside.TargetDscp)
		}

		return nil
	}
}
