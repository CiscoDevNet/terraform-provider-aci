package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAAADomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_aaa_domain.test"
	dataSourceName := "data.aci_aaa_domain.test"
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAAADomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAAADomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAAADomainDSConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
				),
			},
			{
				Config:      CreateAccAAADomainDSConfigWithRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccAAADomainDSConfigWithUpdatedResource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccAAADomainDSConfigWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAAADomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing aaa_domain data source without ", attrName)
	rBlock := `
	resource "aci_aaa_domain" "test" {
	
		name  = "%s"
	}
	`

	switch attrName {
	case "name":
		rBlock += `
	data "aci_aaa_domain" "test" {
	
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccAAADomainDSConfigWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing aaa_domain data source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
	}

	data "aci_aaa_domain" "test" {
		name  = "${aci_aaa_domain.test.name}xyz"
	}
	`, rName)
	return resource
}

func CreateAccAAADomainDSConfigWithUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing aaa_domain data source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
		%s = "%s"
	}

	data "aci_aaa_domain" "test" {
		name  = aci_aaa_domain.test.name
	}
	`, rName, key, value)
	return resource
}

func CreateAccAAADomainDSConfigWithRandomAttr(rName, key, value string) string {
	fmt.Println("=== STEP  testing aaa_domain data source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
	}

	data "aci_aaa_domain" "test" {
		name  = aci_aaa_domain.test.name
		%s = "%s"
	}
	`, rName, key, value)
	return resource
}

func CreateAccAAADomainDSConfig(rName string) string {
	fmt.Println("=== STEP  testing aaa_domain data source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_aaa_domain" "test" {
		name  = "%s"
	}

	data "aci_aaa_domain" "test" {
		name  = aci_aaa_domain.test.name
	}
	`, rName)
	return resource
}
