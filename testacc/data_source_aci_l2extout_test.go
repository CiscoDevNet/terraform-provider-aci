package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL2OutsideDataSource_Basic(t *testing.T) {
	resourceName := "aci_l2_outside.test"
	dataSourceName := "data.aci_l2_outside.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL2OutsideDSWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL2OutsideDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2OutsideConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
				),
			},
			{
				Config:      CreateAccL2OutsideDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL2OutsideDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccL2OutsideDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL2OutsideConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l2_outside Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l2_outside.test.name
		depends_on = [ aci_l2_outside.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateL2OutsideDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_outside creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_l2_outside" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_l2_outside.test.name
		depends_on = [ aci_l2_outside.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_l2_outside.test.name
		depends_on = [ aci_l2_outside.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccL2OutsideDSWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l2_outside Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_l2_outside.test.name}_invalid"
		depends_on = [ aci_l2_outside.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL2OutsideDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_outside Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l2_outside.test.name
		%s = "%s"
		depends_on = [ aci_l2_outside.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccL2OutsideDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_outside Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_l2_outside" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_l2_outside.test.name
		depends_on = [ aci_l2_outside.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
