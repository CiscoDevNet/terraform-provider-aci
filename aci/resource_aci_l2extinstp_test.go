package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL2outExternalEpg_Basic(t *testing.T) {
	var l2out_extepg models.L2outExternalEpg
	description := "l2out_extepg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2outExternalEpgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists("aci_l2out_extepg.fool2out_extepg", &l2out_extepg),
					testAccCheckAciL2outExternalEpgAttributes(description, &l2out_extepg),
				),
			},
		},
	})
}

func TestAccAciL2outExternalEpg_update(t *testing.T) {
	var l2out_extepg models.L2outExternalEpg
	description := "l2out_extepg created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2outExternalEpgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2outExternalEpgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists("aci_l2out_extepg.fool2out_extepg", &l2out_extepg),
					testAccCheckAciL2outExternalEpgAttributes(description, &l2out_extepg),
				),
			},
			{
				Config: testAccCheckAciL2outExternalEpgConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2outExternalEpgExists("aci_l2out_extepg.fool2out_extepg", &l2out_extepg),
					testAccCheckAciL2outExternalEpgAttributes(description, &l2out_extepg),
				),
			},
		},
	})
}

func testAccCheckAciL2outExternalEpgConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l2out_extepg" "fool2out_extepg" {
		l2_outside_dn  = "${aci_l2_outside.example.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		exception_tag  = "example"
  		flood_on_encap = "disabled"
  		match_t = "All"
  		name_alias  = "example"
  		pref_gr_memb = "exclude"
  		prio = "level1"
  		target_dscp = "AF11"
	}
	`, description)
}

func testAccCheckAciL2outExternalEpgExists(name string, l2out_extepg *models.L2outExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2-Out External EPG %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2-Out External EPG dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2out_extepgFound := models.L2outExternalEpgFromContainer(cont)
		if l2out_extepgFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2-Out External EPG %s not found", rs.Primary.ID)
		}
		*l2out_extepg = *l2out_extepgFound
		return nil
	}
}

func testAccCheckAciL2outExternalEpgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l2out_extepg" {
			cont, err := client.Get(rs.Primary.ID)
			l2out_extepg := models.L2outExternalEpgFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2-Out External EPG %s Still exists", l2out_extepg.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL2outExternalEpgAttributes(description string, l2out_extepg *models.L2outExternalEpg) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l2out_extepg.Description {
			return fmt.Errorf("Bad l2out_extepg Description %s", l2out_extepg.Description)
		}

		if "example" != l2out_extepg.Name {
			return fmt.Errorf("Bad l2out_extepg name %s", l2out_extepg.Name)
		}

		if "example" != l2out_extepg.Annotation {
			return fmt.Errorf("Bad l2out_extepg annotation %s", l2out_extepg.Annotation)
		}

		if "example" != l2out_extepg.ExceptionTag {
			return fmt.Errorf("Bad l2out_extepg exception_tag %s", l2out_extepg.ExceptionTag)
		}

		if "disabled" != l2out_extepg.FloodOnEncap {
			return fmt.Errorf("Bad l2out_extepg flood_on_encap %s", l2out_extepg.FloodOnEncap)
		}

		if "All" != l2out_extepg.MatchT {
			return fmt.Errorf("Bad l2out_extepg match_t %s", l2out_extepg.MatchT)
		}

		if "example" != l2out_extepg.NameAlias {
			return fmt.Errorf("Bad l2out_extepg name_alias %s", l2out_extepg.NameAlias)
		}

		if "exclude" != l2out_extepg.PrefGrMemb {
			return fmt.Errorf("Bad l2out_extepg pref_gr_memb %s", l2out_extepg.PrefGrMemb)
		}

		if "level1" != l2out_extepg.Prio {
			return fmt.Errorf("Bad l2out_extepg prio %s", l2out_extepg.Prio)
		}

		if "AF11" != l2out_extepg.TargetDscp {
			return fmt.Errorf("Bad l2out_extepg target_dscp %s", l2out_extepg.TargetDscp)
		}

		return nil
	}
}
