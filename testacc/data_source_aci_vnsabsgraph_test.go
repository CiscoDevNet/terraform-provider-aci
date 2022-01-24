package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL4L7ServiceGraphTemplateDataSource_Basic(t *testing.T) {
	resourceName := "aci_l4_l7_service_graph_template.test"
	dataSourceName := "data.aci_l4_l7_service_graph_template.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL4L7ServiceGraphTemplateDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL4L7ServiceGraphTemplateDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "l4-l7_service_graph_template_type", resourceName, "l4-l7_service_graph_template_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ui_template_type", resourceName, "ui_template_type"),
				),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL4L7ServiceGraphTemplateDSWithInvalidParentDn(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL4L7ServiceGraphTemplateDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL4L7ServiceGraphTemplateConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l4_l7_service_graph_template.test.name
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateL4L7ServiceGraphTemplateDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l4_l7_service_graph_template Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_l4_l7_service_graph_template" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_l4_l7_service_graph_template.test.name
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_l4_l7_service_graph_template.test.name
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccL4L7ServiceGraphTemplateDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_l4_l7_service_graph_template.test.name}_invalid"
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l4_l7_service_graph_template.test.name
		%s = "%s"
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l4_l7_service_graph_template.test.name
		depends_on = [ aci_l4_l7_service_graph_template.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
