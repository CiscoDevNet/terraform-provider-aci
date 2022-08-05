package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTagDataSource_Basic(t *testing.T) {
	resourceName := "aci_tag.test"
	dataSourceName := "data.aci_tag.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)

	key := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateTagDSWithoutRequired(fvTenantName, key, "parent_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateTagDSWithoutRequired(fvTenantName, key, "key"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTagConfigDataSource(fvTenantName, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key", resourceName, "key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "value", resourceName, "value"),
				),
			},
			{
				Config:      CreateAccTagDataSourceUpdate(fvTenantName, key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTagDSWithInvalidParentDn(fvTenantName, key),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccTagDataSourceUpdatedResource(fvTenantName, key, "test_value"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "value", resourceName, "value"),
				),
			},
		},
	})
}

func CreateAccTagConfigDataSource(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing tag Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
	`, fvTenantName, key)
	return resource
}

func CreateTagDSWithoutRequired(fvTenantName, key, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tag Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}
	`
	switch attrName {
	case "parent_dn":
		rBlock += `
	data "aci_tag" "test" {
	#	parent_dn  = aci_tenant.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
		`
	case "key":
		rBlock += `
	data "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
	#	key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, key)
}

func CreateAccTagDSWithInvalidParentDn(fvTenantName, key string) string {
	fmt.Println("=== STEP  testing tag Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "${aci_tag.test.key}_invalid"
		depends_on = [ aci_tag.test ]
	}
	`, fvTenantName, key)
	return resource
}

func CreateAccTagDataSourceUpdate(fvTenantName, key, attr, value string) string {
	fmt.Println("=== STEP  testing tag Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "val"
	}

	data "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_tag.test.key
		%s = "%s"
		depends_on = [ aci_tag.test ]
	}
	`, fvTenantName, key, attr, value)
	return resource
}

func CreateAccTagDataSourceUpdatedResource(fvTenantName, key, value string) string {
	fmt.Println("=== STEP  testing tag Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = "%s"
		value = "%s"
	}

	data "aci_tag" "test" {
		parent_dn  = aci_tenant.test.id
		key  = aci_tag.test.key
		depends_on = [ aci_tag.test ]
	}
	`, fvTenantName, key, value)
	return resource
}
