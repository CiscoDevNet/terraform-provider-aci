package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAnyDataSource_Basic(t *testing.T) {
	resourceName := "aci_any.test"
	dataSourceName := "data.aci_any.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciAnyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccAnyDSWithoutVRFdn(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAnyDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "match_t", resourceName, "match_t"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref_gr_memb", resourceName, "pref_gr_memb"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vrf_dn", resourceName, "vrf_dn"),
				),
			},
			{
				Config:      CreateAccAnyDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccAnyDSWithInvalidVRFdn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccAnyDataSourceUpdate(rName, "description", "test_annotation_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccAnyDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing any data source reading with giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	data "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		depends_on = [aci_any.test]
	}
	`, rName, rName)
	return resource
}

func CreateAccAnyDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing any data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	data "aci_any" "test" {
		vrf_dn =  aci_vrf.test.id
		%s = "%s"
		depends_on = [aci_any.test]
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccAnyDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing any data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
		%s = "%s"
	}
	data "aci_any" "test" {
		vrf_dn =  aci_vrf.test.id
		depends_on = [aci_any.test]
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccAnyDSWithInvalidVRFdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing any data source reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_vrf" "test" {
		tenant_dn = aci_tenant.test.id
		name  = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	data "aci_any" "test" {
		vrf_dn =  "${aci_vrf.test.id}xyz"
		depends_on = [aci_any.test]
	}
	`, rName, rName)
	return resource
}

func CreateAccAnyDSWithoutVRFdn(rName string) string {
	fmt.Println("=== STEP  Basic: testing any data source reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = %s
	}
	resource "aci_vrf" "test" {
		tenant_dn              = aci_tenant.test.id
		name                   = "%s"
	}
	resource "aci_any" "test" {
		vrf_dn = aci_vrf.test.id
	}
	data "aci_any" "test" {
	} 
	`, rName, rName)
	return resource
}
