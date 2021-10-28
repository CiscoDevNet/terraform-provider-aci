package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTACACSSource_Basic(t *testing.T) {
	var tacacs_source models.TACACSSource
	description := "tacacs_source created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSSourceConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists("aci_tacacs_source.footacacs_source", &tacacs_source),
					testAccCheckAciTACACSSourceAttributes(description, &tacacs_source),
				),
			},
		},
	})
}

func TestAccAciTACACSSource_Update(t *testing.T) {
	var tacacs_source models.TACACSSource
	description := "tacacs_source created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSSourceConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists("aci_tacacs_source.footacacs_source", &tacacs_source),
					testAccCheckAciTACACSSourceAttributes(description, &tacacs_source),
				),
			},
			{
				Config: testAccCheckAciTACACSSourceConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSSourceExists("aci_tacacs_source.footacacs_source", &tacacs_source),
					testAccCheckAciTACACSSourceAttributes(description, &tacacs_source),
				),
			},
		},
	})
}

func testAccCheckAciTACACSSourceConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tacacs_source" "footacacs_source" {
		name 		= "test"
		parent_dn   = "uni/tn-aaaaa/monepg-any"
		description = "%s"
		name_alias	= "tacacs_source_alias"
		incl 		= ["audit"]
		min_sev		= "major"
	}

	`, description)
}

func testAccCheckAciTACACSSourceExists(name string, tacacs_source *models.TACACSSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Source %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Source dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_sourceFound := models.TACACSSourceFromContainer(cont)
		if tacacs_sourceFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Source %s not found", rs.Primary.ID)
		}
		*tacacs_source = *tacacs_sourceFound
		return nil
	}
}

func testAccCheckAciTACACSSourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_source" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_source := models.TACACSSourceFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Source %s Still exists", tacacs_source.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSSourceAttributes(description string, tacacs_source *models.TACACSSource) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != tacacs_source.Name {
			return fmt.Errorf("Bad tacacs_src %s", tacacs_source.Name)
		}
		if description != tacacs_source.Description {
			return fmt.Errorf("Bad tacacs_source Description %s", tacacs_source.Description)
		}
		if "tacacs_source_alias" != tacacs_source.NameAlias {
			return fmt.Errorf("Bad tacacs_source Name Alias %s", tacacs_source.NameAlias)
		}
		if "major" != tacacs_source.MinSev {
			return fmt.Errorf("Bad tacacs_source Min Sev %s", tacacs_source.MinSev)
		}
		return nil
	}
}
