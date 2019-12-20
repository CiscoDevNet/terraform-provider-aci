package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciVRF_Basic(t *testing.T) {
	var vrf models.VRF
	description := "vrf created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVRFConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists("aci_vrf.foovrf", &vrf),
					testAccCheckAciVRFAttributes(description, "enabled", &vrf),
				),
			},
			{
				ResourceName:      "aci_vrf",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciVRF_update(t *testing.T) {
	var vrf models.VRF
	description := "vrf created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVRFDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVRFConfig_basic(description, "enabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists("aci_vrf.foovrf", &vrf),
					testAccCheckAciVRFAttributes(description, "enabled", &vrf),
				),
			},
			{
				Config: testAccCheckAciVRFConfig_basic(description, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFExists("aci_vrf.foovrf", &vrf),
					testAccCheckAciVRFAttributes(description, "disabled", &vrf),
				),
			},
		},
	})
}

func testAccCheckAciVRFConfig_basic(description, ip_data_plane_learning string) string {
	return fmt.Sprintf(`
	
	resource "aci_tenant" "tenant_for_vrf" {
		name        = "tenant_for_vrf"
		description = "This tenant is created by terraform ACI provider"
	}

	resource "aci_vrf" "foovrf" {
		tenant_dn   		   = "${aci_tenant.tenant_for_vrf.id}"
		description 		   = "%s"
		name                   = "demo_vrf"
		annotation             = "tag_vrf"
		bd_enforced_enable     = "no"
		ip_data_plane_learning = "%s"
		knw_mcast_act          = "permit"
		name_alias             = "alias_vrf"
		pc_enf_dir             = "egress"
		pc_enf_pref            = "unenforced"
	}
	  
	`, description, ip_data_plane_learning)
}

func testAccCheckAciVRFExists(name string, vrf *models.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VRF %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VRF dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vrfFound := models.VRFFromContainer(cont)
		if vrfFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VRF %s not found", rs.Primary.ID)
		}
		*vrf = *vrfFound
		return nil
	}
}

func testAccCheckAciVRFDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vrf" {
			cont, err := client.Get(rs.Primary.ID)
			vrf := models.VRFFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VRF %s Still exists", vrf.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVRFAttributes(description, ip_data_plane_learning string, vrf *models.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vrf.Description {
			return fmt.Errorf("Bad vrf Description %s", vrf.Description)
		}

		if "demo_vrf" != vrf.Name {
			return fmt.Errorf("Bad vrf name %s", vrf.Name)
		}

		if "tag_vrf" != vrf.Annotation {
			return fmt.Errorf("Bad vrf annotation %s", vrf.Annotation)
		}

		if "no" != vrf.BdEnforcedEnable {
			return fmt.Errorf("Bad vrf bd_enforced_enable %s", vrf.BdEnforcedEnable)
		}

		if ip_data_plane_learning != vrf.IpDataPlaneLearning {
			return fmt.Errorf("Bad vrf ip_data_plane_learning %s", vrf.IpDataPlaneLearning)
		}

		if "permit" != vrf.KnwMcastAct {
			return fmt.Errorf("Bad vrf knw_mcast_act %s", vrf.KnwMcastAct)
		}

		if "alias_vrf" != vrf.NameAlias {
			return fmt.Errorf("Bad vrf name_alias %s", vrf.NameAlias)
		}

		if "egress" != vrf.PcEnfDir {
			return fmt.Errorf("Bad vrf pc_enf_dir %s", vrf.PcEnfDir)
		}

		if "unenforced" != vrf.PcEnfPref {
			return fmt.Errorf("Bad vrf pc_enf_pref %s", vrf.PcEnfPref)
		}

		return nil
	}
}
