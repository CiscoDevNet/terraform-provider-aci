package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudApplicationcontainerDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_applicationcontainer.test"
	dataSourceName := "data.aci_cloud_applicationcontainer.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudApplicationcontainerDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudApplicationcontainerDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudApplicationcontainerDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudApplicationcontainerConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccCloudApplicationcontainerDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccCloudApplicationcontainerDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccCloudApplicationcontainerDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudApplicationcontainerConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_applicationcontainer.test.name
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateCloudApplicationcontainerDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_applicationcontainer Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_cloud_applicationcontainer" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccCloudApplicationcontainerDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_cloud_applicationcontainer.test.name}_invalid"
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccCloudApplicationcontainerDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_applicationcontainer.test.name
		%s = "%s"
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccCloudApplicationcontainerDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_applicationcontainer Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_applicationcontainer" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_cloud_applicationcontainer.test.name
		depends_on = [ aci_cloud_applicationcontainer.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
