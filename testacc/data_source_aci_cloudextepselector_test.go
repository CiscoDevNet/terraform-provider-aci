package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudEndpointSelectorforExternalEPgsDataSource_Basic(t *testing.T) {
	resourceName := "aci_cloud_endpoint_selectorfor_external_epgs.test"
	dataSourceName := "data.aci_cloud_endpoint_selectorfor_external_epgs.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	subnet, _ := acctest.RandIpAddress("10.4.0.0/19")
	subnet = fmt.Sprintf("%s/19", subnet)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciCloudEndpointSelectorforExternalEPgsDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsDSWithoutRequired(rName, rName, rName, subnet, rName, "cloud_external_epg_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsDSWithoutRequired(rName, rName, rName, subnet, rName, "subnet"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudEndpointSelectorforExternalEPgsDSWithoutRequired(rName, rName, rName, subnet, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsConfigDataSource(rName, rName, rName, subnet, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "cloud_external_epg_dn", resourceName, "cloud_external_epg_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "subnet", resourceName, "subnet"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_shared", resourceName, "is_shared"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_expression", resourceName, "match_expression"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsDataSourceUpdate(rName, rName, rName, subnet, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudEndpointSelectorforExternalEPgsDSWithInvalidParentDn(rName, rName, rName, subnet, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccCloudEndpointSelectorforExternalEPgsDataSourceUpdatedResource(rName, rName, rName, subnet, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccCloudEndpointSelectorforExternalEPgsConfigDataSource(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}

	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)
	return resource
}

func CreateCloudEndpointSelectorforExternalEPgsDSWithoutRequired(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_endpoint_selectorfor_external_epgs Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}

	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}
	`
	switch attrName {
	case "cloud_external_epg_dn":
		rBlock += `
	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
	#	cloud_external_epg_dn  = aci_cloud_endpoint_selectorfor_external_epgs.test.cloud_external_epg_dn
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
		`
	case "subnet":
		rBlock += `
	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_endpoint_selectorfor_external_epgs.test.cloud_external_epg_dn
	#	subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_endpoint_selectorfor_external_epgs.test.cloud_external_epg_dn
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
	#	name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)
}

func CreateAccCloudEndpointSelectorforExternalEPgsDSWithInvalidParentDn(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}

	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = "${aci_cloud_external_epg.test.id}_invalid"
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsDataSourceUpdate(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, key, value string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
	}

	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		%s = "%s"
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, key, value)
	return resource
}

func CreateAccCloudEndpointSelectorforExternalEPgsDataSourceUpdatedResource(fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, key, value string) string {
	fmt.Println("=== STEP  testing cloud_endpoint_selectorfor_external_epgs Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_applicationcontainer" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_cloud_external_epg" "test" {
		name 		= "%s"
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.test.id
	}
	
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = "%s"
		name = "%s"
		%s = "%s"
	}

	data "aci_cloud_endpoint_selectorfor_external_epgs" "test" {
		cloud_external_epg_dn  = aci_cloud_external_epg.test.id
		subnet  = aci_cloud_endpoint_selectorfor_external_epgs.test.subnet
		name = aci_cloud_endpoint_selectorfor_external_epgs.test.name
		depends_on = [ aci_cloud_endpoint_selectorfor_external_epgs.test ]
	}
	`, fvTenantName, cloudAppName, cloudExtEPgName, subnet, name, key, value)
	return resource
}
