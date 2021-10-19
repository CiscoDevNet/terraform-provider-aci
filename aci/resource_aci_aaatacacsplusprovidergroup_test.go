package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTACACSPlusProviderGroup_Basic(t *testing.T) {
	var tacacs_provider_group models.TACACSPlusProviderGroup
	description := "tacacs_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSPlusProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists("aci_tacacs_provider_group.footacacs_provider_group", &tacacs_provider_group),
					testAccCheckAciTACACSPlusProviderGroupAttributes(description, &tacacs_provider_group),
				),
			},
		},
	})
}

func TestAccAciTACACSPlusProviderGroup_Update(t *testing.T) {
	var tacacs_provider_group models.TACACSPlusProviderGroup
	description := "tacacs_provider_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSPlusProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists("aci_tacacs_provider_group.footacacs_provider_group", &tacacs_provider_group),
					testAccCheckAciTACACSPlusProviderGroupAttributes(description, &tacacs_provider_group),
				),
			},
			{
				Config: testAccCheckAciTACACSPlusProviderGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSPlusProviderGroupExists("aci_tacacs_provider_group.footacacs_provider_group", &tacacs_provider_group),
					testAccCheckAciTACACSPlusProviderGroupAttributes(description, &tacacs_provider_group),
				),
			},
		},
	})
}

func testAccCheckAciTACACSPlusProviderGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tacacs_provider_group" "footacacs_provider_group" {
		name = "test"
		description = "%s"
		annotation = "test_annotation_value"
		name_alias = "test_name_alias"
	}

	`, description)
}

func testAccCheckAciTACACSPlusProviderGroupExists(name string, tacacs_provider_group *models.TACACSPlusProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACSPlus Provider Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACSPlus Provider Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tacacs_provider_groupFound := models.TACACSPlusProviderGroupFromContainer(cont)
		if tacacs_provider_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACSPlus Provider Group %s not found", rs.Primary.ID)
		}
		*tacacs_provider_group = *tacacs_provider_groupFound
		return nil
	}
}

func testAccCheckAciTACACSPlusProviderGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_provider_group" {
			cont, err := client.Get(rs.Primary.ID)
			tacacs_provider_group := models.TACACSPlusProviderGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACSPlus Provider Group %s Still exists", tacacs_provider_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSPlusProviderGroupAttributes(description string, tacacs_provider_group *models.TACACSPlusProviderGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(tacacs_provider_group.DistinguishedName) {
			return fmt.Errorf("Bad tacacs_provider_group %s", GetMOName(tacacs_provider_group.DistinguishedName))
		}

		if description != tacacs_provider_group.Description {
			return fmt.Errorf("Bad tacacs_provider_group Description %s", tacacs_provider_group.Description)
		}

		if "test_annotation_value" != tacacs_provider_group.Annotation {
			return fmt.Errorf("Bad tacacs_provider_group Annotation %s", tacacs_provider_group.Annotation)
		}

		if "test_name_alias" != tacacs_provider_group.NameAlias {
			return fmt.Errorf("Bad tacacs_provider_group NameAlias %s", tacacs_provider_group.NameAlias)
		}
		return nil
	}
}
