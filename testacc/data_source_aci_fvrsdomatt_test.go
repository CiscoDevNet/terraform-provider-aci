package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFVDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_epg_to_domain.test"
	dataSourceName := "data.aci_epg_to_domain.test"
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFVDomainDSWithoutTdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccFVDomainDSWithoutApplicationEPG(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFVDomainDSConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "application_epg_dn", dataSourceName, "application_epg_dn"),
					resource.TestCheckResourceAttrPair(resourceName, "tdn", dataSourceName, "tdn"),
					resource.TestCheckResourceAttrPair(resourceName, "annotation", dataSourceName, "annotation"),
					resource.TestCheckResourceAttrPair(resourceName, "binding_type", dataSourceName, "binding_type"),
					resource.TestCheckResourceAttrPair(resourceName, "allow_micro_seg", dataSourceName, "allow_micro_seg"),
					resource.TestCheckResourceAttrPair(resourceName, "delimiter", dataSourceName, "delimiter"),
					resource.TestCheckResourceAttrPair(resourceName, "encap", dataSourceName, "encap"),
					resource.TestCheckResourceAttrPair(resourceName, "encap_mode", dataSourceName, "encap_mode"),
					resource.TestCheckResourceAttrPair(resourceName, "epg_cos", dataSourceName, "epg_cos"),
					resource.TestCheckResourceAttrPair(resourceName, "epg_cos_pref", dataSourceName, "epg_cos_pref"),
					resource.TestCheckResourceAttrPair(resourceName, "instr_imedcy", dataSourceName, "instr_imedcy"),
					resource.TestCheckResourceAttrPair(resourceName, "lag_policy_name", dataSourceName, "lag_policy_name"),
					resource.TestCheckResourceAttrPair(resourceName, "netflow_dir", dataSourceName, "netflow_dir"),
					resource.TestCheckResourceAttrPair(resourceName, "netflow_pref", dataSourceName, "netflow_pref"),
					resource.TestCheckResourceAttrPair(resourceName, "num_ports", dataSourceName, "num_ports"),
					resource.TestCheckResourceAttrPair(resourceName, "port_allocation", dataSourceName, "port_allocation"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_encap", dataSourceName, "primary_encap"),
					resource.TestCheckResourceAttrPair(resourceName, "primary_encap_inner", dataSourceName, "primary_encap_inner"),
					resource.TestCheckResourceAttrPair(resourceName, "res_imedcy", dataSourceName, "res_imedcy"),
					resource.TestCheckResourceAttrPair(resourceName, "secondary_encap_inner", dataSourceName, "secondary_encap_inner"),
					resource.TestCheckResourceAttrPair(resourceName, "switching_mode", dataSourceName, "switching_mode"),
				),
			},
			{
				Config:      CreateAccFVDomainDSConfigUpdatedAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFVDomainDSConfigWithInvalidEPgDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFVDomainDSConfigWithUpdatedResource(rName, "annotation", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "annotation", dataSourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFVDomainDSWithoutTdn(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain Data Source without tdn")
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

	  data "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_epg_to_domain.test.application_epg_dn
	  }
	`, rName, rName, rName, rName)
	return resource

}

func CreateAccFVDomainDSWithoutApplicationEPG(rName string) string {
	fmt.Println("=== STEP  testing aci_epg_to_domain Data Source without application_epg_dn")
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

	  data "aci_epg_to_domain" "test" {
		tdn    = aci_epg_to_domain.test.tdn
	  }
	`, rName, rName, rName, rName)
	return resource

}

func CreateAccFVDomainDSConfigWithInvalidEPgDn(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain Data Source with invalid application_epg_dn")
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
		provider_profile_dn = "%s"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
	  }

	  data "aci_epg_to_domain" "test" {
		application_epg_dn    = "${aci_epg_to_domain.test.application_epg_dn}xyz"
		tdn                   = aci_epg_to_domain.test.tdn
	  }
	`, rName, rName, rName, rName, vmmProvProfileDn)
	return resource
}

func CreateAccFVDomainDSConfigWithUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing epg_to_domain Data Source with updated resource")
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
		provider_profile_dn = "%s"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
		%s = "%s"
	  }

	  data "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_epg_to_domain.test.application_epg_dn
		tdn                   = aci_epg_to_domain.test.tdn
	  }
	`, rName, rName, rName, rName, vmmProvProfileDn, key, value)
	return resource
}

func CreateAccFVDomainDSConfigUpdatedAttr(rName, key, value string) string {
	fmt.Println("=== STEP  testing epg_to_domain Data Source with random attribute")
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
		provider_profile_dn = "%s"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
	  }

	  data "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_epg_to_domain.test.application_epg_dn
		tdn                   = aci_epg_to_domain.test.tdn
		%s                    = "%s"
	  }
	`, rName, rName, rName, rName, vmmProvProfileDn, key, value)
	return resource
}

func CreateAccFVDomainDSConfig(rName string) string {
	fmt.Println("=== STEP  testing epg_to_domain Data Source required parameters only")
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
		provider_profile_dn = "%s"
		enable_ave = "yes"
	  }
	  
	 resource "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_application_epg.test.id
		tdn                   = aci_vmm_domain.test.id
	  }

	  data "aci_epg_to_domain" "test" {
		application_epg_dn    = aci_epg_to_domain.test.application_epg_dn
		tdn                   = aci_epg_to_domain.test.tdn
	  }
	`, rName, rName, rName, rName, vmmProvProfileDn)
	return resource
}
