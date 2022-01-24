package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMonitoringPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_monitoring_policy.test"
	dataSourceName := "data.aci_monitoring_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMonitoringPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateMonitoringPolicyDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateMonitoringPolicyDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMonitoringPolicyConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccMonitoringPolicyDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMonitoringPolicyDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccMonitoringPolicyDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMonitoringPolicyConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing monitoring_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_monitoring_policy.test.name
		depends_on = [ aci_monitoring_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateMonitoringPolicyDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing monitoring_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_monitoring_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_monitoring_policy.test.name
		depends_on = [ aci_monitoring_policy.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_monitoring_policy.test.name
		depends_on = [ aci_monitoring_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccMonitoringPolicyDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing monitoring_policy Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_monitoring_policy.test.name}_invalid"
		depends_on = [ aci_monitoring_policy.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccMonitoringPolicyDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing monitoring_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_monitoring_policy.test.name
		%s = "%s"
		depends_on = [ aci_monitoring_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccMonitoringPolicyDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing monitoring_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_monitoring_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_monitoring_policy.test.name
		depends_on = [ aci_monitoring_policy.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
