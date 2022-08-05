package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciFCDomainDataSource_Basic(t *testing.T) {
	resourceName := "aci_fc_domain.test"
	dataSourceName := "data.aci_fc_domain.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateFCDomainDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFCDomainConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccFCDomainDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccFCDomainDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccFCDomainDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccFCDomainConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing fc_domain Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}

	data "aci_fc_domain" "test" {
	
		name  = aci_fc_domain.test.name
		depends_on = [ aci_fc_domain.test ]
	}
	`, rName)
	return resource
}

func CreateFCDomainDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing fc_domain Data Source without ", attrName)
	rBlock := `
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_fc_domain" "test" {
	
	#	name  = aci_fc_domain.test.name
		depends_on = [ aci_fc_domain.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccFCDomainDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing fc_domain Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}

	data "aci_fc_domain" "test" {
	
		name  = "${aci_fc_domain.test.name}_invalid"
		depends_on = [ aci_fc_domain.test ]
	}
	`, rName)
	return resource
}

func CreateAccFCDomainDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing fc_domain Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
	}

	data "aci_fc_domain" "test" {
	
		name  = aci_fc_domain.test.name
		%s = "%s"
		depends_on = [ aci_fc_domain.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccFCDomainDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing fc_domain Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_fc_domain" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_fc_domain" "test" {
	
		name  = aci_fc_domain.test.name
		depends_on = [ aci_fc_domain.test ]
	}
	`, rName, key, value)
	return resource
}
