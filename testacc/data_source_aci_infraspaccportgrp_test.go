package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSpinePortPolicyGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_spine_port_policy_group.test"
	dataSourceName := "data.aci_spine_port_policy_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSpinePortPolicyGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateSpinePortPolicyGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSpinePortPolicyGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccSpinePortPolicyGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSpinePortPolicyGroupDSWithInvalidName(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccSpinePortPolicyGroupDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSpinePortPolicyGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}

	data "aci_spine_port_policy_group" "test" {
	
		name  = aci_spine_port_policy_group.test.name
		depends_on = [ aci_spine_port_policy_group.test ]
	}
	`, rName)
	return resource
}

func CreateSpinePortPolicyGroupDSWithoutRequired(rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing spine_port_policy_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_spine_port_policy_group" "test" {
	
	#	name  = aci_spine_port_policy_group.test.name
		depends_on = [ aci_spine_port_policy_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName)
}

func CreateAccSpinePortPolicyGroupDSWithInvalidName(rName string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}

	data "aci_spine_port_policy_group" "test" {
	
		name  = "${aci_spine_port_policy_group.test.name}_invalid"
		depends_on = [ aci_spine_port_policy_group.test ]
	}
	`, rName)
	return resource
}

func CreateAccSpinePortPolicyGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
	}

	data "aci_spine_port_policy_group" "test" {
	
		name  = aci_spine_port_policy_group.test.name
		%s = "%s"
		depends_on = [ aci_spine_port_policy_group.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccSpinePortPolicyGroupDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing spine_port_policy_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_spine_port_policy_group" "test" {
	
		name  = "%s"
		%s = "%s"
	}

	data "aci_spine_port_policy_group" "test" {
	
		name  = aci_spine_port_policy_group.test.name
		depends_on = [ aci_spine_port_policy_group.test ]
	}
	`, rName, key, value)
	return resource
}
