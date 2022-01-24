package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciLeafBreakoutPortGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_leaf_breakout_port_group.test"
	dataSourceName := "data.aci_leaf_breakout_port_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLeafBreakoutPortGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateLeafBreakoutPortGroupDSWithoutRequired(rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLeafBreakoutPortGroupConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "brkout_map", resourceName, "brkout_map"),
				),
			},
			{
				Config:      CreateAccLeafBreakoutPortGroupDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccLeafBreakoutPortGroupDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccLeafBreakoutPortGroupConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}

	data "aci_leaf_breakout_port_group" "test" {
	
		name  = aci_leaf_breakout_port_group.test.name
		depends_on = [
			aci_leaf_breakout_port_group.test
		]
	}
	`, rName)
	return resource
}

func CreateLeafBreakoutPortGroupDSWithoutRequired(rName, attr string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group Data Source without required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}

	data "aci_leaf_breakout_port_group" "test" {
	
		depends_on = [
			aci_leaf_breakout_port_group.test
		]
	}
	`, rName)
	return resource
}

func CreateAccLeafBreakoutPortGroupDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing leaf_breakout_port_group creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_breakout_port_group" "test" {
	
		name  = "%s"
	}

	data "aci_leaf_breakout_port_group" "test" {
	
		name  = aci_leaf_breakout_port_group.test.name
		%s = "%s"
		depends_on = [
			aci_leaf_breakout_port_group.test
		]
	}
	`, rName, key, value)
	return resource
}
