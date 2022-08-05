package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudExternalEPgDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_external_epg.test"
	dataSourceName := "data.aci_cloud_external_epg.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudExternalEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudExternalEPgDSWithoutRequired(rName, rName, rName, "cloud_applicationcontainer_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudExternalEPgDSWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudExternalEPgConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_applicationcontainer_dn", resourceName, "cloud_applicationcontainer_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
					resource.TestCheckResourceAttrPair(dataSourceName, "route_reachability", resourceName, "route_reachability"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
				),
			},
			{
				Config:      CreateAccCloudExternalEPgDataSourceUpdate(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudExternalEPgDSWithInvalidName(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudExternalEPgDataSourceUpdatedResource(rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudExternalEPgConfigDataSource(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_external_epg Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_external_epg.test.name
		depends_on = [ aci_cloud_external_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateCloudExternalEPgDSWithoutRequired(fvTenantName, cloudAppName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_external_epg Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "cloud_applicationcontainer_dn":
		rBlock += `
	data "aci_cloud_external_epg" "test" {
	#	cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_external_epg.test.name
		depends_on = [ aci_cloud_external_epg.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
	#	name  = aci_cloud_external_epg.test.name
		depends_on = [ aci_cloud_external_epg.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, rName)
}

func CreateAccCloudExternalEPgDSWithInvalidName(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_external_epg Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "${aci_cloud_external_epg.test.name}_invalid"
		depends_on = [ aci_cloud_external_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudExternalEPgDataSourceUpdate(fvTenantName, cloudAppName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_external_epg Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_external_epg.test.name
		%s = "%s"
		depends_on = [ aci_cloud_external_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName, key, value)
	return resource
}

func CreateAccCloudExternalEPgDataSourceUpdatedResource(fvTenantName, cloudAppName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_external_epg Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_external_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_external_epg.test.name
		depends_on = [ aci_cloud_external_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName, key, value)
	return resource
}
