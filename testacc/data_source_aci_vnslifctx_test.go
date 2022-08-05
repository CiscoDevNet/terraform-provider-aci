package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLogicalInterfaceContextDataSource_Basic(t *testing.T) {
	resourceName := "aci_logical_interface_context.test"
	dataSourceName := "data.aci_logical_interface_context.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	connNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsLDevCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLogicalInterfaceContextDSWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, "logical_device_context_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLogicalInterfaceContextDSWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, "conn_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfigDataSource(fvTenantName, vnsLDevCtxName, connNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_device_context_dn", resourceName, "logical_device_context_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "conn_name_or_lbl", resourceName, "conn_name_or_lbl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "l3_dest", resourceName, "l3_dest"),
					resource.TestCheckResourceAttrPair(dataSourceName, "permit_log", resourceName, "permit_log"),
				),
			},
			{
				Config:      CreateAccLogicalInterfaceContextDataSourceUpdate(fvTenantName, vnsLDevCtxName, connNameOrLbl, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLogicalInterfaceContextDSWithInvalidParentDn(fvTenantName, vnsLDevCtxName, connNameOrLbl),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccLogicalInterfaceContextDataSourceUpdatedResource(fvTenantName, vnsLDevCtxName, connNameOrLbl, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLogicalInterfaceContextConfigDataSource(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_interface_context Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}

	data "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = aci_logical_interface_context.test.conn_name_or_lbl
		depends_on = [ aci_logical_interface_context.test ]
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)
	return resource
}

func CreateLogicalInterfaceContextDSWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_interface_context Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}
	`
	switch attrName {
	case "logical_device_context_dn":
		rBlock += `
	data "aci_logical_interface_context" "test" {
	#	logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = aci_logical_interface_context.test.conn_name_or_lbl
		depends_on = [ aci_logical_interface_context.test ]
	}
		`
	case "conn_name_or_lbl":
		rBlock += `
	data "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
	#	conn_name_or_lbl  = aci_logical_interface_context.test.conn_name_or_lbl
		depends_on = [ aci_logical_interface_context.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsLDevCtxName, connNameOrLbl)
}

func CreateAccLogicalInterfaceContextDSWithInvalidParentDn(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_interface_context Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}

	data "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "${aci_logical_interface_context.test.conn_name_or_lbl}_invalid"
		depends_on = [ aci_logical_interface_context.test ]
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextDataSourceUpdate(fvTenantName, vnsLDevCtxName, connNameOrLbl, key, value string) string {
	fmt.Println("=== STEP  testing logical_interface_context Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}

	data "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = aci_logical_interface_context.test.conn_name_or_lbl
		%s = "%s"
		depends_on = [ aci_logical_interface_context.test ]
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl, key, value)
	return resource
}

func CreateAccLogicalInterfaceContextDataSourceUpdatedResource(fvTenantName, vnsLDevCtxName, connNameOrLbl, key, value string) string {
	fmt.Println("=== STEP  testing logical_interface_context Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
		%s = "%s"
	}

	data "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = aci_logical_interface_context.test.conn_name_or_lbl
		depends_on = [ aci_logical_interface_context.test ]
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl, key, value)
	return resource
}
