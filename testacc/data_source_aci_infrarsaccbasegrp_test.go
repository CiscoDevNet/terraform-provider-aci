package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciAccessGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_access_group.test"
	dataSourceName := "data.aci_access_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciAccessGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccessGroupDSWithoutRequired(rName, rName, "access_port_selector_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccAccessGroupConfigDataSource(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "access_port_selector_dn", resourceName, "access_port_selector_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fex_id", resourceName, "fex_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "t_dn", resourceName, "t_dn"),
				),
			},
			{
				Config:      CreateAccAccessGroupDataSourceUpdate(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccAccessGroupDSWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`Invalid RN`),
			},

			{
				Config: CreateAccAccessGroupDataSourceUpdatedResource(rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccAccessGroupConfigDataSource(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  testing access_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}

	data "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		depends_on = [ aci_access_group.test ]
	}
	`, infraAccPortPName, infraHPortSName)
	return resource
}

func CreateAccessGroupDSWithoutRequired(infraAccPortPName, infraHPortSName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing access_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}
	`
	switch attrName {
	case "access_port_selector_dn":
		rBlock += `
	data "aci_access_group" "test" {
	#	access_port_selector_dn  = aci_access_port_selector.test.id
	
		depends_on = [ aci_access_group.test ]
	}
		`

	}
	return fmt.Sprintf(rBlock, infraAccPortPName, infraHPortSName)
}

func CreateAccAccessGroupDSWithInvalidParentDn(infraAccPortPName, infraHPortSName string) string {
	fmt.Println("=== STEP  testing access_group Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}

	data "aci_access_group" "test" {
		access_port_selector_dn  = "${aci_access_port_selector.test.id}_invalid"
		depends_on = [ aci_access_group.test ]
	}
	`, infraAccPortPName, infraHPortSName)
	return resource
}

func CreateAccAccessGroupDataSourceUpdate(infraAccPortPName, infraHPortSName, key, value string) string {
	fmt.Println("=== STEP  testing access_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
	}

	data "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		%s = "%s"
		depends_on = [ aci_access_group.test ]
	}
	`, infraAccPortPName, infraHPortSName, key, value)
	return resource
}

func CreateAccAccessGroupDataSourceUpdatedResource(infraAccPortPName, infraHPortSName, key, value string) string {
	fmt.Println("=== STEP  testing access_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_interface_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_access_port_selector" "test" {
		name 		= "%s"
		leaf_interface_profile_dn = aci_leaf_interface_profile.test.id
		access_port_selector_type = "ALL"
	}
	
	resource "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		%s = "%s"
	}

	data "aci_access_group" "test" {
		access_port_selector_dn  = aci_access_port_selector.test.id
		depends_on = [ aci_access_group.test ]
	}
	`, infraAccPortPName, infraHPortSName, key, value)
	return resource
}
