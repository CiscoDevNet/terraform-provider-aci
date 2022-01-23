package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSpineAccessPortSelector_Basic(t *testing.T) {
	var spine_access_port_selector models.SpineAccessPortSelector
	infra_sp_acc_port_p_name := acctest.RandString(5)
	infra_sh_port_s_name := acctest.RandString(5)
	description := "spine_access_port_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineAccessPortSelectorConfig_basic(infra_sp_acc_port_p_name, infra_sh_port_s_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortSelectorExists("aci_spine_access_port_selector.foospine_access_port_selector", &spine_access_port_selector),
					testAccCheckAciSpineAccessPortSelectorAttributes(infra_sp_acc_port_p_name, infra_sh_port_s_name, description, &spine_access_port_selector),
				),
			},
		},
	})
}

func TestAccAciSpineAccessPortSelector_Update(t *testing.T) {
	var spine_access_port_selector models.SpineAccessPortSelector
	infra_sp_acc_port_p_name := acctest.RandString(5)
	infra_sh_port_s_name := acctest.RandString(5)
	description := "spine_access_port_selector created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSpineAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSpineAccessPortSelectorConfig_basic(infra_sp_acc_port_p_name, infra_sh_port_s_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortSelectorExists("aci_spine_access_port_selector.foospine_access_port_selector", &spine_access_port_selector),
					testAccCheckAciSpineAccessPortSelectorAttributes(infra_sp_acc_port_p_name, infra_sh_port_s_name, description, &spine_access_port_selector),
				),
			},
			{
				Config: testAccCheckAciSpineAccessPortSelectorConfig_basic(infra_sp_acc_port_p_name, infra_sh_port_s_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSpineAccessPortSelectorExists("aci_spine_access_port_selector.foospine_access_port_selector", &spine_access_port_selector),
					testAccCheckAciSpineAccessPortSelectorAttributes(infra_sp_acc_port_p_name, infra_sh_port_s_name, description, &spine_access_port_selector),
				),
			},
		},
	})
}

func testAccCheckAciSpineAccessPortSelectorConfig_basic(infra_sp_acc_port_p_name, infra_sh_port_s_name string) string {
	return fmt.Sprintf(`

	resource "aci_spine_interface_profile" "foospine_interface_profile" {
		name 		= "%s"
		description = "spine_interface_profile created while acceptance testing"

	}

	resource "aci_spine_access_port_selector" "foospine_access_port_selector" {
		name 		= "%s"
		description = "spine_access_port_selector created while acceptance testing"
		spine_interface_profile_dn = aci_spine_interface_profile.foospine_interface_profile.id
	}

	`, infra_sp_acc_port_p_name, infra_sh_port_s_name)
}

func testAccCheckAciSpineAccessPortSelectorExists(name string, spine_access_port_selector *models.SpineAccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Spine Access Port Selector %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spine Access Port Selector dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		spine_access_port_selectorFound := models.SpineAccessPortSelectorFromContainer(cont)
		if spine_access_port_selectorFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Spine Access Port Selector %s not found", rs.Primary.ID)
		}
		*spine_access_port_selector = *spine_access_port_selectorFound
		return nil
	}
}

func testAccCheckAciSpineAccessPortSelectorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_spine_access_port_selector" {
			cont, err := client.Get(rs.Primary.ID)
			spine_access_port_selector := models.SpineAccessPortSelectorFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Spine Access Port Selector %s Still exists", spine_access_port_selector.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSpineAccessPortSelectorAttributes(infra_sp_acc_port_p_name, infra_sh_port_s_name, description string, spine_access_port_selector *models.SpineAccessPortSelector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if infra_sh_port_s_name != GetMOName(spine_access_port_selector.DistinguishedName) {
			return fmt.Errorf("Bad infrash_port_s %s", GetMOName(spine_access_port_selector.DistinguishedName))
		}

		if infra_sp_acc_port_p_name != GetMOName(GetParentDn(spine_access_port_selector.DistinguishedName)) {
			return fmt.Errorf(" Bad infra_sp_acc_port_p %s", GetMOName(GetParentDn(spine_access_port_selector.DistinguishedName)))
		}
		if description != spine_access_port_selector.Description {
			return fmt.Errorf("Bad spine_access_port_selector Description %s", spine_access_port_selector.Description)
		}
		return nil
	}
}
