package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSpineProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_spine_profile.test"
	dataSourceName := "data.aci_spine_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpineProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpineProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccSpineProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpineProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSpineProfileDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpineProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing spine_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_profile" "test" {
	
		name  = aci_spine_profile.test.name
		depends_on = [ aci_spine_profile.test ]
	}
	`, rName)
	return resource
}

func CreateSpineProfileDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_profile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_spine_profile" "test" {
	
	#	name  = aci_spine_profile.test.name
		depends_on = [ aci_spine_profile.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpineProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing spine_profile Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_profile" "test" {
	
		name  = "${aci_spine_profile.test.name}_invalid"
		depends_on = [ aci_spine_profile.test ]
	}
	`, rName)
	return resource
}

func CreateAccSpineProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
	}

	data "aci_spine_profile" "test" {
	
		name  = aci_spine_profile.test.name
		%s = "%s"
		depends_on = [ aci_spine_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSpineProfileDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_spine_profile" "test" {
	
		name  = aci_spine_profile.test.name
		depends_on = [ aci_spine_profile.test ]
	}
	`, rName, key, value)
	return resource
}
