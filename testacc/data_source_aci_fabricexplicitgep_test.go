package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVPCExplicitProtectionGroupDataSource_Basic(t *testing.T) {
	resourceName := "aci_vpc_explicit_protection_group.test"
	dataSourceName := "data.aci_vpc_explicit_protection_group.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	switch1 := "501"
	switch2 := "502"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVPCExplicitProtectionGroupDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateVPCExplicitProtectionGroupDSWithoutRequired(switch1, switch2, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupConfigDataSource(switch1, switch2, rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "vpc_explicit_protection_group_id", resourceName, "vpc_explicit_protection_group_id"),
				),
			},
			{
				Config:      CreateAccVPCExplicitProtectionGroupDataSourceUpdate(switch1, switch2, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVPCExplicitProtectionGroupDSWithInvalidName(switch1, switch2, rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccVPCExplicitProtectionGroupDataSourceUpdatedResource(switch1, switch2, rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVPCExplicitProtectionGroupConfigDataSource(sw1, sw2, rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}

	data "aci_vpc_explicit_protection_group" "test" {
	
		name  = aci_vpc_explicit_protection_group.test.name
		depends_on = [ aci_vpc_explicit_protection_group.test ]
	}
	`, sw1, sw2, rName)
	return resource
}

func CreateVPCExplicitProtectionGroupDSWithoutRequired(sw1, sw2, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vpc_explicit_protection_group Data Source without ", attrName)
	rBlock := `
	
	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}
	`
	switch attrName {
	case "name":
		rBlock += `
	data "aci_vpc_explicit_protection_group" "test" {
	
	#	name  = aci_vpc_explicit_protection_group.test.name
		depends_on = [ aci_vpc_explicit_protection_group.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, sw1, sw2, rName)
}

func CreateAccVPCExplicitProtectionGroupDSWithInvalidName(sw1, sw2, rName string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group Data Source with invalid name")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}

	data "aci_vpc_explicit_protection_group" "test" {
	
		name  = "${aci_vpc_explicit_protection_group.test.name}_invalid"
		depends_on = [ aci_vpc_explicit_protection_group.test ]
	}
	`, sw1, sw2, rName)
	return resource
}

func CreateAccVPCExplicitProtectionGroupDataSourceUpdate(sw1, sw2, rName, key, value string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
	}

	data "aci_vpc_explicit_protection_group" "test" {
	
		name  = aci_vpc_explicit_protection_group.test.name
		%s = "%s"
		depends_on = [ aci_vpc_explicit_protection_group.test ]
	}
	`, sw1, sw2, rName, key, value)
	return resource
}

func CreateAccVPCExplicitProtectionGroupDataSourceUpdatedResource(sw1, sw2, rName, key, value string) string {
	fmt.Println("=== STEP  testing vpc_explicit_protection_group Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_vpc_explicit_protection_group" "test" {
		switch1 = "%s"
		switch2 = "%s"
		name  = "%s"
		%s = "%s"
	}

	data "aci_vpc_explicit_protection_group" "test" {
	
		name  = aci_vpc_explicit_protection_group.test.name
		depends_on = [ aci_vpc_explicit_protection_group.test ]
	}
	`, sw1, sw2, rName, key, value)
	return resource
}
