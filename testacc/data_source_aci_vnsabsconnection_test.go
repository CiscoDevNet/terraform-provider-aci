package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciConnectionDataSource_Basic(t *testing.T) {
	resourceName := "aci_connection.test"
	dataSourceName := "data.aci_connection.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateConnectionDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "l4_l7_service_graph_template_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateConnectionDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccConnectionConfigDataSource(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l4_l7_service_graph_template_dn", resourceName, "l4_l7_service_graph_template_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "adj_type", resourceName, "adj_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "conn_dir", resourceName, "conn_dir"),
					resource.TestCheckResourceAttrPair(dataSourceName, "conn_type", resourceName, "conn_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "direct_connect", resourceName, "direct_connect"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unicast_route", resourceName, "unicast_route"),
				),
			},
			{
				Config:      CreateAccConnectionDataSourceUpdate(fvTenantName, vnsAbsGraphName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccConnectionDSWithInvalidParentDn(fvTenantName, vnsAbsGraphName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccConnectionDataSourceUpdatedResource(fvTenantName, vnsAbsGraphName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccConnectionConfigDataSource(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing connection Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_connection.test.name
		depends_on = [ aci_connection.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateConnectionDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing connection Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "l4_l7_service_graph_template_dn":
		rBlock += `
	data "aci_connection" "test" {
	#	l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_connection.test.name
		depends_on = [ aci_connection.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
	#	name  = aci_connection.test.name
		depends_on = [ aci_connection.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsAbsGraphName, rName)
}

func CreateAccConnectionDSWithInvalidParentDn(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing connection Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "${aci_connection.test.name}_invalid"
		depends_on = [ aci_connection.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccConnectionDataSourceUpdate(fvTenantName, vnsAbsGraphName, rName, key, value string) string {
	fmt.Println("=== STEP  testing connection Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_connection.test.name
		%s = "%s"
		depends_on = [ aci_connection.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName, key, value)
	return resource
}

func CreateAccConnectionDataSourceUpdatedResource(fvTenantName, vnsAbsGraphName, rName, key, value string) string {
	fmt.Println("=== STEP  testing connection Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_connection" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_connection.test.name
		depends_on = [ aci_connection.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName, key, value)
	return resource
}
