package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL2DomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_l2_domain.test"
	dataSourceName := "data.aci_l2_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL2DomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2DomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccL2DomainDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL2DomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL2DomainDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL2DomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing l2_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}

	data "aci_l2_domain" "test" {
	
		name  = aci_l2_domain.test.name
		depends_on = [ aci_l2_domain.test ]
	}
	`, rName)
	return resource
}

func CreateL2DomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_domain Data Source without ", attrName)
	rBlock := `
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_l2_domain" "test" {
	
	#	name  = aci_l2_domain.test.name
		depends_on = [ aci_l2_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL2DomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing l2_domain Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}

	data "aci_l2_domain" "test" {
	
		name  = "${aci_l2_domain.test.name}_invalid"
		depends_on = [ aci_l2_domain.test ]
	}
	`, rName)
	return resource
}

func CreateAccL2DomainDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_domain Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
	}

	data "aci_l2_domain" "test" {
	
		name  = aci_l2_domain.test.name
		%s = "%s"
		depends_on = [ aci_l2_domain.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccL2DomainDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_domain Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_domain" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_l2_domain" "test" {
	
		name  = aci_l2_domain.test.name
		depends_on = [ aci_l2_domain.test ]
	}
	`, rName, key, value)
	return resource
}
