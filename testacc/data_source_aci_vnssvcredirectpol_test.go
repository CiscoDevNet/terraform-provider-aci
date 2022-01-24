package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciServiceRedirectPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_service_redirect_policy.test"
	dataSourceName := "data.aci_service_redirect_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciServiceRedirectPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateServiceRedirectPolicyDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateServiceRedirectPolicyDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccServiceRedirectPolicyConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "anycast_enabled", resourceName, "anycast_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "dest_type", resourceName, "dest_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hashing_algorithm", resourceName, "hashing_algorithm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_threshold_percent", resourceName, "max_threshold_percent"),
					resource.TestCheckResourceAttrPair(dataSourceName, "min_threshold_percent", resourceName, "min_threshold_percent"),
					resource.TestCheckResourceAttrPair(dataSourceName, "program_local_pod_only", resourceName, "program_local_pod_only"),
					resource.TestCheckResourceAttrPair(dataSourceName, "resilient_hash_enabled", resourceName, "resilient_hash_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "threshold_down_action", resourceName, "threshold_down_action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "threshold_enable", resourceName, "threshold_enable"),
				),
			},
			{
				Config:      CreateAccServiceRedirectPolicyDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccServiceRedirectPolicyDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccServiceRedirectPolicyDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccServiceRedirectPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing service_redirect_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_service_redirect_policy.test.name
		depends_on = [ aci_service_redirect_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateServiceRedirectPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing service_redirect_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_service_redirect_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_service_redirect_policy.test.name
		depends_on = [ aci_service_redirect_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_service_redirect_policy.test.name
		depends_on = [ aci_service_redirect_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccServiceRedirectPolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing service_redirect_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_service_redirect_policy.test.name}_invalid"
		depends_on = [ aci_service_redirect_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccServiceRedirectPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing service_redirect_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_service_redirect_policy.test.name
		%s = "%s"
		depends_on = [ aci_service_redirect_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccServiceRedirectPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing service_redirect_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_service_redirect_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_service_redirect_policy.test.name
		depends_on = [ aci_service_redirect_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
