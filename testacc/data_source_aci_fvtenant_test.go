package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTenantDataSource_Basic(t *testing.T) {
	resourceName := "aci_tenant.test"
	dataSourceName := "data.aci_tenant.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccTenantDSWithoutName(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTenantDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config:      CreateAccTenantDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccTenantDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTenantDataSourceUpdate(rName, "description", "test_annotation_1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccTenantDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading with giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
		annotation = "test_annotation"
		name_alias = "testing_name_alias"
		description = "testing_description"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}"
	}
	`, rName)
	return resource
}

func CreateAccTenantDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tenant data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}"
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccTenantDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing tenant data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
		%s = "%s"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccTenantDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
		name = "${aci_tenant.test.name}xyz"
	}
	`, rName)
	return resource
}

func CreateAccTenantDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing tenant data source reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	data "aci_tenant" "test" {
	}
	`, rName)
	return resource
}
