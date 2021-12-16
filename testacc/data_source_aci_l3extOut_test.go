package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3OutsideDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3_outside.test"
	dataSourceName := "data.aci_l3_outside.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3OutsideDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccL3OutsideDSConfigWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3OutsideDSConfigWithoutTenantdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3OutsideDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_dscp", resourceName, "target_dscp"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enforce_rtctrl", resourceName, "enforce_rtctrl"),
				),
			},
			{
				Config:      CreateAccL3OutsideDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccL3OutsideDSWithInvalidTenantdn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config:      CreateAccL3OutsideDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3OutsideDataSourceUpdate(rName, "description", "test_annotation_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccL3OutsideDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside data source reading with giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	  data "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = aci_tenant.test.name
		depends_on = [aci_l3_outside.test]
	  }
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3outside data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	  data "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = aci_l3_outside.test.name
		%s = "%s"
		depends_on = [aci_l3_outside.test]
	  }
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccL3OutsideDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3outside data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
		%s = "%s"
	  }
	  data "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = aci_l3_outside.test.name
		depends_on = [aci_l3_outside.test]
	  }
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccL3OutsideDSWithInvalidTenantdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside data source reading with invalid tenant dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	data "aci_l3_outside" "test" {
	    tenant_dn      = "${aci_tenant.test.id}xyz"
		name           = aci_l3_outside.test.name
		depends_on = [aci_l3_outside.test]
	}
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside data source reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	}
	data "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "${aci_l3_outside.test.name}xyz"
		depends_on = [aci_l3_outside.test]
	}
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideDSConfigWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside data source reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	  data "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		depends_on = [aci_l3_outside.test]
	  }
	`, rName, rName)
	return resource
}

func CreateAccL3OutsideDSConfigWithoutTenantdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing l3outside data source reading without giving tenant dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name        = "%s"
	  }
	  resource "aci_l3_outside" "test" {
	    tenant_dn      = aci_tenant.test.id
		name           = "%s"
	  }
	  data "aci_l3_outside" "test" {
	    name = aci_l3_outside.test.name
		depends_on = [aci_l3_outside.test]
	  }
	`, rName, rName)
	return resource
}
