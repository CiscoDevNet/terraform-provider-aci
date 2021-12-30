package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outRouteTagPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_route_tag_policy.test"
	dataSourceName := "data.aci_l3out_route_tag_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outRouteTagPolicyDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outRouteTagPolicyDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outRouteTagPolicyConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccL3outRouteTagPolicyDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outRouteTagPolicyDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outRouteTagPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outRouteTagPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_hsrp_interface_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	
	data "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l3out_route_tag_policy.test.name
		depends_on = [
			aci_l3out_route_tag_policy.test
		]
	}

	`, rName, rName, key, value)
	return resource
}

func CreateL3outRouteTagPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_route_tag_policy Data Source without", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_l3out_route_tag_policy" "test" {
	#	tenant_dn  = aci_l3out_route_tag_policy.test.tenant_dn
		name  = aci_l3out_route_tag_policy.test.name
	}
		`
	case "name":
		rBlock += `
	data "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_l3out_route_tag_policy.test.tenant_dn
	#	name  = "aci_l3out_route_tag_policy.test.name"
	}	`
	}

	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccL3outRouteTagPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l3out_route_tag_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l3out_route_tag_policy.test.name
		depends_on = [
			aci_l3out_route_tag_policy.test
		]
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccL3outRouteTagPolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l3out_route_tag_policy creation with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_l3out_route_tag_policy.test.name}_invalid"
		depends_on = [
			aci_l3out_route_tag_policy.test
		]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL3outRouteTagPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l3out_route_tag_policy Data Source with random parameter")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l3out_route_tag_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l3out_route_tag_policy.test.name
		%s = "%s"
		depends_on = [
			aci_l3out_route_tag_policy.test
		]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
