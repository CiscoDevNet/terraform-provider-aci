package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAccessPortSelectorDataSource_Basic(t *testing.T) {
	resourceName := "aci_access_port_selector.test"
	dataSourceName := "data.aci_access_port_selector.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	accessPortSelectorType := "ALL"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessPortSelectorDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessPortSelectorDSWithoutRequired(rName, rName, accessPortSelectorType, "leaf_interface_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccessPortSelectorDSWithoutRequired(rName, rName, accessPortSelectorType, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessPortSelectorConfigDataSource(rName, rName, accessPortSelectorType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "leaf_interface_profile_dn", resourceName, "leaf_interface_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "access_port_selector_type", resourceName, "access_port_selector_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccAccessPortSelectorDataSourceUpdate(rName, rName, accessPortSelectorType, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAccessPortSelectorDSWithInvalidParentDn(rName, rName, accessPortSelectorType),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccAccessPortSelectorDataSourceUpdatedResource(rName, rName, accessPortSelectorType, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccAccessPortSelectorConfigDataSource(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing access_port_selector Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}

	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = aci_access_port_selector.test.name
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccessPortSelectorDSWithoutRequired(infraAccPortPName, rName, access_port_selector_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_port_selector Data Source without ", attrName)
	rBlock := `
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}
	`
	switch attrName {
	case "leaf_interface_profile_dn":
		rBlock += `
	data "aci_access_port_selector" "test" {
	#	leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = aci_access_port_selector.test.name
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
	#	name  = aci_access_port_selector.test.name
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
		`
	case "access_port_selector_type":
		rBlock += `
	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = aci_access_port_selector.test.name
	#	access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraAccPortPName, rName, access_port_selector_type)
}

func CreateAccAccessPortSelectorDSWithInvalidParentDn(infraAccPortPName, rName, access_port_selector_type string) string {
	fmt.Println("=== STEP  testing access_port_selector Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}

	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "${aci_access_port_selector.test.name}_invalid"
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
	`, infraAccPortPName, rName, access_port_selector_type)
	return resource
}

func CreateAccAccessPortSelectorDataSourceUpdate(infraAccPortPName, rName, access_port_selector_type, key, value string) string {
	fmt.Println("=== STEP  testing access_port_selector Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
	}

	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = aci_access_port_selector.test.name
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		%s = "%s"
		depends_on = [ aci_access_port_selector.test ]
	}
	`, infraAccPortPName, rName, access_port_selector_type, key, value)
	return resource
}

func CreateAccAccessPortSelectorDataSourceUpdatedResource(infraAccPortPName, rName, access_port_selector_type, key, value string) string {
	fmt.Println("=== STEP  testing access_port_selector Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = "%s"
		access_port_selector_type  = "%s"
		%s = "%s"
	}

	data "aci_access_port_selector" "test" {
		leaf_interface_profile_dn  = aci_leaf_interface_profile.test.id
		name  = aci_access_port_selector.test.name
		access_port_selector_type  = aci_access_port_selector.test.access_port_selector_type
		depends_on = [ aci_access_port_selector.test ]
	}
	`, infraAccPortPName, rName, access_port_selector_type, key, value)
	return resource
}
