package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciBgpRouteSummarizationDataSource_Basic(t *testing.T) {
	resourceName := "aci_bgp_route_summarization.test"
	dataSourceName := "data.aci_bgp_route_summarization.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpRouteSummarizationDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpRouteSummarizationDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "attrmap", resourceName, "attrmap"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
				),
			},
			{
				Config:      CreateAccBgpRouteSummarizationDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccBgpRouteSummarizationDSWithInvalidName(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccBgpRouteSummarizationDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccBgpRouteSummarizationConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_route_summarization.test.name
		depends_on = [ aci_bgp_route_summarization.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateBgpRouteSummarizationDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_route_summarization Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_bgp_route_summarization" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		depends_on = [ aci_bgp_route_summarization.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		depends_on = [ aci_bgp_route_summarization.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpRouteSummarizationDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_bgp_route_summarization.test.name}_invalid"
		depends_on = [ aci_bgp_route_summarization.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteSummarizationDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_route_summarization.test.name
		%s = "%s"
		depends_on = [ aci_bgp_route_summarization.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccBgpRouteSummarizationDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_bgp_route_summarization.test.name
		depends_on = [ aci_bgp_route_summarization.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
