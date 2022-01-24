package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciVlanEncapsulationforVxlanTrafficDataSource_Basic(t *testing.T) {
	resourceName := "aci_vlan_encapsulationfor_vxlan_traffic.test"
	dataSourceName := "data.aci_vlan_encapsulationfor_vxlan_traffic.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVlanEncapsulationforVxlanTrafficDSWithoutRequired(rName, "attachable_access_entity_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfigDataSource(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "attachable_access_entity_profile_dn", resourceName, "attachable_access_entity_profile_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficDataSourceUpdate(rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficDSWithInvalidParentDn(rName),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},

			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficDataSourceUpdatedResource(rName, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccVlanEncapsulationforVxlanTrafficConfigDataSource(infraAttEntityPName string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}

	data "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		depends_on = [
			aci_vlan_encapsulationfor_vxlan_traffic.test
		]
	}
	`, infraAttEntityPName)
	return resource
}

func CreateVlanEncapsulationforVxlanTrafficDSWithoutRequired(infraAttEntityPName, attr string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic Data Source without required argument")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}

	data "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		# attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		depends_on = [
			aci_vlan_encapsulationfor_vxlan_traffic.test
		]
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficDSWithInvalidParentDn(infraAttEntityPName string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic Data Source with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}

	data "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = "${aci_attachable_access_entity_profile.test.id}_invalid"
		depends_on = [
			aci_vlan_encapsulationfor_vxlan_traffic.test
		]
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficDataSourceUpdate(infraAttEntityPName, key, value string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}

	data "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		%s = "%s"
		depends_on = [
			aci_vlan_encapsulationfor_vxlan_traffic.test
		]
	}
	`, infraAttEntityPName, key, value)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficDataSourceUpdatedResource(infraAttEntityPName, key, value string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		%s = "%s"
	}

	data "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		depends_on = [
			aci_vlan_encapsulationfor_vxlan_traffic.test
		]
	}
	`, infraAttEntityPName, key, value)
	return resource
}
