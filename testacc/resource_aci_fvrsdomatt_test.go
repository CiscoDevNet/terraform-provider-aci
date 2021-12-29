package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFVDomain_Basic(t *testing.T) {
	var epg_to_domain_default models.FVDomain
	var epg_to_domain_updated models.FVDomain
	resourceName := "aci_epg_to_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	rother := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFVDomainWithoutTdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFVDomainWithoutApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFVDomainConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "binding_type", "none"),
					resource.TestCheckResourceAttr(resourceName, "allow_micro_seg", "false"),
					resource.TestCheckResourceAttr(resourceName, "delimiter", ""),
					resource.TestCheckResourceAttr(resourceName, "encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "auto"),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos0"),
					resource.TestCheckResourceAttr(resourceName, "epg_cos_pref", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "lag_policy_name", ""),
					resource.TestCheckResourceAttr(resourceName, "netflow_dir", "both"),
					resource.TestCheckResourceAttr(resourceName, "netflow_pref", "disabled"),
					resource.TestCheckResourceAttr(resourceName, "num_ports", "0"),
					resource.TestCheckResourceAttr(resourceName, "port_allocation", "none"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap_inner", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "res_imedcy", "lazy"),
					resource.TestCheckResourceAttr(resourceName, "secondary_encap_inner", "unknown"),
					resource.TestCheckResourceAttr(resourceName, "switching_mode", "native"),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/vmmp-VMware/dom-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vmm_id", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[uni/vmmp-VMware/dom-%s]/sec", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "vmm_allow_promiscuous", "reject"),
					resource.TestCheckResourceAttr(resourceName, "vmm_forged_transmits", "reject"),
					resource.TestCheckResourceAttr(resourceName, "vmm_mac_changes", "reject"),
				),
			},
			{
				Config: CreateAccFVDomainConfigWithOptionalValues(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "from_terraform"),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "binding_type", "dynamicBinding"),
					resource.TestCheckResourceAttr(resourceName, "allow_micro_seg", "true"),
					resource.TestCheckResourceAttr(resourceName, "delimiter", ""),
					resource.TestCheckResourceAttr(resourceName, "encap", "vlan-5"),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos5"),
					resource.TestCheckResourceAttr(resourceName, "epg_cos_pref", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "instr_imedcy", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "lag_policy_name", "lag_policy_name"),
					resource.TestCheckResourceAttr(resourceName, "netflow_dir", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "netflow_pref", "enabled"),
					resource.TestCheckResourceAttr(resourceName, "num_ports", "3"),
					resource.TestCheckResourceAttr(resourceName, "port_allocation", "fixed"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap", "vlan-4"),
					resource.TestCheckResourceAttr(resourceName, "primary_encap_inner", "vlan-6"),
					resource.TestCheckResourceAttr(resourceName, "res_imedcy", "immediate"),
					resource.TestCheckResourceAttr(resourceName, "secondary_encap_inner", "vlan-7"),
					resource.TestCheckResourceAttr(resourceName, "switching_mode", "AVE"),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/vmmp-VMware/dom-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "vmm_id", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s/rsdomAtt-[uni/vmmp-VMware/dom-%s]/sec", rName, rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "vmm_allow_promiscuous", "accept"),
					resource.TestCheckResourceAttr(resourceName, "vmm_forged_transmits", "accept"),
					resource.TestCheckResourceAttr(resourceName, "vmm_mac_changes", "accept"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vmm_allow_promiscuous", "vmm_forged_transmits", "vmm_mac_changes", "vmm_id"},
			},
			{
				Config: CreateAccFVDomainConfigWithEpgAndDomainName(rName, rother),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/vmmp-VMware/dom-%s", rother)),
					testAccCheckAciFVDomainIdNotEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainConfig(rName),
			},
			{
				Config: CreateAccFVDomainConfigWithEpgAndDomainName(rother, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rother, rother, rother)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/vmmp-VMware/dom-%s", rName)),
					testAccCheckAciFVDomainIdNotEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
		},
	})
}

func TestAccAciFVDomain_Update(t *testing.T) {
	var epg_to_domain_default models.FVDomain
	var epg_to_domain_updated models.FVDomain
	var epg_to_domain_default_fc_domain models.FVDomain
	var epg_to_domain_updated_fc_domain models.FVDomain
	resourceName := "aci_epg_to_domain.test"
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFVDomainConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_default),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "binding_type", "staticBinding"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "binding_type", "staticBinding"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "binding_type", "ephemeral"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "binding_type", "ephemeral"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos1"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos2"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos3"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos4"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos4"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated)),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos6"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos6"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "epg_cos", "Cos7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "epg_cos", "Cos7"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "netflow_dir", "egress"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "netflow_dir", "egress"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "port_allocation", "elastic"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "port_allocation", "elastic"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "res_imedcy", "pre-provision"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated),
					resource.TestCheckResourceAttr(resourceName, "res_imedcy", "pre-provision"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default, &epg_to_domain_updated),
				),
			},
			{
				Config: CreateAccFVDomainConfigFCDomain(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_default_fc_domain),
					resource.TestCheckResourceAttr(resourceName, "application_epg_dn", fmt.Sprintf("uni/tn-%s/ap-%s/epg-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fmt.Sprintf("uni/fc-%s", rName)),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttrFCDomain(rName, "encap_mode", "vlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated_fc_domain),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "vlan"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default_fc_domain, &epg_to_domain_updated_fc_domain),
				),
			},
			{
				Config: CreateAccFVDomainUpdatedAttrFCDomain(rName, "encap_mode", "vxlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists(resourceName, &epg_to_domain_updated_fc_domain),
					resource.TestCheckResourceAttr(resourceName, "encap_mode", "vxlan"),
					testAccCheckAciFVDomainIdEqual(&epg_to_domain_default_fc_domain, &epg_to_domain_updated_fc_domain),
				),
			},
		},
	})
}

func TestAccAciFVDomain_NegativeCases(t *testing.T) {
	rName := acctest.RandString(5)
	longAnnotation := acctest.RandString(129)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFVDomainConfig(rName),
			},
			{
				Config:      CreateAccFVDomainWithInvalidApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainWithInvalidTDn(rName),
				ExpectError: regexp.MustCompile(`Invalid target DN`),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "switching_mode", "AVE"),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "encap_mode", "vlan"),
				ExpectError: regexp.MustCompile(`VLAN encap mode is not allowed for AVE Non-Local switching domain`),
			},
			{
				Config: CreateAccFVDomainUpdatedAttr(rName, "encap", "vlan-5"),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "encap_mode", "vxlan"),
				ExpectError: regexp.MustCompile(`static vlan setting is only allowed when encap mode is vlan.`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "annotation", longAnnotation),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "binding_type", randomValue),
				ExpectError: regexp.MustCompile(`expected binding_type to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "allow_micro_seg", randomValue),
				ExpectError: regexp.MustCompile(`Incorrect attribute value type`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name encap, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "encap_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected encap_mode to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "epg_cos", randomValue),
				ExpectError: regexp.MustCompile(`expected epg_cos to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "epg_cos_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected epg_cos_pref to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "instr_imedcy", randomValue),
				ExpectError: regexp.MustCompile(`expected instr_imedcy to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "netflow_dir", randomValue),
				ExpectError: regexp.MustCompile(`expected netflow_dir to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "netflow_pref", randomValue),
				ExpectError: regexp.MustCompile(`expected netflow_pref to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "lag_policy_name", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`failed validation for value`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "num_ports", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name numPorts, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "port_allocation", randomValue),
				ExpectError: regexp.MustCompile(`expected port_allocation to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "primary_encap", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name primaryEncap, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "primary_encap_inner", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name primaryEncapInner, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "res_imedcy", randomValue),
				ExpectError: regexp.MustCompile(`expected res_imedcy to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "secondary_encap_inner", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name secondaryEncapInner, class fvRsDomAtt (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "switching_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected switching_mode to be one of (.)+, got (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "vmm_allow_promiscuous", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name allowPromiscuous, class vmmSecP (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "vmm_forged_transmits", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name forgedTransmits, class vmmSecP (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, "vmm_mac_changes", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name macChanges, class vmmSecP (.)+`),
			},
			{
				Config:      CreateAccFVDomainUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFVDomainConfig(rName),
			},
		},
	})
}

func TestAccAciFVDomain_MultipleCreateDelete(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFVDomainConfigs(rName),
			},
		},
	})
}

func testAccCheckAciFVDomainDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing epg_to_domain destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_epg_to_domain" {
			cont, err := client.Get(rs.Primary.ID)
			epg_to_domain := models.FVDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("EPG to Domain %s still exists", epg_to_domain.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func CreateAccFVDomainWithInvalidTDn(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain creation with invalid tdn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_application_epg" "test" {
		name = "%s"
		application_profile_dn = aci_application_profile.test.id
	}

	resource "aci_epg_to_domain" "test" {
		application_epg_dn = aci_application_epg.test.id
		tdn = aci_tenant.test.id
	}

	`, rName, rName, rName)
	return resource
}

func CreateAccFVDomainWithInvalidApplicationEPG(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain creation with invalid application_epg_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	
	resource "aci_fc_domain" "test" {
		name = "%s"
	}

	resource "aci_epg_to_domain" "test" {
		application_epg_dn = aci_tenant.test.id
		tdn = aci_fc_domain.test.id
	}

	`, rName, rName)
	return resource
}

func CreateAccFVDomainWithoutTdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg_to_domain without tdn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
	  }
	`, rName, rName, rName)
	return resource

}

func CreateAccFVDomainWithoutApplicationEPG(rName string) string {
	fmt.Println("=== STEP  Basic: testing aci_epg_to_domain without creating application_epg_dn")
	resource := fmt.Sprintf(`
	  resource "aci_vmm_domain" "test" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		tdn                   = aci_vmm_domain.test.id
	  }
	`, rName)
	return resource

}

func CreateAccFVDomainConfigFCDomain(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain creation with required parameters only while tdn has reference of fc_domain")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_fc_domain" "test" {
		name = "%s"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_fc_domain.test.id
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccFVDomainConfigs(rName string) string {
	fmt.Println("=== STEP  testing mutiple epg_to_domain creation")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_vmm_domain" "test1" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	  }

	  resource "aci_vmm_domain" "test2" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	  }

	  resource "aci_vmm_domain" "test3" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	  }
	  
	 resource "aci_epg_to_domain" "test1" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test1.id
	  }

	  resource "aci_epg_to_domain" "test2" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test2.id
	  }

	  resource "aci_epg_to_domain" "test3" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test3.id
	  }
	`, rName, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource
}

func CreateAccFVDomainConfig(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain creation with required parameters only")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_vmm_domain" "test" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
	  }
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccFVDomainConfigWithEpgAndDomainName(r1, r2 string) string {
	fmt.Printf("=== STEP  Basic: testing epg_to_domain creation with epg name %s and vmm_domain name %s\n", r1, r2)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_vmm_domain" "test" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
	  }
	`, r1, r1, r1, r2)
	return resource
}

func CreateAccFVDomainConfigWithOptionalValues(rName string) string {
	fmt.Println("=== STEP  Basic: testing epg to domain creation with optional parameters")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_vmm_domain" "test" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
		annotation 			= "from_terraform"
  		binding_type          = "dynamicBinding"
  		allow_micro_seg       = "true"
  		delimiter             = ""
  		encap                 = "vlan-5"
  		epg_cos               = "Cos5"
  		epg_cos_pref          = "enabled"
  		instr_imedcy          = "immediate"
  		lag_policy_name       = "lag_policy_name"
  		netflow_dir           = "ingress"
  		netflow_pref          = "enabled"
  		num_ports             = "3"
  		port_allocation       = "fixed"
  		primary_encap         = "vlan-4"
  		primary_encap_inner   = "vlan-6"
  		res_imedcy            = "immediate"
  		secondary_encap_inner = "vlan-7"
  		switching_mode        = "AVE"
  		vmm_allow_promiscuous = "accept"
  		vmm_forged_transmits  = "accept"
  		vmm_mac_changes       = "accept"
	  }
	`, rName, rName, rName, rName)
	return resource
}

func testAccCheckAciFVDomainIdEqual(m1, m2 *models.FVDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("epg_to_domain DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFVDomainIdNotEqual(m1, m2 *models.FVDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("epg_to_domain DNs are equal")
		}
		return nil
	}
}

func testAccCheckAciFVDomainExists(name string, domain *models.FVDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		domainFound := models.FVDomainFromContainer(cont)
		if domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Domain %s not found", rs.Primary.ID)
		}
		*domain = *domainFound
		return nil
	}
}

func CreateAccFVDomainUpdatedAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing epg_to_domain updation with %s = %s\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_vmm_domain" "test" {
		name = "%s"
		provider_profile_dn = "uni/vmmp-VMware"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
		%s = "%s"
	  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccFVDomainUpdatedAttrFCDomain(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing epg_to_domain updation with %s = %s while tdn having referece of fc_domain\n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  
	  resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name      = "%s"
	  }

	  resource "aci_application_epg" "test" {
		application_profile_dn = aci_application_profile.test.id
		name                   = "%s"
	  }

	  resource "aci_fc_domain" "test" {
		name = "%s"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_fc_domain.test.id
		%s = "%s"
	  }
	`, rName, rName, rName, rName, attribute, value)
	return resource
}
