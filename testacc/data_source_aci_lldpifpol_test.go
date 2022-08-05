package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLLDPInterfacePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_lldp_interface_policy.test"
	dataSourceName := "data.aci_lldp_interface_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLLDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLLDPInterfacePolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLLDPInterfacePolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_rx_st", resourceName, "admin_rx_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_tx_st", resourceName, "admin_tx_st"),
				),
			},
			{
				Config:      CreateAccLLDPInterfacePolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLLDPInterfacePolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccLLDPInterfacePolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLLDPInterfacePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lldp_interface_policy" "test" {
	
		name  = aci_lldp_interface_policy.test.name
		depends_on = [
			aci_lldp_interface_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy Data Source with invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lldp_interface_policy" "test" {
	
		name  = "${aci_lldp_interface_policy.test.name}_invalid"
		depends_on = [
			aci_lldp_interface_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateLLDPInterfacePolicyDSWithoutRequired(rName, attr string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy Data Source without required argument")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lldp_interface_policy" "test" {
	
		depends_on = [
			aci_lldp_interface_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLLDPInterfacePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lldp_interface_policy" "test" {
	
		name  = aci_lldp_interface_policy.test.name
		%s = "%s"
		depends_on = [
			aci_lldp_interface_policy.test
		]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLLDPInterfacePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing lldp_interface_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_lldp_interface_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_lldp_interface_policy" "test" {
	
		name  = aci_lldp_interface_policy.test.name
		depends_on = [
			aci_lldp_interface_policy.test
		]
	}
	`, rName, key, value)
	return resource
}
