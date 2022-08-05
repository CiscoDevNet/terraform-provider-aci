package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVMMCredentialDataSource_Basic(t *testing.T) {
	resourceName := "aci_vmm_credential.test"
	dataSourceName := "data.aci_vmm_credential.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	vmmDomPName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVMMCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVMMCredentialDSWithoutRequired(vmmDomPName, rName, "vmm_domain_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVMMCredentialDSWithoutRequired(vmmDomPName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVMMCredentialConfigDataSource(vmmDomPName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "vmm_domain_dn", resourceName, "vmm_domain_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "usr", resourceName, "usr"),
				),
			},
			{
				Config:      CreateAccVMMCredentialDataSourceUpdate(vmmDomPName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVMMCredentialDSWithInvalidParentDn(vmmDomPName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccVMMCredentialDataSourceUpdatedResource(vmmDomPName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVMMCredentialConfigDataSource(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing vmm_credential Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}

	data "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_credential.test.name
		depends_on = [ aci_vmm_credential.test ]
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}

func CreateVMMCredentialDSWithoutRequired(vmmDomPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vmm_credential Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "vmm_domain_dn":
		rBlock += `
	data "aci_vmm_credential" "test" {
	#	vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		depends_on = [ aci_vmm_credential.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
	#	name  = "%s"
		depends_on = [ aci_vmm_credential.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, vmmDomPName, providerProfileDn, rName)
}

func CreateAccVMMCredentialDSWithInvalidParentDn(vmmDomPName, rName string) string {
	fmt.Println("=== STEP  testing vmm_credential Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}

	data "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "${aci_vmm_credential.test.name}_invalid"
		depends_on = [ aci_vmm_credential.test ]
	}
	`, vmmDomPName, providerProfileDn, rName)
	return resource
}

func CreateAccVMMCredentialDataSourceUpdate(vmmDomPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_credential Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
	}

	data "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_credential.test.name
		%s = "%s"
		depends_on = [ aci_vmm_credential.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, key, value)
	return resource
}

func CreateAccVMMCredentialDataSourceUpdatedResource(vmmDomPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing vmm_credential Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vmm_domain" "test" {
		name 		= "%s"
		provider_profile_dn = "%v"
	}
	
	resource "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_vmm_credential" "test" {
		vmm_domain_dn  = aci_vmm_domain.test.id
		name  = aci_vmm_credential.test.name
		depends_on = [ aci_vmm_credential.test ]
	}
	`, vmmDomPName, providerProfileDn, rName, key, value)
	return resource
}
