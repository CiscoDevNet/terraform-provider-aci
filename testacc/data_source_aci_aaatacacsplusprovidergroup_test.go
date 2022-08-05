package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciTACACSPlusProviderGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_tacacs_provider_group.test"
	dataSourceName := "data.aci_tacacs_provider_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciTACACSPlusProviderGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateTACACSPlusProviderGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccTACACSPlusProviderGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccTACACSPlusProviderGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccTACACSPlusProviderGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccTACACSPlusProviderGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_provider_group" "test" {
	
		name  = aci_tacacs_provider_group.test.name
		depends_on = [ aci_tacacs_provider_group.test ]
	}
	`, rName)
	return resource
}

func CreateTACACSPlusProviderGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing tacacs_provider_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_tacacs_provider_group" "test" {
	
	#	name  = aci_tacacs_provider_group.test.name
		depends_on = [ aci_tacacs_provider_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccTACACSPlusProviderGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_provider_group" "test" {
	
		name  = "${aci_tacacs_provider_group.test.name}_invalid"
		depends_on = [ aci_tacacs_provider_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccTACACSPlusProviderGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
	}

	data "aci_tacacs_provider_group" "test" {
	
		name  = aci_tacacs_provider_group.test.name
		%s = "%s"
		depends_on = [ aci_tacacs_provider_group.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccTACACSPlusProviderGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing tacacs_provider_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tacacs_provider_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_tacacs_provider_group" "test" {
	
		name  = aci_tacacs_provider_group.test.name
		depends_on = [ aci_tacacs_provider_group.test ]
	}
	`, rName, key, value)
	return resource
}
