package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLogicalDeviceContextDataSource_Basic(t *testing.T) {
	resourceName := "aci_logical_device_context.test"
	dataSourceName := "data.aci_logical_device_context.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	ctrctNameOrLbl := makeTestVariable(acctest.RandString(5))
	graphNameOrLbl := makeTestVariable(acctest.RandString(5))
	nodeNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLogicalDeviceContextDSWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLogicalDeviceContextDSWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "ctrct_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateLogicalDeviceContextDSWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "graph_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateLogicalDeviceContextDSWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "node_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigDataSource(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrct_name_or_lbl", resourceName, "ctrct_name_or_lbl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "graph_name_or_lbl", resourceName, "graph_name_or_lbl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "node_name_or_lbl", resourceName, "node_name_or_lbl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "context", resourceName, "context"),
				),
			},
			{
				Config:      CreateAccLogicalDeviceContextDataSourceUpdate(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLogicalDeviceContextDSWithInvalidParentDn(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccLogicalDeviceContextDataSourceUpdatedResource(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLogicalDeviceContextConfigDataSource(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_device_context Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}

	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateLogicalDeviceContextDSWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_device_context Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_logical_device_context" "test" {
	#	tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl	
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl	
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
		`
	case "ctrct_name_or_lbl":
		rBlock += `
	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
	#	ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
		`
	case "graph_name_or_lbl":
		rBlock += `
	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
	#	graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
		`
	case "node_name_or_lbl":
		rBlock += `
	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
	#	node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
}

func CreateAccLogicalDeviceContextDSWithInvalidParentDn(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_device_context Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}

	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "${aci_logical_device_context.test.ctrct_name_or_lbl}_invalid"
		graph_name_or_lbl  = "${aci_logical_device_context.test.graph_name_or_lbl}_invalid"
		node_name_or_lbl  = "${aci_logical_device_context.test.node_name_or_lbl}_invalid"
		depends_on = [ aci_logical_device_context.test ]
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateAccLogicalDeviceContextDataSourceUpdate(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, key, value string) string {
	fmt.Println("=== STEP  testing logical_device_context Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}

	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		%s = "%s"
		depends_on = [ aci_logical_device_context.test ]
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, key, value)
	return resource
}

func CreateAccLogicalDeviceContextDataSourceUpdatedResource(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, key, value string) string {
	fmt.Println("=== STEP  testing logical_device_context Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
		%s = "%s"
	}

	data "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = aci_logical_device_context.test.ctrct_name_or_lbl
		graph_name_or_lbl  = aci_logical_device_context.test.graph_name_or_lbl
		node_name_or_lbl  = aci_logical_device_context.test.node_name_or_lbl
		depends_on = [ aci_logical_device_context.test ]
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, key, value)
	return resource
}
