package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3InterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3_interface_policy.test"
	dataSourceName := "data.aci_l3_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3InterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3InterfacePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3InterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bfd_isis", resourceName, "bfd_isis"),
				),
			},
			{
				Config:      CreateAccL3InterfacePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3InterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3InterfacePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3InterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing l3_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l3_interface_policy" "test" {
	
		name  = aci_l3_interface_policy.test.name
		depends_on = [ aci_l3_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateL3InterfacePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3_interface_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_l3_interface_policy" "test" {
	
	#	name  = aci_l3_interface_policy.test.name
		depends_on = [ aci_l3_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccL3InterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing l3_interface_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l3_interface_policy" "test" {
	
		name  = "${aci_l3_interface_policy.test.name}_invalid"
		depends_on = [ aci_l3_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccL3InterfacePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_l3_interface_policy" "test" {
	
		name  = aci_l3_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_l3_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccL3InterfacePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing l3_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_l3_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_l3_interface_policy" "test" {
	
		name  = aci_l3_interface_policy.test.name
		depends_on = [ aci_l3_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}
