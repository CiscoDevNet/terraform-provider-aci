package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMiscablingProtocolInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_miscabling_protocol_interface_policy.test"
	dataSourceName := "data.aci_miscabling_protocol_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMiscablingProtocolInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateMiscablingProtocolInterfacePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
				),
			},
			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccMiscablingProtocolInterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccMiscablingProtocolInterfacePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMiscablingProtocolInterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = aci_miscabling_protocol_interface_policy.test.name
		depends_on = [ aci_miscabling_protocol_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "${aci_miscabling_protocol_interface_policy.test.name}_invalid"
		depends_on = [ aci_miscabling_protocol_interface_policy.test ]
	}
	`, rName)
	return resource
}

func CreateMiscablingProtocolInterfacePolicyDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing miscabling_protocol_interface_policy Data Source without ", attrName)
	rBlock := `
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_miscabling_protocol_interface_policy" "test" {
	
	#	name  = "%s"
		depends_on = [ aci_miscabling_protocol_interface_policy.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccMiscablingProtocolInterfacePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = aci_miscabling_protocol_interface_policy.test.name
		%s = "%s"
		depends_on = [ aci_miscabling_protocol_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccMiscablingProtocolInterfacePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing miscabling_protocol_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_miscabling_protocol_interface_policy" "test" {
	
		name  = aci_miscabling_protocol_interface_policy.test.name
		depends_on = [ aci_miscabling_protocol_interface_policy.test ]
	}
	`, rName, key, value)
	return resource
}
