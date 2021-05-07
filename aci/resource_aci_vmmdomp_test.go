package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVMMDomain_Basic(t *testing.T) {
	var vmm_domain models.VMMDomain
	description := "vmm_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVMMDomainConfig_basic(description, "epDpVerify"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists("aci_vmm_domain.foovmm_domain", &vmm_domain),
					testAccCheckAciVMMDomainAttributes(description, "epDpVerify", &vmm_domain),
				),
			},
			{
				ResourceName:      "aci_vmm_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciVMMDomain_update(t *testing.T) {
	var vmm_domain models.VMMDomain
	description := "vmm_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVMMDomainConfig_basic(description, "epDpVerify"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists("aci_vmm_domain.foovmm_domain", &vmm_domain),
					testAccCheckAciVMMDomainAttributes(description, "epDpVerify", &vmm_domain),
				),
			},
			{
				Config: testAccCheckAciVMMDomainConfig_basic(description, "none"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMDomainExists("aci_vmm_domain.foovmm_domain", &vmm_domain),
					testAccCheckAciVMMDomainAttributes(description, "none", &vmm_domain),
				),
			},
		},
	})
}

func testAccCheckAciVMMDomainConfig_basic(description, ctrl_knob string) string {
	return fmt.Sprintf(`

	resource "aci_vmm_domain" "foovmm_domain" {
		provider_profile_dn = "${aci_provider_profile.example.id}"
		description         = "%s"
		name                = "demo_domp"
		access_mode         = "read-write"
		annotation          = "tag_dom"
		arp_learning        = "disabled"
		ave_time_out        = "30"
		config_infra_pg     = "no"
		ctrl_knob           = "%s"
		delimiter           = ";"
		enable_ave          = "no"
		enable_tag          = "no"
		encap_mode          = "unknown"
		enf_pref            = "hw"
		ep_inventory_type   = "on-link"
		ep_ret_time         = "0"
		hv_avail_monitor    = "no"
		mcast_addr          = "224.0.1.0"
		mode                = "default"
		name_alias          = "alias_dom"
		pref_encap_mode     = "unspecified"
	}  
	`, description, ctrl_knob)
}

func testAccCheckAciVMMDomainExists(name string, vmm_domain *models.VMMDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_domainFound := models.VMMDomainFromContainer(cont)
		if vmm_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Domain %s not found", rs.Primary.ID)
		}
		*vmm_domain = *vmm_domainFound
		return nil
	}
}

func testAccCheckAciVMMDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vmm_domain" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_domain := models.VMMDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Domain %s Still exists", vmm_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVMMDomainAttributes(description, ctrl_knob string, vmm_domain *models.VMMDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vmm_domain.Description {
			return fmt.Errorf("Bad vmm_domain Description %s", vmm_domain.Description)
		}

		if "demo_domp" != vmm_domain.Name {
			return fmt.Errorf("Bad vmm_domain name %s", vmm_domain.Name)
		}

		if "read-write" != vmm_domain.AccessMode {
			return fmt.Errorf("Bad vmm_domain access_mode %s", vmm_domain.AccessMode)
		}

		if "tag_dom" != vmm_domain.Annotation {
			return fmt.Errorf("Bad vmm_domain annotation %s", vmm_domain.Annotation)
		}

		if "disabled" != vmm_domain.ArpLearning {
			return fmt.Errorf("Bad vmm_domain arp_learning %s", vmm_domain.ArpLearning)
		}

		if "30" != vmm_domain.AveTimeOut {
			return fmt.Errorf("Bad vmm_domain ave_time_out %s", vmm_domain.AveTimeOut)
		}

		if "no" != vmm_domain.ConfigInfraPg {
			return fmt.Errorf("Bad vmm_domain config_infra_pg %s", vmm_domain.ConfigInfraPg)
		}

		if ctrl_knob != vmm_domain.CtrlKnob {
			return fmt.Errorf("Bad vmm_domain ctrl_knob %s", vmm_domain.CtrlKnob)
		}

		if ";" != vmm_domain.Delimiter {
			return fmt.Errorf("Bad vmm_domain delimiter %s", vmm_domain.Delimiter)
		}

		if "no" != vmm_domain.EnableAVE {
			return fmt.Errorf("Bad vmm_domain enable_ave %s", vmm_domain.EnableAVE)
		}

		if "no" != vmm_domain.EnableTag {
			return fmt.Errorf("Bad vmm_domain enable_tag %s", vmm_domain.EnableTag)
		}

		if "unknown" != vmm_domain.EncapMode {
			return fmt.Errorf("Bad vmm_domain encap_mode %s", vmm_domain.EncapMode)
		}

		if "hw" != vmm_domain.EnfPref {
			return fmt.Errorf("Bad vmm_domain enf_pref %s", vmm_domain.EnfPref)
		}

		if "on-link" != vmm_domain.EpInventoryType {
			return fmt.Errorf("Bad vmm_domain ep_inventory_type %s", vmm_domain.EpInventoryType)
		}

		if "0" != vmm_domain.EpRetTime {
			return fmt.Errorf("Bad vmm_domain ep_ret_time %s", vmm_domain.EpRetTime)
		}

		if "no" != vmm_domain.HvAvailMonitor {
			return fmt.Errorf("Bad vmm_domain hv_avail_monitor %s", vmm_domain.HvAvailMonitor)
		}

		if "224.0.1.0" != vmm_domain.McastAddr {
			return fmt.Errorf("Bad vmm_domain mcast_addr %s", vmm_domain.McastAddr)
		}

		if "default" != vmm_domain.Mode {
			return fmt.Errorf("Bad vmm_domain mode %s", vmm_domain.Mode)
		}

		if "alias_dom" != vmm_domain.NameAlias {
			return fmt.Errorf("Bad vmm_domain name_alias %s", vmm_domain.NameAlias)
		}

		if "unspecified" != vmm_domain.PrefEncapMode {
			return fmt.Errorf("Bad vmm_domain pref_encap_mode %s", vmm_domain.PrefEncapMode)
		}

		return nil
	}
}
