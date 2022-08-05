package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVMMDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_vmm_domain.test"
	dataSourceName := "data.aci_vmm_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMDomainDSWithoutRequired(rName, "provider_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMDomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMDomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "provider_profile_dn", resourceName, "provider_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "access_mode", resourceName, "access_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "arp_learning", resourceName, "arp_learning"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ave_time_out", resourceName, "ave_time_out"),
					resource.TestCheckResourceAttrPair(dataSourceName, "config_infra_pg", resourceName, "config_infra_pg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl_knob", resourceName, "ctrl_knob"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delimiter", resourceName, "delimiter"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_ave", resourceName, "enable_ave"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_tag", resourceName, "enable_tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "encap_mode", resourceName, "encap_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enf_pref", resourceName, "enf_pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ep_inventory_type", resourceName, "ep_inventory_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ep_ret_time", resourceName, "ep_ret_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hv_avail_monitor", resourceName, "hv_avail_monitor"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mcast_addr", resourceName, "mcast_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_encap_mode", resourceName, "pref_encap_mode"),
				),
			},
			{
				Config:      CreateAccVMMDomainDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVMMDomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccVMMDomainDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVMMDomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing vmm_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}

	data "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
		name  = aci_vmm_domain.test.name
		depends_on = [ aci_vmm_domain.test ]
	}
	`, vmmProvProfileDn, rName)
	return resource
}

func CreateVMMDomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_domain Data Source without ", attrName)
	rBlock := `
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}
	`
	switch attrName {
	case "provider_profile_dn":
		rBlock += `
	data "aci_vmm_domain" "test" {
	#	provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
		name  = aci_vmm_domain.test.name
		depends_on = [ aci_vmm_domain.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
	#	name  = aci_vmm_domain.test.name
		depends_on = [ aci_vmm_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, vmmProvProfileDn, rName)
}

func CreateAccVMMDomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing vmm_domain Data Source with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}

	data "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
		name  = "${aci_vmm_domain.test.name}_invalid"
		depends_on = [ aci_vmm_domain.test ]
	}
	`, vmmProvProfileDn, rName)
	return resource
}

func CreateAccVMMDomainDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_domain Data Source with random attribute")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
	}

	data "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
		name  = aci_vmm_domain.test.name
		%s = "%s"
		depends_on = [ aci_vmm_domain.test ]
	}
	`, vmmProvProfileDn, rName, key, value)
	return resource
}

func CreateAccVMMDomainDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_domain Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_vmm_domain" "test" {
		provider_profile_dn  = "%s"
		name  = "%s"
		%s = "%s"
	}

	data "aci_vmm_domain" "test" {
		provider_profile_dn  = aci_vmm_domain.test.provider_profile_dn
		name  = aci_vmm_domain.test.name
		depends_on = [ aci_vmm_domain.test ]
	}
	`, vmmProvProfileDn, rName, key, value)
	return resource
}
