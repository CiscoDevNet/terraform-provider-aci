package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outBfdInterfaceProfileDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_bfd_interface_profile.test"
	dataSourceName := "data.aci_l3out_bfd_interface_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outBfdInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outBfdInterfaceProfileDataSourceWithoutRequired(rName),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			}, {
				Config: CreateAccL3outBfdInterfaceProfileConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "logical_interface_profile_dn", resourceName, "logical_interface_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "key_id", resourceName, "key_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "interface_profile_type", resourceName, "interface_profile_type"),
				),
			},
			{
				Config:      CreateAccL3outBfdInterfaceProfileDataSourceUpdateRandomAttr(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outBfdInterfaceProfileDataSourceWithInvalidParentDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outBfdInterfaceProfileDataSourceUpdate(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateL3outBfdInterfaceProfileDataSourceWithoutRequired(rName string) string {
	fmt.Println("=== STEP  Basic: Testing l3out_bfd_interface_profile data soruce creation without required attribute")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_bfd_interface_profile" "test" {
		depends_on = [aci_l3out_bfd_interface_profile.test]
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outBfdInterfaceProfileConfigDataSource(rName string) string {
	fmt.Println("=== STEP  Testing l3out_bfd_interface_profile data source creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		depends_on = [
			aci_l3out_bfd_interface_profile.test
		]
	}
	`, rName, rName, rName, rName)
	return resource
}
func CreateAccL3outBfdInterfaceProfileDataSourceWithInvalidParentDn(rName string) string {
	fmt.Println("=== STEP  testing l3out_bfd_interface_profile creation with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = "${aci_logical_interface_profile.test.id}xyz"
		depends_on = [
			aci_l3out_bfd_interface_profile.test
		]
	}
	`, rName, rName, rName, rName)
	return resource
}

func CreateAccL3outBfdInterfaceProfileDataSourceUpdateRandomAttr(rName, key, value string) string {
	fmt.Println("=== STEP  Testing l3out_bfd_interface_profile data source update with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
	}

	data "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
		depends_on = [
			aci_l3out_bfd_interface_profile.test
		]
	}
	`, rName, rName, rName, rName, key, value)
	return resource
}

func CreateAccL3outBfdInterfaceProfileDataSourceUpdate(rName, key, value string) string {
	fmt.Printf("=== STEP  Testing l3out_bfd_interface_profile data source update with %s =  %s\n", key, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_l3_outside" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_node_profile" "test" {
		name 		= "%s"
		l3_outside_dn = aci_l3_outside.test.id
	}
	
	resource "aci_logical_interface_profile" "test" {
		name 		= "%s"
		logical_node_profile_dn = aci_logical_node_profile.test.id
	}
	
	resource "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		%s = "%s"
	}

	data "aci_l3out_bfd_interface_profile" "test" {
		logical_interface_profile_dn  = aci_logical_interface_profile.test.id
		depends_on = [
			aci_l3out_bfd_interface_profile.test
		]
	}
	`, rName, rName, rName, rName, key, value)
	return resource
}
