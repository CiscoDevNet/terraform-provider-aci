package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciPhysicalDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_physical_domain.test"
	dataSourceName := "data.aci_physical_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreatePhysicalDomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccPhysicalDomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccPhysicalDomainDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccPhysicalDomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccPhysicalDomainDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccPhysicalDomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing physical_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_physical_domain" "test" {

		name  = "%s"
	}

	data "aci_physical_domain" "test" {

		name  = aci_physical_domain.test.name
		depends_on = [ aci_physical_domain.test ]
	}
	`, rName)
	return resource
}

func CreateAccPhysicalDomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing physical_domain Data Source with Invalid Name")
	resource := fmt.Sprintf(`

	resource "aci_physical_domain" "test" {
		name  = "%s"
	}

	data "aci_physical_domain" "test" {

		name  = "${aci_physical_domain.test.name}_invalid"
		depends_on = [ aci_physical_domain.test ]
	}
	`, rName)
	return resource
}

func CreatePhysicalDomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing physical_domain Data Source without ", attrName)
	rBlock := `

	resource "aci_physical_domain" "test" {
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_physical_domain" "test" {

	#	name  = "%s"
		depends_on = [ aci_physical_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccPhysicalDomainDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing physical_domain Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_physical_domain" "test" {
		name  = "%s"
	}

	data "aci_physical_domain" "test" {
		name  = aci_physical_domain.test.name
		%s = "%s"
		depends_on = [ aci_physical_domain.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccPhysicalDomainDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing physical_domain Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_physical_domain" "test" {
		name  = "%s"
		%s = "%s"
	}

	data "aci_physical_domain" "test" {
		name  = aci_physical_domain.test.name
		depends_on = [ aci_physical_domain.test ]
	}
	`, rName, key, value)
	return resource
}
