package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciSwitchAssociationDataSource_Basic(t *testing.T) {
	resourceName := "aci_leaf_selector.test"
	dataSourceName := "data.aci_leaf_selector.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciSwitchAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "leaf_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateSwitchAssociationDSWithoutRequired(rName, rName, "ALL", "switch_association_type"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccSwitchAssociationConfigDataSource(rName, rName, "ALL"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "leaf_profile_dn", resourceName, "leaf_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "switch_association_type", resourceName, "switch_association_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccSwitchAssociationDataSourceUpdate(rName, rName, "ALL", randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccSwitchAssociationDSWithInvalidName(rName, rName, "ALL"),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccSwitchAssociationDataSourceUpdatedResource(rName, rName, "ALL", "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccSwitchAssociationConfigDataSource(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  testing leaf_selector Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}

	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = aci_leaf_selector.test.name
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}

func CreateSwitchAssociationDSWithoutRequired(infraNodePName, rName, switch_association_type, attrName string) string {
	fmt.Println("=== STEP  Basic: testing leaf_selector Data Source without ", attrName)
	rBlock := `
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}
	`
	switch attrName {
	case "leaf_profile_dn":
		rBlock += `
	data "aci_leaf_selector" "test" {
	#	leaf_profile_dn  = aci_leaf_profile.test.id
		name  = aci_leaf_selector.test.name	
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
		`
	case "name":
		rBlock += `
	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
	#	name  = aci_leaf_selector.test.name
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
		`
	case "switch_association_type":
		rBlock += `
	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = aci_leaf_selector.test.name
	#	switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, infraNodePName, rName, switch_association_type)
}

func CreateAccSwitchAssociationDSWithInvalidName(infraNodePName, rName, switch_association_type string) string {
	fmt.Println("=== STEP  testing leaf_selector Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}

	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "${aci_leaf_selector.test.name}_invalid"
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
	`, infraNodePName, rName, switch_association_type)
	return resource
}

func CreateAccSwitchAssociationDataSourceUpdate(infraNodePName, rName, switch_association_type, key, value string) string {
	fmt.Println("=== STEP  testing leaf_selector Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
	}

	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = aci_leaf_selector.test.name
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		%s = "%s"
		depends_on = [ aci_leaf_selector.test ]
	}
	`, infraNodePName, rName, switch_association_type, key, value)
	return resource
}

func CreateAccSwitchAssociationDataSourceUpdatedResource(infraNodePName, rName, switch_association_type, key, value string) string {
	fmt.Println("=== STEP  testing leaf_selector Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_leaf_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = "%s"
		switch_association_type  = "%s"
		%s = "%s"
	}

	data "aci_leaf_selector" "test" {
		leaf_profile_dn  = aci_leaf_profile.test.id
		name  = aci_leaf_selector.test.name
		switch_association_type  = aci_leaf_selector.test.switch_association_type
		depends_on = [ aci_leaf_selector.test ]
	}
	`, infraNodePName, rName, switch_association_type, key, value)
	return resource
}
