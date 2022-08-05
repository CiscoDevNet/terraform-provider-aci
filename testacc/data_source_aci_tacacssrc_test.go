package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSSourceDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_source.test"
	dataSourceName := "data.aci_tacacs_source.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSSourceDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSSourceDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSSourceConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "parent_dn", resourceName, "parent_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "incl.#", resourceName, "incl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "incl.0", resourceName, "incl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "incl.1", resourceName, "incl.1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "min_sev", resourceName, "min_sev"),
				),
			},
			{
				Config:      CreateAccTACACSSourceDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSSourceDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTACACSSourceDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSSourceConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing tacacs_source Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
		parent_dn   = "uni/fabric/moncommon"
		name  = "%s"
	}

	data "aci_tacacs_source" "test" {
	
		name  = aci_tacacs_source.test.name
		depends_on = [ aci_tacacs_source.test ]
	}
	`, rName)
	return resource
}

func CreateTACACSSourceDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_source Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_source" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_tacacs_source" "test" {
	
	#	name  = aci_tacacs_source.test.name
		depends_on = [ aci_tacacs_source.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSSourceDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_source Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_source" "test" {
	
		name  = "${aci_tacacs_source.test.name}_invalid"
		depends_on = [ aci_tacacs_source.test ]
	}
	`, rName)
	return resource
}

func CreateAccTACACSSourceDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_source Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_source" "test" {
	
		name  = aci_tacacs_source.test.name
		%s = "%s"
		depends_on = [ aci_tacacs_source.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccTACACSSourceDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_source Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_source" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_source" "test" {
	
		name  = aci_tacacs_source.test.name
		depends_on = [ aci_tacacs_source.test ]
	}
	`, rName, key, value)
	return resource
}
