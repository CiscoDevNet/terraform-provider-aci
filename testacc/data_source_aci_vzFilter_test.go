package acctest

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFilterDataSource_Basic(t *testing.T) {
	resourceName := "aci_filter.test"                                                 // defining name of resource
	dataSourceName := "data.aci_filter.test"                                          // defining name of data source
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccFilterDSWithoutTenant(rName),           // creating data source for application profile without required arguement tenant_dn
				ExpectError: regexp.MustCompile(`Missing required argument`), // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccFilterDSWithoutName(rName), // creating data source for application profile without required arguement name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFilterConfigDataSource(rName), // creating data source with required arguements from the resource
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"), // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),   // comparing value of parameter description in data source and resoruce

					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccFilterDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFilterDSWithInvalidName(rName),          // data source configuration with invalid application profile profile name
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`), // test step expect error which should be match with defined regex
			},
			{
				Config: CreateAccFilterDataSourceUpdate(rName, "description", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccFilterDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing filter data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	data "aci_filter" "test" {
		name = aci_filter.test.name
		tenant_dn = aci_filter.test.tenant_dn
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing filter data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_filter" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "aci_filter" "test" {
		name = aci_filter.test.name
		tenant_dn = aci_filter.test.tenant_dn
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccFilterConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_filter.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccFilterDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_filter.test.name}xyz"
	}
	`, rName, rName)
	return resource
}

func CreateAccFilterDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_filter" "test" {
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccFilterDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing filter reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_filter" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName)
	return resource
}
