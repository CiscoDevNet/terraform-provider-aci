package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciOspfRouteSummarizationDataSource_Basic(t *testing.T) {
	resourceName := "aci_ospf_route_summarization.test"
	dataSourceName := "data.aci_ospf_route_summarization.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOspfRouteSummarizationDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOspfRouteSummarizationDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cost", resourceName, "cost"),
					resource.TestCheckResourceAttrPair(dataSourceName, "inter_area_enabled", resourceName, "inter_area_enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tag", resourceName, "tag"),
				),
			},
			{
				Config:      CreateAccOspfRouteSummarizationDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccOspfRouteSummarizationDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccOspfRouteSummarizationDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccOspfRouteSummarizationConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_route_summarization Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_route_summarization.test.name
		depends_on = [ aci_ospf_route_summarization.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateOspfRouteSummarizationDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_route_summarization creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_ospf_route_summarization" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_ospf_route_summarization.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_ospf_route_summarization.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOspfRouteSummarizationDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing ospf_route_summarization Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_ospf_route_summarization.test.name}_invalid"
		depends_on = [ aci_ospf_route_summarization.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOspfRouteSummarizationDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_route_summarization Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_route_summarization.test.name
		%s = "%s"
		depends_on = [ aci_ospf_route_summarization.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccOspfRouteSummarizationDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing ospf_route_summarization Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_ospf_route_summarization.test.name
		depends_on = [ aci_ospf_route_summarization.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
