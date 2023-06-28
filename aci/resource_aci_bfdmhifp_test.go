package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBfdMultihopInterfaceProfile_Basic(t *testing.T) {
	var bfd_multihop_interface_profile models.AciBfdMultihopInterfaceProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_l_node_p_name := acctest.RandString(5)
	l3ext_l_if_p_name := acctest.RandString(5)
	bfd_mh_if_p_name := acctest.RandString(5)
	description := "bfd_multihop_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBfdMultihopInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBfdMultihopInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfaceProfileExists("aci_bfd_multihop_interface_profile.foointerface_profile", &bfd_multihop_interface_profile),
					testAccCheckAciBfdMultihopInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name, description, &bfd_multihop_interface_profile),
				),
			},
		},
	})
}

func TestAccAciBfdMultihopInterfaceProfile_Update(t *testing.T) {
	var bfd_multihop_interface_profile models.AciBfdMultihopInterfaceProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	l3ext_l_node_p_name := acctest.RandString(5)
	l3ext_l_if_p_name := acctest.RandString(5)
	bfd_mh_if_p_name := acctest.RandString(5)
	description := "bfd_multihop_interface_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBfdMultihopInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBfdMultihopInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfaceProfileExists("aci_bfd_multihop_interface_profile.foointerface_profile", &bfd_multihop_interface_profile),
					testAccCheckAciBfdMultihopInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name, description, &bfd_multihop_interface_profile),
				),
			},
			{
				Config: testAccCheckAciBfdMultihopInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfaceProfileExists("aci_bfd_multihop_interface_profile.foointerface_profile", &bfd_multihop_interface_profile),
					testAccCheckAciBfdMultihopInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name, description, &bfd_multihop_interface_profile),
				),
			},
		},
	})
}

func testAccCheckAciBfdMultihopInterfaceProfileConfig_basic(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3_outside" "fool3_outside" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	resource "aci_logical_node_profile" "foological_node_profile" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.fool3_outside.id
	}

	resource "aci_logical_interface_profile" "foological_interface_profile" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.foological_node_profile.id
	}

	resource "aci_bfd_multihop_interface_profile" "foointerface_profile" {
		name 		= "%s"
		description = "bfd_multihop_interface_profile created while acceptance testing"
		logical_interface_profile_dn = aci_logical_interface_profile.foological_interface_profile.id
	}

	`, fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name)
}

func testAccCheckAciBfdMultihopInterfaceProfileExists(name string, bfd_multihop_interface_profile *models.AciBfdMultihopInterfaceProfile) resource.TestCheckFunc {
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

		interface_profileFound := models.AciBfdMultihopInterfaceProfileFromContainer(cont)
		if interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Interface Profile %s not found", rs.Primary.ID)
		}
		*bfd_multihop_interface_profile = *interface_profileFound
		return nil
	}
}

func testAccCheckAciBfdMultihopInterfaceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bfd_multihop_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			bfd_multihop_interface_profile := models.AciBfdMultihopInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Interface Profile %s Still exists", bfd_multihop_interface_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBfdMultihopInterfaceProfileAttributes(fv_tenant_name, l3ext_out_name, l3ext_l_node_p_name, l3ext_l_if_p_name, bfd_mh_if_p_name, description string, bfd_multihop_interface_profile *models.AciBfdMultihopInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if bfd_mh_if_p_name != GetMOName(bfd_multihop_interface_profile.DistinguishedName) {
			return fmt.Errorf("Bad bfd_mh_if_p %s", GetMOName(bfd_multihop_interface_profile.DistinguishedName))
		}

		if l3ext_l_if_p_name != GetMOName(GetParentDn(bfd_multihop_interface_profile.DistinguishedName)) {
			return fmt.Errorf(" Bad l3extl_if_p %s", GetMOName(GetParentDn(bfd_multihop_interface_profile.DistinguishedName)))
		}
		if description != bfd_multihop_interface_profile.Description {
			return fmt.Errorf("Bad bfd_multihop_interface_profile Description %s", bfd_multihop_interface_profile.Description)
		}
		return nil
	}
}
