package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVMMControllerDataSource_Basic(t *testing.T) {
	resourceName := "aci_vmm_controller.test"
	dataSourceName := "data.aci_vmm_controller.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("10.2.0.0/16")
	rootContName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMControllerDSWithoutRequired(vmmDomPName, rName, ip, rootContName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMControllerDSWithoutRequired(vmmDomPName, rName, ip, rootContName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMControllerConfigDataSource(vmmDomPName, rName, ip, rootContName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "vmm_domain_dn", resourceName, "vmm_domain_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dvs_version", resourceName, "dvs_version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "host_or_ip", resourceName, "host_or_ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "inventory_trig_st", resourceName, "inventory_trig_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "msft_config_err_msg", resourceName, "msft_config_err_msg"),
					resource.TestCheckResourceAttrPair(dataSourceName, "msft_config_issues.#", resourceName, "msft_config_issues.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "msft_config_issues.0", resourceName, "msft_config_issues.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "n1kv_stats_mode", resourceName, "n1kv_stats_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "port", resourceName, "port"),
					resource.TestCheckResourceAttrPair(dataSourceName, "root_cont_name", resourceName, "root_cont_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "scope", resourceName, "scope"),
					resource.TestCheckResourceAttrPair(dataSourceName, "seq_num", resourceName, "seq_num"),
					resource.TestCheckResourceAttrPair(dataSourceName, "stats_mode", resourceName, "stats_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vxlan_depl_pref", resourceName, "vxlan_depl_pref"),
				),
			},
			{
				Config:      CreateAccVMMControllerDataSourceUpdate(vmmDomPName, rName, ip, rootContName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccVMMControllerDSWithInvalidParentDn(vmmDomPName, rName, ip, rootContName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVMMControllerDataSourceUpdatedResource(vmmDomPName, rName, ip, rootContName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVMMControllerConfigDataSource(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}

	data "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_controller.test.name
		host_or_ip  = aci_vmm_controller.test.host_or_ip
		root_cont_name  = aci_vmm_controller.test.root_cont_name
		depends_on = [ aci_vmm_controller.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)
	return resource
}

func CreateVMMControllerDSWithoutRequired(vmmDomPName, rName, ip, rootContName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_controller Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
		data "aci_vmm_controller" "test" {
		#	vmm_domain_dn  = aci_vmm_domain.test.id
			name  = aci_vmm_controller.test.name
			host_or_ip  = aci_vmm_controller.test.host_or_ip
			root_cont_name  = aci_vmm_controller.test.root_cont_name
			depends_on = [ aci_vmm_controller.test ]
		}
		`
	case "name":
		rBlock += `
		data "aci_vmm_controller" "test" {
			vmm_domain_dn  = aci_vmm_domain.test.id
		#	name  = aci_vmm_controller.test.name
			host_or_ip  = aci_vmm_controller.test.host_or_ip
			root_cont_name  = aci_vmm_controller.test.root_cont_name
			depends_on = [ aci_vmm_controller.test ]
		}
		`
	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn, rName, ip, rootContName)
}

func CreateAccVMMControllerDSWithInvalidParentDn(vmmDomPName, rName, ip, rootContName string) string {
	fmt.Println("=== STEP  testing vmm_controller Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}

	data "aci_vmm_controller" "test" {
		vmm_domain_dn  = "${aci_vmm_domain.test.id}_invalid"
		name  = aci_vmm_controller.test.name
		host_or_ip  = aci_vmm_controller.test.host_or_ip
		root_cont_name  = aci_vmm_controller.test.root_cont_name
		depends_on = [ aci_vmm_controller.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName)
	return resource
}

func CreateAccVMMControllerDataSourceUpdate(vmmDomPName, rName, ip, rootContName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_controller Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
	}

	data "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_controller.test.name
		%s = "%s"
		depends_on = [ aci_vmm_controller.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName, key, value)
	return resource
}

func CreateAccVMMControllerDataSourceUpdatedResource(vmmDomPName, rName, ip, rootContName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_controller Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		host_or_ip = "%s"
		root_cont_name = "%s"
		%s = "%s"
	}

	data "aci_vmm_controller" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_controller.test.name
		host_or_ip  = aci_vmm_controller.test.host_or_ip
		root_cont_name  = aci_vmm_controller.test.root_cont_name
		depends_on = [ aci_vmm_controller.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, ip, rootContName, key, value)
	return resource
}
