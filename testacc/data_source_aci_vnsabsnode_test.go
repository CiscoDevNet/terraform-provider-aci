package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFunctionNodeDataSource_Basic(t *testing.T) {
	resourceName := "aci_function_node.test"
	dataSourceName := "data.aci_function_node.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFunctionNodeDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "l4_l7_service_graph_template_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFunctionNodeDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFunctionNodeConfigDataSource(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "l4_l7_service_graph_template_dn", resourceName, "l4_l7_service_graph_template_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "func_template_type", resourceName, "func_template_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "func_type", resourceName, "func_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "is_copy", resourceName, "is_copy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "managed", resourceName, "managed"),
					resource.TestCheckResourceAttrPair(dataSourceName, "routing_mode", resourceName, "routing_mode"),
					resource.TestCheckResourceAttrPair(dataSourceName, "sequence_number", resourceName, "sequence_number"),
					resource.TestCheckResourceAttrPair(dataSourceName, "share_encap", resourceName, "share_encap"),
				),
			},
			{
				Config:      CreateAccFunctionNodeDataSourceUpdate(fvTenantName, vnsAbsGraphName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFunctionNodeDSWithInvalidParentDn(fvTenantName, vnsAbsGraphName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccFunctionNodeDataSourceUpdatedResource(fvTenantName, vnsAbsGraphName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFunctionNodeConfigDataSource(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing function_node Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_function_node.test.name
		depends_on = [ aci_function_node.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateFunctionNodeDSWithoutRequired(fvTenantName, vnsAbsGraphName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing function_node Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "l4_l7_service_graph_template_dn":
		rBlock += `
	data "aci_function_node" "test" {
	#	l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_function_node.test.name
		depends_on = [ aci_function_node.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
	#	name  = aci_function_node.test.name
		depends_on = [ aci_function_node.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsAbsGraphName, rName)
}

func CreateAccFunctionNodeDSWithInvalidParentDn(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing function_node Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "${aci_function_node.test.name}_invalid"
		depends_on = [ aci_function_node.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccFunctionNodeDataSourceUpdate(fvTenantName, vnsAbsGraphName, rName, key, value string) string {
	fmt.Println("=== STEP  testing function_node Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}

	data "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_function_node.test.name
		%s = "%s"
		depends_on = [ aci_function_node.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName, key, value)
	return resource
}

func CreateAccFunctionNodeDataSourceUpdatedResource(fvTenantName, vnsAbsGraphName, rName, key, value string) string {
	fmt.Println("=== STEP  testing function_node Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = aci_function_node.test.name
		depends_on = [ aci_function_node.test ]
	}
	`, fvTenantName, vnsAbsGraphName, rName, key, value)
	return resource
}
