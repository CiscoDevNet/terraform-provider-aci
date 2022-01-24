package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudEPgDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_epg.test"
	dataSourceName := "data.aci_cloud_epg.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEPgDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEPgDSWithoutRequired(rName, rName, rName, "cloud_applicationcontainer_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEPgDSWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEPgConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_applicationcontainer_dn", resourceName, "cloud_applicationcontainer_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "flood_on_encap", resourceName, "flood_on_encap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "exception_tag", resourceName, "exception_tag"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),
				),
			},
			{
				Config:      CreateAccCloudEPgDataSourceUpdate(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudEPgDSWithInvalidName(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudEPgDataSourceUpdatedResource(rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudEPgConfigDataSource(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_epg Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_epg.test.name
		depends_on = [ aci_cloud_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateCloudEPgDSWithoutRequired(fvTenantName, cloudAppName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_epg Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "cloud_applicationcontainer_dn":
		rBlock += `
	data "aci_cloud_epg" "test" {
	#	cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_epg.test.name
		depends_on = [ aci_cloud_epg.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
	#	name  = aci_cloud_epg.test.name
		depends_on = [ aci_cloud_epg.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, rName)
}

func CreateAccCloudEPgDSWithInvalidName(fvTenantName, cloudAppName, rName string) string {
	fmt.Println("=== STEP  testing cloud_epg Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "${aci_cloud_epg.test.name}_invalid"
		depends_on = [ aci_cloud_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName)
	return resource
}

func CreateAccCloudEPgDataSourceUpdate(fvTenantName, cloudAppName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_epg Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
	}

	data "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_epg.test.name
		%s = "%s"
		depends_on = [ aci_cloud_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName, key, value)
	return resource
}

func CreateAccCloudEPgDataSourceUpdatedResource(fvTenantName, cloudAppName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_epg Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_epg" "test" {
		cloud_applicationcontainer_dn  = aci_cloud_applicationcontainer.test.id
		name  = aci_cloud_epg.test.name
		depends_on = [ aci_cloud_epg.test ]
	}
	`, fvTenantName, cloudAppName, rName, key, value)
	return resource
}
