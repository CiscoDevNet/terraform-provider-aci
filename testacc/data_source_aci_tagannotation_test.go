package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAnnotationDataSource_Basic(t *testing.T) {
	resourceName := "aci_annotation.test"
	dataSourceName := "data.aci_annotation.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	key := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAnnotationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAnnotationDSWithoutRequired(fvTenantName, key, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAnnotationDSWithoutRequired(fvTenantName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnnotationConfigDataSource(fvTenantName, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "value", resourceName, "value"),
				),
			},
			{
				Config:      CreateAccAnnotationDataSourceUpdate(fvTenantName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAnnotationDSWithInvalidParentDn(fvTenantName, key),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccAnnotationDataSourceUpdatedResource(fvTenantName, key, "test_value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "value", resourceName, "value"),
				),
			},
		},
	})
}

func CreateAccAnnotationConfigDataSource(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing annotation Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
	`, fvTenantName, key)
	return resource
}

func CreateAnnotationDSWithoutRequired(fvTenantName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing annotation Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	data "aci_annotation" "test" {
	#	parent_dn  = aci_tenant.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
		`
	case "key":
		rBlock += `
	data "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
	#	key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, key)
}

func CreateAccAnnotationDSWithInvalidParentDn(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing annotation Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "${aci_annotation.test.key}_invalid"
		depends_on = [ aci_annotation.test ]
	}
	`, fvTenantName, key)
	return resource
}

func CreateAccAnnotationDataSourceUpdate(fvTenantName, key, attr, value string) string {
	fmt.Println("=== STEP  testing annotation Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_annotation.test.key
		%s = "%s"
		depends_on = [ aci_annotation.test ]
	}
	`, fvTenantName, key, attr, value)
	return resource
}

func CreateAccAnnotationDataSourceUpdatedResource(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing annotation Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}

	data "aci_annotation" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_annotation.test.key
		depends_on = [ aci_annotation.test ]
	}
	`, fvTenantName, key, value)
	return resource
}
