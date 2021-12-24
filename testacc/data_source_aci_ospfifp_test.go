package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outOspfInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_ospf_interface_profile.test"
	dataSourceName := "data.aci_l3out_ospf_interface_profile.test"
	rName := acctest.RandString(5)
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outOspfInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateAccL3outOspfInterfaceProfileDSWithoutLogicalInterfaceProfile(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileDSWithoutAuthKey(rName),
				ExpectError: regexp.MustCompile(`Object may not exists`),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_key", resourceName, "auth_key"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_key_id", resourceName, "auth_key_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "auth_type", resourceName, "auth_type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_interface_profile_dn", resourceName, "logical_interface_profile_dn"),
				),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileUpdatedConfigDataSourceRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outOspfInterfaceProfileUpdatedConfigDataSource(rName, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
				),
			},
			{
				Config:      CreateAccL3outOspfInterfaceProfileDSWithInvalidLogicalInterfaceProfileDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccL3outOspfInterfaceProfileUpdatedConfigDataSourceRandomAttr(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing L3out Ospf Interface Profile data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = aci_l3out_ospf_interface_profile.test.logical_interface_profile_dn
		auth_key = aci_l3out_ospf_interface_profile.test.auth_key
		%s = "%s"
	}
	`, rName, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccL3outOspfInterfaceProfileUpdatedConfigDataSource(rName, attribute, value string) string {
	fmt.Println("=== STEP  Basic: Testing L3out Ospf Interface Profile data source with updated resource")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
		%s = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = aci_l3out_ospf_interface_profile.test.logical_interface_profile_dn
		auth_key = aci_l3out_ospf_interface_profile.test.auth_key
	}

	`, rName, rName, rName, rName, rName, attribute, value)
	return resource
}

func CreateAccL3outOspfInterfaceProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Basic: L3out Ospf Interface Profile subject data source")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = aci_l3out_ospf_interface_profile.test.logical_interface_profile_dn
		auth_key = aci_l3out_ospf_interface_profile.test.auth_key
	}
	`, rName, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outOspfInterfaceProfileDSWithInvalidLogicalInterfaceProfileDn(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3out Ospf Interface Profile reading with Invalid Logical Interface Profile Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = "${aci_l3out_ospf_interface_profile.test.logical_interface_profile_dn}abc"
		auth_key = aci_l3out_ospf_interface_profile.test.auth_key
	}

	`, rName, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outOspfInterfaceProfileDSWithoutLogicalInterfaceProfile(rName string) string {
	fmt.Println("=== STEP  Basic: testing L3out Ospf Interface Profile reading without giving Logical Interface Profile")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		auth_key = aci_l3out_ospf_interface_profile.test.auth_key
	}
	`, rName, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outOspfInterfaceProfileDSWithoutAuthKey(rName string) string {
	fmt.Println("=== STEP  Basic: testing contract subject reading without giving auth key")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		description = "l3_outside created while acceptance testing"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		description = "logical_node_profile created while acceptance testing"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		description = "logical_interface_profile created while acceptance testing"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		auth_key = "%s"
	}

	data "aci_l3out_ospf_interface_profile" "test" {
		logical_interface_profile_dn = "${aci_l3out_ospf_interface_profile.test.logical_interface_profile_dn}abc"
	}
	`, rName, rName, rName, rName, rName)
	return resource
}
