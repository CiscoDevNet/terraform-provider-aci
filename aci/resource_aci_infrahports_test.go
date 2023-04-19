package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAccessPortSelector_Basic(t *testing.T) {
	var access_port_selector models.AccessPortSelector
	description := "access_port_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessPortSelectorConfig_basic(description, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists("aci_access_port_selector.fooaccess_port_selector", &access_port_selector),
					testAccCheckAciAccessPortSelectorAttributes(description, "ALL", &access_port_selector),
				),
			},
			{
				ResourceName:      "aci_access_port_selector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciAccessPortSelector_update(t *testing.T) {
	var access_port_selector models.AccessPortSelector
	description := "access_port_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciAccessPortSelectorConfig_basic(description, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists("aci_access_port_selector.fooaccess_port_selector", &access_port_selector),
					testAccCheckAciAccessPortSelectorAttributes(description, "ALL", &access_port_selector),
				),
			},
			{
				Config: testAccCheckAciAccessPortSelectorConfig_basic(description, "range"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAccessPortSelectorExists("aci_access_port_selector.fooaccess_port_selector", &access_port_selector),
					testAccCheckAciAccessPortSelectorAttributes(description, "range", &access_port_selector),
				),
			},
		},
	})
}

func testAccCheckAciAccessPortSelectorConfig_basic(description, access_port_selector_type string) string {
	return fmt.Sprintf(`

	resource "aci_leaf_interface_profile" "example" {
		name        = "demo_leaf_profile"
	}	
	resource "aci_access_port_selector" "fooaccess_port_selector" {
		leaf_interface_profile_dn = aci_leaf_interface_profile.example.id
		description               = "%s"
		name                      = "demo_port_selector"
		access_port_selector_type = "%s"
		annotation                = "tag_port_selector"
		name_alias                = "alias_port_selector"
	} 
	`, description, access_port_selector_type)
}

func testAccCheckAciAccessPortSelectorExists(name string, access_port_selector *models.AccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Access Port Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Access Port Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		access_port_selectorFound := models.AccessPortSelectorFromContainer(cont)
		if access_port_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Access Port Selector %s not found", rs.Primary.ID)
		}
		*access_port_selector = *access_port_selectorFound
		return nil
	}
}

func testAccCheckAciAccessPortSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_access_port_selector" {
			cont, err := client.Get(rs.Primary.ID)
			access_port_selector := models.AccessPortSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Access Port Selector %s Still exists", access_port_selector.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciAccessPortSelectorAttributes(description, access_port_selector_type string, access_port_selector *models.AccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != access_port_selector.Description {
			return fmt.Errorf("Bad access_port_selector Description %s", access_port_selector.Description)
		}

		if "demo_port_selector" != access_port_selector.Name {
			return fmt.Errorf("Bad access_port_selector name %s", access_port_selector.Name)
		}

		if access_port_selector_type != access_port_selector.AccessPortSelector_type {
			return fmt.Errorf("Bad access_port_selector access_port_selector_type %s", access_port_selector.AccessPortSelector_type)
		}

		if "tag_port_selector" != access_port_selector.Annotation {
			return fmt.Errorf("Bad access_port_selector annotation %s", access_port_selector.Annotation)
		}

		if "alias_port_selector" != access_port_selector.NameAlias {
			return fmt.Errorf("Bad access_port_selector name_alias %s", access_port_selector.NameAlias)
		}

		return nil
	}
}
