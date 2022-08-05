package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSwitchSpineAssociationDataSource_Basic(t *testing.T) {
	resourceName := "aci_spine_switch_association.test"
	dataSourceName := "data.aci_spine_switch_association.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchSpineAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSpineSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "spine_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSpineSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSpineSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "spine_switch_association_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpineSwitchAssociationConfigDataSource(rName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "spine_profile_dn", resourceName, "spine_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spine_switch_association_type", resourceName, "spine_switch_association_type"),
				),
			},
			{
				Config:      CreateAccSpineSwitchAssociationDataSourceUpdate(rName, rName, "ALL", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpineSwitchAssociationDSWithInvalidParentDn(rName, rName, "ALL"),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},

			{
				Config: CreateAccSpineSwitchAssociationDataSourceUpdatedResource(rName, rName, "ALL", "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpineSwitchAssociationConfigDataSource(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing spine_switch_association Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}

	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = aci_spine_switch_association.test.name
		spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		depends_on = [ aci_spine_switch_association.test ]
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}

func CreateSpineSwitchAssociationDSWithoutRequired(infraSpinePName, rName, spine_switch_association_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_switch_association Data Source without ", attrName)
	rBlock := `
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}
	`
	switch attrName {
	case "spine_profile_dn":
		rBlock += `
	data "aci_spine_switch_association" "test" {
	#	spine_profile_dn  = aci_spine_profile.test.id
		name  = aci_spine_switch_association.test.name	
		spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		depends_on = [ aci_spine_switch_association.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
	#	name  = aci_spine_switch_association.test.name
		spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		depends_on = [ aci_spine_switch_association.test ]
	}
		`
	case "spine_switch_association_type":
		rBlock += `
	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = aci_spine_switch_association.test.name
	#	spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		depends_on = [ aci_spine_switch_association.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraSpinePName, rName, spine_switch_association_type)
}

func CreateAccSpineSwitchAssociationDSWithInvalidParentDn(infraSpinePName, rName, spine_switch_association_type string) string {
	fmt.Println("=== STEP  testing spine_switch_association Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}

	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "${aci_spine_switch_association.test.name}_invalid"
		spine_switch_association_type  = "${aci_spine_switch_association.test.spine_switch_association_type}_invalid"
		depends_on = [ aci_spine_switch_association.test ]
	}
	`, infraSpinePName, rName, spine_switch_association_type)
	return resource
}

func CreateAccSpineSwitchAssociationDataSourceUpdate(infraSpinePName, rName, spine_switch_association_type, key, value string) string {
	fmt.Println("=== STEP  testing spine_switch_association Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
	}

	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = aci_spine_switch_association.test.name
		spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		%s = "%s"
		depends_on = [ aci_spine_switch_association.test ]
	}
	`, infraSpinePName, rName, spine_switch_association_type, key, value)
	return resource
}

func CreateAccSpineSwitchAssociationDataSourceUpdatedResource(infraSpinePName, rName, spine_switch_association_type, key, value string) string {
	fmt.Println("=== STEP  testing spine_switch_association Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = "%s"
		spine_switch_association_type  = "%s"
		%s = "%s"
	}

	data "aci_spine_switch_association" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		name  = aci_spine_switch_association.test.name
		spine_switch_association_type  = aci_spine_switch_association.test.spine_switch_association_type
		depends_on = [ aci_spine_switch_association.test ]
	}
	`, infraSpinePName, rName, spine_switch_association_type, key, value)
	return resource
}
