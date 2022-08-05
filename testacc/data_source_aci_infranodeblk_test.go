package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciNodeBlockDataSource_Basic(t *testing.T) {
	resourceName := "aci_node_block.test"
	dataSourceName := "data.aci_node_block.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciNodeBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateNodeBlockDSWithoutRequired(rName, rName, rName, "switch_association_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateNodeBlockDSWithoutRequired(rName, rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccNodeBlockConfigDataSource(rName, rName, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "switch_association_dn", resourceName, "switch_association_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "from_", resourceName, "from_"),
					resource.TestCheckResourceAttrPair(dataSourceName, "to_", resourceName, "to_"),
				),
			},
			{
				Config:      CreateAccNodeBlockDataSourceUpdate(rName, rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccNodeBlockDSWithInvalidName(rName, rName, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccNodeBlockDataSourceUpdatedResource(rName, rName, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccNodeBlockConfigDataSource(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}

	resource "aci_node_block" "test"{
  		switch_association_dn = aci_leaf_selector.test.id
  		name = "%s"
	}

	data "aci_node_block" "test" {
		switch_association_dn  = aci_node_block.test.switch_association_dn
		name  = aci_node_block.test.name
		depends_on = [ aci_node_block.test ]
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
	return resource
}

func CreateNodeBlockDSWithoutRequired(mgmtNodeGrpName, infrazoneNodeGrpName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing node_block Data Source without ", attrName)
	rBlock := `
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}

	resource "aci_node_block" "test"{
  		switch_association_dn = aci_leaf_selector.test.id
  		name = "%s"
	}
	`
	switch attrName {
	case "switch_association_dn":
		rBlock += `
	data "aci_node_block" "test" {
	#	switch_association_dn  = aci_node_block.test.switch_association_dn
		name  = aci_node_block.test.name
		depends_on = [ aci_node_block.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_node_block" "test" {
		switch_association_dn  = aci_node_block.test.switch_association_dn
	#	name  = aci_node_block.test.name
		depends_on = [ aci_node_block.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
}

func CreateAccNodeBlockDSWithInvalidName(mgmtNodeGrpName, infrazoneNodeGrpName, rName string) string {
	fmt.Println("=== STEP  testing node_block Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}

	resource "aci_node_block" "test"{
  		switch_association_dn = aci_leaf_selector.test.id
  		name = "%s"
	}

	data "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = "${aci_node_block.test.name}_invalid"
		depends_on = [ aci_node_block.test ]
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName)
	return resource
}

func CreateAccNodeBlockDataSourceUpdate(mgmtNodeGrpName, infrazoneNodeGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_block Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}

	resource "aci_node_block" "test"{
  		switch_association_dn = aci_leaf_selector.test.id
  		name = "%s"
	}

	data "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = aci_node_block.test.name
		%s = "%s"
		depends_on = [ aci_node_block.test ]
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName, key, value)
	return resource
}

func CreateAccNodeBlockDataSourceUpdatedResource(mgmtNodeGrpName, infrazoneNodeGrpName, rName, key, value string) string {
	fmt.Println("=== STEP  testing node_block Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}

  	resource "aci_leaf_selector" "test" {
    	name = "%s"
    	leaf_profile_dn = aci_leaf_profile.test.id
    	switch_association_type = "ALL"
  	}

	resource "aci_node_block" "test"{
  		switch_association_dn = aci_leaf_selector.test.id
  		name = "%s"
		%s = "%s"
	}

	data "aci_node_block" "test" {
		switch_association_dn  = aci_leaf_selector.test.id
		name  = aci_node_block.test.name
		depends_on = [ aci_node_block.test ]
	}
	`, mgmtNodeGrpName, infrazoneNodeGrpName, rName, key, value)
	return resource
}
