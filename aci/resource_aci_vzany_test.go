package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciAny_Basic(t *testing.T) {
	var any models.Any
	description := "any created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAnyConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists("aci_any.fooany", &any),
					testAccCheckAciAnyAttributes(description, "AtleastOne", &any),
				),
			},
			{
				ResourceName:      "aci_any",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciAny_update(t *testing.T) {
	var any models.Any
	description := "any created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAnyConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists("aci_any.fooany", &any),
					testAccCheckAciAnyAttributes(description, "AtleastOne", &any),
				),
			},
			{
				Config: testAccCheckAciAnyConfig_basic(description, "AtmostOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAnyExists("aci_any.fooany", &any),
					testAccCheckAciAnyAttributes(description, "AtmostOne", &any),
				),
			},
		},
	})
}

func testAccCheckAciAnyConfig_basic(description, match_t string) string {
	return fmt.Sprintf(`

	resource "aci_any" "fooany" {
		vrf_dn       = "${aci_vrf.example.id}"
		description  = "%s"
		annotation   = "tag_any"
		match_t      = "%s"
		name_alias   = "alias_any"
		pref_gr_memb = "disabled"
	}
	  
	`, description, match_t)
}

func testAccCheckAciAnyExists(name string, any *models.Any) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Any %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Any dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		anyFound := models.AnyFromContainer(cont)
		if anyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Any %s not found", rs.Primary.ID)
		}
		*any = *anyFound
		return nil
	}
}

func testAccCheckAciAnyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_any" {
			cont, err := client.Get(rs.Primary.ID)
			any := models.AnyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Any %s Still exists", any.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAnyAttributes(description, match_t string, any *models.Any) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != any.Description {
			return fmt.Errorf("Bad any Description %s", any.Description)
		}

		if "tag_any" != any.Annotation {
			return fmt.Errorf("Bad any annotation %s", any.Annotation)
		}

		if match_t != any.MatchT {
			return fmt.Errorf("Bad any match_t %s", any.MatchT)
		}

		if "alias_any" != any.NameAlias {
			return fmt.Errorf("Bad any name_alias %s", any.NameAlias)
		}

		if "disabled" != any.PrefGrMemb {
			return fmt.Errorf("Bad any pref_gr_memb %s", any.PrefGrMemb)
		}

		return nil
	}
}
