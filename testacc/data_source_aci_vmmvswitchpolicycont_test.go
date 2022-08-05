package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVSwitchPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_vswitch_policy.test"
	dataSourceName := "data.aci_vswitch_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVSwitchPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVSwitchPolicyDSWithoutRequired(rName, rName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			}, {
				Config: CreateAccVSwitchPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "vmm_domain_dn", resourceName, "vmm_domain_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVSwitchPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccVSwitchPolicyDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVSwitchPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVSwitchPolicyConfigDataSource(vmmProvPName, vmmDomPName string) string {
	fmt.Println("=== STEP  testing vswitch_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}

	data "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		depends_on = [ aci_vswitch_policy.test ]
	}
	`, vmmDomPName, providerProfileDn)
	return resource
}

func CreateVSwitchPolicyDSWithoutRequired(vmmProvPName, vmmDomPName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vswitch_policy Data Source without ", attrName)
	rBlock := `
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%s"
	}
	
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
	data "aci_vswitch_policy" "test" {
	#	vmm_domain_dn  = aci_vmm_domain.test.id
		depends_on = [ aci_vswitch_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn)
}

func CreateAccVSwitchPolicyDSWithInvalidParentDn(vmmProvPName, vmmDomPName string) string {
	fmt.Println("=== STEP  testing vswitch_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}

	data "aci_vswitch_policy" "test" {
		vmm_domain_dn  = "${aci_vmm_domain.test.id}_invalid"
		depends_on = [ aci_vswitch_policy.test ]
	}
	`, vmmDomPName, providerProfileDn)
	return resource
}

func CreateAccVSwitchPolicyDataSourceUpdate(vmmProvPName, vmmDomPName, key, value string) string {
	fmt.Println("=== STEP  testing vswitch_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	}

	data "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		%s = "%s"
		depends_on = [ aci_vswitch_policy.test ]
	}
	`, vmmDomPName, providerProfileDn, key, value)
	return resource
}

func CreateAccVSwitchPolicyDataSourceUpdatedResource(vmmProvPName, vmmDomPName, key, value string) string {
	fmt.Println("=== STEP  testing vswitch_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}

	resource "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		%s = "%s"
	}

	data "aci_vswitch_policy" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		depends_on = [ aci_vswitch_policy.test ]
	}
	`, vmmDomPName, providerProfileDn, key, value)
	return resource
}
