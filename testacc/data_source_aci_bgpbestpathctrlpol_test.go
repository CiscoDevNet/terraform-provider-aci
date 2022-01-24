package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBgpBestPathPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_bgp_best_path_policy.test"
	dataSourceName := "data.aci_bgp_best_path_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpBestPathPolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpBestPathPolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
				),
			},
			{
				Config:      CreateAccBgpBestPathPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccBgpBestPathPolicyDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccBgpBestPathPolicyDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBgpBestPathPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_best_path_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_best_path_policy.test.name
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateBgpBestPathPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_best_path_policy Data Source without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bgp_best_path_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpBestPathPolicyDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_best_path_policy Data Source with invalid name")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_bgp_best_path_policy.test.name}_invalid"
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpBestPathPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_best_path_policy Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_best_path_policy.test.name
		%s = "%s"
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBgpBestPathPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_best_path_policy Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_best_path_policy.test.name
		depends_on = [ aci_bgp_best_path_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
