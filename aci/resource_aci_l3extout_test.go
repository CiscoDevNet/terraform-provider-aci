package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3Outside_Basic(t *testing.T) {
	var l3_outside models.L3Outside
	description := "l3_outside created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3OutsideConfig_basic(description, "export"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3OutsideExists("aci_l3_outside.fool3_outside", &l3_outside),
					testAccCheckAciL3OutsideAttributes(description, "export", &l3_outside),
				),
			},
		},
	})
}

func TestAccAciL3Outside_update(t *testing.T) {
	var l3_outside models.L3Outside
	description := "l3_outside created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3OutsideConfig_basic(description, "export"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3OutsideExists("aci_l3_outside.fool3_outside", &l3_outside),
					testAccCheckAciL3OutsideAttributes(description, "export", &l3_outside),
				),
			},
			{
				Config: testAccCheckAciL3OutsideConfig_basic(description, "export"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3OutsideExists("aci_l3_outside.fool3_outside", &l3_outside),
					testAccCheckAciL3OutsideAttributes(description, "export", &l3_outside),
				),
			},
		},
	})
}

func testAccCheckAciL3OutsideConfig_basic(description, enforce_rtctrl string) string {
	return fmt.Sprintf(`

	resource "aci_l3_outside" "fool3_outside" {
		tenant_dn      = "uni/tn-crest_test_kishan_tenant"
		description    = "%s"
		name           = "demo_l3out"
		annotation     = "tag_l3out"
		enforce_rtctrl = "%s"
		name_alias     = "alias_out"
		target_dscp    = "unspecified"
	}
	  
	`, description, enforce_rtctrl)
}

func testAccCheckAciL3OutsideExists(name string, l3_outside *models.L3Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Outside %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Outside dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_outsideFound := models.L3OutsideFromContainer(cont)
		if l3_outsideFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Outside %s not found", rs.Primary.ID)
		}
		*l3_outside = *l3_outsideFound
		return nil
	}
}

func testAccCheckAciL3OutsideDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3_outside" {
			cont, err := client.Get(rs.Primary.ID)
			l3_outside := models.L3OutsideFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Outside %s Still exists", l3_outside.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3OutsideAttributes(description, enforce_rtctrl string, l3_outside *models.L3Outside) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3_outside.Description {
			return fmt.Errorf("Bad l3_outside Description %s", l3_outside.Description)
		}

		if "demo_l3out" != l3_outside.Name {
			return fmt.Errorf("Bad l3_outside name %s", l3_outside.Name)
		}

		if "tag_l3out" != l3_outside.Annotation {
			return fmt.Errorf("Bad l3_outside annotation %s", l3_outside.Annotation)
		}

		if enforce_rtctrl != l3_outside.EnforceRtctrl {
			return fmt.Errorf("Bad l3_outside enforce_rtctrl %s", l3_outside.EnforceRtctrl)
		}

		if "alias_out" != l3_outside.NameAlias {
			return fmt.Errorf("Bad l3_outside name_alias %s", l3_outside.NameAlias)
		}

		if "unspecified" != l3_outside.TargetDscp {
			return fmt.Errorf("Bad l3_outside target_dscp %s", l3_outside.TargetDscp)
		}

		return nil
	}
}
