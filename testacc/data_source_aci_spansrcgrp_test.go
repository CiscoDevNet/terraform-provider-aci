package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSPANSourceGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_span_source_group.test"
	dataSourceName := "data.aci_span_source_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSPANSourceGroupDSWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSPANSourceGroupDSWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSPANSourceGroupConfigDataSource(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
				),
			},
			{
				Config:      CreateAccSPANSourceGroupDataSourceUpdate(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSPANSourceGroupDSWithInvalidName(fvTenantName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccSPANSourceGroupDataSourceUpdatedResource(fvTenantName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSPANSourceGroupConfigDataSource(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_source_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_span_source_group.test.name
		depends_on = [ aci_span_source_group.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateSPANSourceGroupDSWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing span_source_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	data "aci_span_source_group" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = aci_span_source_group.test.name
		depends_on = [ aci_span_source_group.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = aci_span_source_group.test.name
		depends_on = [ aci_span_source_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccSPANSourceGroupDSWithInvalidName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing span_source_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "${aci_span_source_group.test.name}_invalid"
		depends_on = [ aci_span_source_group.test ]
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccSPANSourceGroupDataSourceUpdate(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing span_source_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	data "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_span_source_group.test.name
		%s = "%s"
		depends_on = [ aci_span_source_group.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}

func CreateAccSPANSourceGroupDataSourceUpdatedResource(fvTenantName, rName, key, value string) string {
	fmt.Println("=== STEP  testing span_source_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_span_source_group" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = aci_span_source_group.test.name
		depends_on = [ aci_span_source_group.test ]
	}
	`, fvTenantName, rName, key, value)
	return resource
}
