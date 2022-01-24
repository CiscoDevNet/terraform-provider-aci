package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAttachableAccessEntityProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_attachable_access_entity_profile.test"
	dataSourceName := "data.aci_attachable_access_entity_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAttachableAccessEntityProfileDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateAttachableAccessEntityProfileDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccAttachableAccessEntityProfileDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAttachableAccessEntityProfileDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccAttachableAccessEntityProfileDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccAttachableAccessEntityProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}
	
	data "aci_attachable_access_entity_profile" "test" {
		name  = aci_attachable_access_entity_profile.test.name
		depends_on = [ aci_attachable_access_entity_profile.test ]
	}
	`, rName)
	return resource
}

func CreateAccAttachableAccessEntityProfileDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile Data Source with Invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}

	data "aci_attachable_access_entity_profile" "test" {
		name  = "${aci_attachable_access_entity_profile.test.name}_invalid"
		depends_on = [ aci_attachable_access_entity_profile.test ]
	}
	`, rName)
	return resource
}

func CreateAttachableAccessEntityProfileDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing attachable_access_entity_profile Data Source without ", attrName)
	rBlock := `
	resource "aci_attachable_access_entity_profile" "test" {
	name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_attachable_access_entity_profile" "test" {
	#	name  = "%s"
		depends_on = [ aci_attachable_access_entity_profile.test ]
	}
	`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccAttachableAccessEntityProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name  = "%s"
	}

	data "aci_attachable_access_entity_profile" "test" {
		name  = aci_attachable_access_entity_profile.test.name
		%s = "%s"
		depends_on = [ aci_attachable_access_entity_profile.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccAttachableAccessEntityProfileDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing attachable_access_entity_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_attachable_access_entity_profile" "test" {	
		name  = aci_attachable_access_entity_profile.test.name
		depends_on = [ aci_attachable_access_entity_profile.test ]
	}
	`, rName, key, value)
	return resource
}
