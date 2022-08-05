package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL2InterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_l2_interface_policy.test"
	dataSourceName := "data.aci_l2_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL2InterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL2InterfacePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL2InterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "qinq", resourceName, "qinq"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vepa", resourceName, "vepa"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vlan_scope", resourceName, "vlan_scope"),
				),
			},
			{
				Config:      CreateAccL2InterfacePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL2InterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL2InterfacePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL2InterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing l2_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l2_interface_policy" "test" {
	
		name  = aci_l2_interface_policy.test.name
		depends_on = [ aci_l2_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccL2InterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing l2_interface_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l2_interface_policy" "test" {
	
		name  = "${aci_l2_interface_policy.test.name}_invalid"
		depends_on = [ aci_l2_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateL2InterfacePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l2_interface_policy creation without ", attrName)
	rBlock := `
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_l2_interface_policy" "test" {
	
	#	name  = "%s"
		depends_on = [ aci_l2_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL2InterfacePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l2_interface_policy" "test" {
	
		name  = aci_l2_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_l2_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccL2InterfacePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l2_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_l2_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_l2_interface_policy" "test" {
	
		name  = aci_l2_interface_policy.test.name
		depends_on = [ aci_l2_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}
