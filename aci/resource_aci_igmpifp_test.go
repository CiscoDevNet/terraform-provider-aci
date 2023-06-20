package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciInterfaceProfile_Basic(t *testing.T) {
	var interface_profile models.InterfaceProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_lnode_p_name := acctest.RandString(5)
	l3ext_lif_p_name := acctest.RandString(5)
	igmp_if_p_name := acctest.RandString(5)
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_interface_profile.foo_interface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name, description, &interface_profile),
				),
			},
		},
	})
}

func TestAccAciInterfaceProfile_Update(t *testing.T) {
	var interface_profile models.InterfaceProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_lnode_p_name := acctest.RandString(5)
	l3ext_lif_p_name := acctest.RandString(5)
	igmp_if_p_name := acctest.RandString(5)
	description := "interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_interface_profile.foo_interface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name, description, &interface_profile),
				),
			},
			{
				Config: testAccCheckAciInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInterfaceProfileExists("aci_interface_profile.foo_interface_profile", &interface_profile),
					testAccCheckAciInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name, description, &interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3outside" "foo_l3outside" {
		name 		= "%s"
		description = "l3outside created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_logical_node_profile" "foo_logical_node_profile" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3outside_dn = aci_l3outside.foo_l3outside.id
	}

	resource "aci_logical_interface_profile" "foo_logical_interface_profile" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.foo_logical_node_profile.id
	}

	resource "aci_interface_profile" "foo_interface_profile" {
		name 		= "%s"
		description = "interface_profile created while acceptance testing"
		logical_interface_profile_dn = aci_logical_interface_profile.foo_logical_interface_profile.id
	}

	`, fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name)
}

func testAccCheckAciInterfaceProfileExists(name string, interface_profile *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		interface_profileFound := models.InterfaceProfileFromContainer(cont)
		if interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Profile %s not found", rs.Primary.ID)
		}
		*interface_profile = *interface_profileFound
		return nil
	}
}

func testAccCheckAciInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			interface_profile := models.InterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Profile %s Still exists", interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_lnode_p_name, l3ext_lif_p_name, igmp_if_p_name, description string, interface_profile *models.InterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if igmp_if_p_name != GetMOName(interface_profile.DistinguishedName) {
			return fmt.Errorf("Bad igmpif_p %s", GetMOName(interface_profile.DistinguishedName))
		}

		if l3ext_lif_p_name != GetMOName(GetParentDn(interface_profile.DistinguishedName, interface_profile.Rn)) {
			return fmt.Errorf(" Bad l3extlif_p %s", GetMOName(GetParentDn(interface_profile.DistinguishedName, interface_profile.Rn)))
		}
		if description != interface_profile.Description {
			return fmt.Errorf("Bad interface_profile Description %s", interface_profile.Description)
		}
		return nil
	}
}
