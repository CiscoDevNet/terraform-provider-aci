package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSpineInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_spine_interface_profile.test"
	dataSourceName := "data.aci_spine_interface_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineInterfaceProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpineInterfaceProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineInterfaceProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccSpineInterfaceProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpineInterfaceProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSpineInterfaceProfileDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpineInterfaceProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing spine_interface_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_interface_profile" "test" {
	
		name  = aci_spine_interface_profile.test.name
		depends_on = [ aci_spine_interface_profile.test ]
	}
	`, rName)
	return resource
}

func CreateSpineInterfaceProfileDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_interface_profile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_spine_interface_profile" "test" {
	
	#	name  = aci_spine_interface_profile.test.name
		depends_on = [ aci_spine_interface_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpineInterfaceProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing spine_interface_profile Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_interface_profile" "test" {
	
		name  = "${aci_spine_interface_profile.test.name}_invalid"
		depends_on = [ aci_spine_interface_profile.test ]
	}
	`, rName)
	return resource
}

func CreateAccSpineInterfaceProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_interface_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_interface_profile" "test" {
	
		name  = aci_spine_interface_profile.test.name
		%s = "%s"
		depends_on = [ aci_spine_interface_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSpineInterfaceProfileDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_interface_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_interface_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_spine_interface_profile" "test" {
	
		name  = aci_spine_interface_profile.test.name
		depends_on = [ aci_spine_interface_profile.test ]
	}
	`, rName, key, value)
	return resource
}
