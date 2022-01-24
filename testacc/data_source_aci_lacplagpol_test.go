package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLACPPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_lacp_policy.test"
	dataSourceName := "data.aci_lacp_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLACPPolicyDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLACPPolicyDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLACPPolicyConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.0", resourceName, "ctrl.0"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.1", resourceName, "ctrl.1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.2", resourceName, "ctrl.2"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_links", resourceName, "max_links"),
					resource.TestCheckResourceAttrPair(dataSourceName, "min_links", resourceName, "min_links"),
					resource.TestCheckResourceAttrPair(dataSourceName, "mode", resourceName, "mode"),
				),
			},
			{
				Config:      CreateAccLACPPolicyDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccLACPPolicyDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLACPPolicyDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateLACPPolicyDSWithoutRequired(rName, attr string) string {
	fmt.Println("=== STEP  testing lacp_policy Data Source without required argument")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lacp_policy" "test" {
		# name  = aci_lacp_policy.test.name
		depends_on = [
			aci_lacp_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLACPPolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing lacp_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lacp_policy" "test" {
	
		name  = aci_lacp_policy.test.name
		depends_on = [
			aci_lacp_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLACPPolicyDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing lacp_policy Data Source with invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lacp_policy" "test" {
	
		name  = "${aci_lacp_policy.test.name}_invalid"
		depends_on = [
			aci_lacp_policy.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLACPPolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing lacp_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
	}

	data "aci_lacp_policy" "test" {
	
		name  = aci_lacp_policy.test.name
		%s = "%s"
		depends_on = [
			aci_lacp_policy.test
		]
	}
	`, rName, key, value)
	return resource
}

func CreateAccLACPPolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing lacp_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_lacp_policy" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_lacp_policy" "test" {
	
		name  = aci_lacp_policy.test.name
		depends_on = [
			aci_lacp_policy.test
		]
	}
	`, rName, key, value)
	return resource
}
