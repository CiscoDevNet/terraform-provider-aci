package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAccessGenericDataSource_Basic(t *testing.T) {
	resourceName := "aci_access_generic.test"
	dataSourceName := "data.aci_access_generic.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGenericDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessGenericDSWithoutRequired(rName, rName, "attachable_access_entity_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccessGenericDSWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGenericConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "attachable_access_entity_profile_dn", resourceName, "attachable_access_entity_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccAccessGenericDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAccessGenericDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccAccessGenericDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccAccessGenericConfigDataSource(infraAttEntityPName, rName string) string {
	fmt.Println("=== STEP  testing access_generic Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "default"
	}

	data "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = aci_access_generic.test.name
		depends_on = [ aci_access_generic.test ]
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccessGenericDSWithoutRequired(infraAttEntityPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_generic Data Source without ", attrName)
	rBlock := `
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "default"
	}
	`
	switch attrName {
	case "attachable_access_entity_profile_dn":
		rBlock += `
	data "aci_access_generic" "test" {
	#	attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = aci_access_generic.test.name
		depends_on = [ aci_access_generic.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	#	name  = aci_access_generic.test.name
		depends_on = [ aci_access_generic.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAttEntityPName)
}

func CreateAccAccessGenericDSWithInvalidParentDn(infraAttEntityPName, rName string) string {
	fmt.Println("=== STEP  testing access_generic Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "default"
	}

	data "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "${aci_access_generic.test.name}_invalid"
		depends_on = [ aci_access_generic.test ]
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccAccessGenericDataSourceUpdate(infraAttEntityPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing access_generic Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "default"
	}

	data "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = aci_access_generic.test.name
		%s = "%s"
		depends_on = [ aci_access_generic.test ]
	}
	`, infraAttEntityPName, key, value)
	return resource
}

func CreateAccAccessGenericDataSourceUpdatedResource(infraAttEntityPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing access_generic Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = "default"
		%s = "%s"
	}

	data "aci_access_generic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		name  = aci_access_generic.test.name
		depends_on = [ aci_access_generic.test ]
	}
	`, infraAttEntityPName, key, value)
	return resource
}
