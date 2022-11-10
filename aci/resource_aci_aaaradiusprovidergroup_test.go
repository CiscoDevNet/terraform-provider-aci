package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRADIUSProviderGroup_Basic(t *testing.T) {
	var radius_provider_group models.RADIUSProviderGroup
	description := "radius_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRADIUSProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRADIUSProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderGroupExists("aci_radius_provider_group.fooradius_provider_group", &radius_provider_group),
					testAccCheckAciRADIUSProviderGroupAttributes(description, &radius_provider_group),
				),
			},
		},
	})
}

func TestAccAciRADIUSProviderGroup_Update(t *testing.T) {
	var radius_provider_group models.RADIUSProviderGroup
	description := "radius_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRADIUSProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRADIUSProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderGroupExists("aci_radius_provider_group.fooradius_provider_group", &radius_provider_group),
					testAccCheckAciRADIUSProviderGroupAttributes(description, &radius_provider_group),
				),
			},
			{
				Config: testAccCheckAciRADIUSProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderGroupExists("aci_radius_provider_group.fooradius_provider_group", &radius_provider_group),
					testAccCheckAciRADIUSProviderGroupAttributes(description, &radius_provider_group),
				),
			},
		},
	})
}

func testAccCheckAciRADIUSProviderGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_radius_provider_group" "fooradius_provider_group" {
		name 		= "test"
		description = "%s"
		name_alias  = "radius_provider_group_alias"
		annotation  = "example"

	}

	`, description)
}

func testAccCheckAciRADIUSProviderGroupExists(name string, radius_provider_group *models.RADIUSProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("RADIUS Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RADIUS Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		radius_provider_groupFound := models.RADIUSProviderGroupFromContainer(cont)
		if radius_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("RADIUS Provider Group %s not found", rs.Primary.ID)
		}
		*radius_provider_group = *radius_provider_groupFound
		return nil
	}
}

func testAccCheckAciRADIUSProviderGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_radius_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			radius_provider_group := models.RADIUSProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("RADIUS Provider Group %s Still exists", radius_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRADIUSProviderGroupAttributes(description string, radius_provider_group *models.RADIUSProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != radius_provider_group.Name {
			return fmt.Errorf("Bad aaa_radius_provider_group %s", radius_provider_group.Name)
		}

		if description != radius_provider_group.Description {
			return fmt.Errorf("Bad radius_provider_group Description %s", radius_provider_group.Description)
		}

		if "radius_provider_group_alias" != radius_provider_group.NameAlias {
			return fmt.Errorf("Bad radius_provider_group Name Alias %s", radius_provider_group.NameAlias)
		}

		if "example" != radius_provider_group.Annotation {
			return fmt.Errorf("Bad radius_provider_group Annotation %s", radius_provider_group.Annotation)
		}
		return nil
	}
}
