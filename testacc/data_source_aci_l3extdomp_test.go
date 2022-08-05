package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3DomainProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3_domain_profile.test"
	dataSourceName := "data.aci_l3_domain_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3DomainProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3DomainProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccL3DomainProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3DomainProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3DomainProfileDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3DomainProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing l3_domain_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}

	data "aci_l3_domain_profile" "test" {

		name  = aci_l3_domain_profile.test.name
		depends_on = [ aci_l3_domain_profile.test ]
	}
	`, rName)
	return resource
}

func CreateL3DomainProfileDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_domain_profile Data Source without ", attrName)
	rBlock := `

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_l3_domain_profile" "test" {

	#	name  = "%s"
		depends_on = [ aci_l3_domain_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL3DomainProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3_domain_profile Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}

	data "aci_l3_domain_profile" "test" {

		name  = aci_l3_domain_profile.test.name
		%s = "%s"
		depends_on = [ aci_l3_domain_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccL3DomainProfileDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3_domain_profile Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
		%s = "%s"
	}

	data "aci_l3_domain_profile" "test" {

		name  = aci_l3_domain_profile.test.name
		depends_on = [ aci_l3_domain_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccL3DomainProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing l3_domain_profile Data Source with invalid name")
	resource := fmt.Sprintf(`

	resource "aci_l3_domain_profile" "test" {

		name  = "%s"
	}

	data "aci_l3_domain_profile" "test" {

		name  = "${aci_l3_domain_profile.test.name}_invalid"
		depends_on = [ aci_l3_domain_profile.test ]
	}
	`, rName)
	return resource
}
