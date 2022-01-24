package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSpinePortSelectorDataSource_Basic(t *testing.T) {
	resourceName := "aci_spine_port_selector.test"
	dataSourceName := "data.aci_spine_port_selector.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	tDn := "aci_spine_interface_profile.test.id"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSpinePortSelectorDSWithoutRequired(rName, rName, tDn, "spine_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSpinePortSelectorDSWithoutRequired(rName, rName, tDn, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpinePortSelectorConfigDataSource(rName, rName, tDn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "spine_profile_dn", resourceName, "spine_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tdn", resourceName, "tdn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config:      CreateAccSpinePortSelectorDataSourceUpdate(rName, rName, tDn, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpinePortSelectorDSWithInvalidParentDn(rName, rName, tDn),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccSpinePortSelectorDataSourceUpdatedResource(rName, rName, tDn, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpinePortSelectorConfigDataSource(infraSpinePName, rName, tDn string) string {
	fmt.Println("=== STEP  testing spine_port_selector Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}

	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = %s
	}

	data "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_port_selector.test.tdn
		depends_on = [ aci_spine_port_selector.test ]
	}
	`, infraSpinePName, rName, tDn)
	return resource
}

func CreateSpinePortSelectorDSWithoutRequired(infraSpinePName, rName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_selector Data Source without ", attrName)
	rBlock := `
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}

	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
	}
	`
	switch attrName {
	case "spine_profile_dn":
		rBlock += `
	data "aci_spine_port_selector" "test" {
	#	spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_port_selector.test.tdn
		depends_on = [ aci_spine_port_selector.test ]
	}
		`
	case "tdn":
		rBlock += `
	data "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
	#	tdn  = aci_spine_port_selector.test.tdn
		depends_on = [ aci_spine_port_selector.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraSpinePName, tDn)
}

func CreateAccSpinePortSelectorDSWithInvalidParentDn(infraSpinePName, rName, tDn string) string {
	fmt.Println("=== STEP  testing spine_port_selector Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}

	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
	}

	data "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = "${aci_spine_port_selector.test.tdn}_invalid"
		depends_on = [ aci_spine_port_selector.test ]
	}
	`, infraSpinePName, tDn)
	return resource
}

func CreateAccSpinePortSelectorDataSourceUpdate(infraSpinePName, rName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing spine_port_selector Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}

	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}
	
	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
	}

	data "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_port_selector.test.tdn
		%s = "%s"
		depends_on = [ aci_spine_port_selector.test ]
	}
	`, infraSpinePName, tDn, key, value)
	return resource
}

func CreateAccSpinePortSelectorDataSourceUpdatedResource(infraSpinePName, rName, tDn, key, value string) string {
	fmt.Println("=== STEP  testing spine_port_selector Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_spine_interface_profile" "test" {
		name        = "%s"
	}

	resource "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_interface_profile.test.id
		%s = "%s"
	}

	data "aci_spine_port_selector" "test" {
		spine_profile_dn  = aci_spine_profile.test.id
		tdn  = aci_spine_port_selector.test.tdn
		depends_on = [ aci_spine_port_selector.test ]
	}
	`, infraSpinePName, tDn, key, value)
	return resource
}
