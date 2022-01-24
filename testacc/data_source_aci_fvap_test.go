package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciApplicationProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_application_profile.test"                                    // defining name of resource
	dataSourceName := "data.aci_application_profile.test"                             // defining name of data source
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz") // creating random string of 5 characters (to give as random parameter)
	randomValue := acctest.RandString(5)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciApplicationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccApplicationProfileDSWithoutTenant(rName), // creating data source for application profile without required argument tenant_dn
				ExpectError: regexp.MustCompile(`Missing required argument`),   // test step expect error which should be match with defined regex
			},
			{
				Config:      CreateAccApplicationProfileDSWithoutName(rName), // creating data source for application profile without required argument name
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccApplicationProfileConfigDataSource(rName), // creating data source with required arguments from the resource
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"), // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),   // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "prio", resourceName, "prio"),               // comparing value of parameter description in data source and resoruce
					resource.TestCheckResourceAttrPair(dataSourceName, "tenant_dn", resourceName, "tenant_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccApplicationProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccApplicationProfileDSWithInvalidName(rName), // data source configuration with invalid application profile profile name
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),    // test step expect error which should be match with defined regex
			},
			{
				Config: CreateAccApplicationProfileDataSourceUpdate(rName, "description", "description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
		},
	})
}

func CreateAccApplicationProfileDataSourceUpdateRandomAttr(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing application profile data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	data "aci_application_profile" "test" {
		name = aci_application_profile.test.name
		tenant_dn = aci_application_profile.test.tenant_dn
		%s = "%s"
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationProfileDataSourceUpdate(rName, attribute, value string) string {
	fmt.Printf("=== STEP  Basic: testing application profile data source update for attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn = aci_tenant.test.id
		name = "%s"
		%s = "%s"
	}
	data "aci_application_profile" "test" {
		name = aci_application_profile.test.name
		tenant_dn = aci_application_profile.test.tenant_dn
	}
	`, rName, rName, attribute, value)
	return resource
}

func CreateAccApplicationProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: testing application profile creation for data source test")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = aci_application_profile.test.name
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  Basic: testing application profile reading with invalid name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "${aci_application_profile.test.name}xyz"
	}
	`, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithoutTenant(rName string) string {
	fmt.Println("=== STEP  Basic: testing application profile reading without giving tenant_dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		name = "%s"
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccApplicationProfileDSWithoutName(rName string) string {
	fmt.Println("=== STEP  Basic: testing application profile reading without giving name")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name = "%s"
	}

	resource "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	data "aci_application_profile" "test" {
		tenant_dn = aci_tenant.test.id
	}
	`, rName, rName)
	return resource
}
