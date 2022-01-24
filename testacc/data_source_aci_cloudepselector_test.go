package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudEndpointSelectorDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_endpoint_selector.test"
	dataSourceName := "data.aci_cloud_endpoint_selector.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEndpointSelectorDSWithoutRequired(rName, rName, rName, rName, "cloud_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorDSWithoutRequired(rName, rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorConfigDataSource(rName, rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_epg_dn", resourceName, "cloud_epg_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_expression", resourceName, "match_expression"),
				),
			},
			{
				Config:      CreateAccCloudEndpointSelectorDataSourceUpdate(rName, rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudEndpointSelectorDSWithInvalidName(rName, rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudEndpointSelectorDataSourceUpdatedResource(rName, rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudEndpointSelectorConfigDataSource(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}

	data "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = aci_cloud_endpoint_selector.test.name
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)
	return resource
}

func CreateCloudEndpointSelectorDSWithoutRequired(fvTenantName, cloudAppName, cloudEPgName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selector Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "cloud_epg_dn":
		rBlock += `
	data "aci_cloud_endpoint_selector" "test" {
	#	cloud_epg_dn  = aci_cloud_epg.test.id
		name  = aci_cloud_endpoint_selector.test.name
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
	#	name  = aci_cloud_endpoint_selector.test.name
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, cloudEPgName, rName)
}

func CreateAccCloudEndpointSelectorDSWithInvalidName(fvTenantName, cloudAppName, cloudEPgName, rName string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}

	data "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "${aci_cloud_endpoint_selector.test.name}_invalid"
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName)
	return resource
}

func CreateAccCloudEndpointSelectorDataSourceUpdate(fvTenantName, cloudAppName, cloudEPgName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
	}

	data "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = aci_cloud_endpoint_selector.test.name
		%s = "%s"
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName, key, value)
	return resource
}

func CreateAccCloudEndpointSelectorDataSourceUpdatedResource(fvTenantName, cloudAppName, cloudEPgName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selector Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_endpoint_selector" "test" {
		cloud_epg_dn  = aci_cloud_epg.test.id
		name  = aci_cloud_endpoint_selector.test.name
		depends_on = [ aci_cloud_endpoint_selector.test ]
	}
	`, fvTenantName, cloudAppName, cloudEPgName, rName, key, value)
	return resource
}
