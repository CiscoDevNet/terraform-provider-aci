package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAcil3extLoopBackIfPDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_loopback_interface_profile.test"
	dataSourceName := "data.aci_l3out_loopback_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("5.5.0.0/16")
	addrother, _ := acctest.RandIpAddress("6.6.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoopBackInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLoopBackInterfaceProfileDSWithoutRequired(rName, fabricNodeDn4, addr, addr, "fabric_node_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLoopBackInterfaceProfileDSWithoutRequired(rName, fabricNodeDn4, addr, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileDSConfig(rName, fabricNodeDn4, addr, addr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "fabric_node_dn", dataSourceName, "fabric_node_dn"),
					resource.TestCheckResourceAttrPair(resourceName, "addr", dataSourceName, "addr"),
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
					resource.TestCheckResourceAttrPair(resourceName, "annotation", dataSourceName, "annotation"),
					resource.TestCheckResourceAttrPair(resourceName, "name_alias", dataSourceName, "name_alias"),
				),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileDSConfigRandomAttr(rName, fabricNodeDn4, addr, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			}, {
				Config:      CreateAccLoopBackInterfaceProfileDSConfigWithInvalidIP(rName, fabricNodeDn4, addr, addr, addrother),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileDSConfigUpdatedResource(rName, fabricNodeDn4, addr, addr, "description", randomValue),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "description", dataSourceName, "description"),
				),
			},
		},
	})
}

func CreateAccLoopBackInterfaceProfileDSConfigUpdatedResource(rName, tdn, parent_addr, addr, key, value string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		}
		resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name = "%s"
		}
		resource "aci_logical_node_profile" "test" {
			l3_outside_dn = aci_l3_outside.test.id
			name ="%s"
		}
		resource "aci_logical_node_to_fabric_node" "test" {
			logical_node_profile_dn  = aci_logical_node_profile.test.id
			tdn               = "%s"
			rtr_id            = "%s"
		}
		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_logical_node_to_fabric_node.test.id
			addr           = "%s"
			%s             = "%s"
		}
		data "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_l3out_loopback_interface_profile.test.fabric_node_dn
			addr           = aci_l3out_loopback_interface_profile.test.addr
		}
	`, rName, rName, rName, tdn, addr, addr, key, value)
	return resource
}

func CreateAccLoopBackInterfaceProfileDSConfigWithInvalidIP(rName, tdn, parent_addr, addr, addrother string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile Data Source with invalid ip")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		}
		resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name = "%s"
		}
		resource "aci_logical_node_profile" "test" {
			l3_outside_dn = aci_l3_outside.test.id
			name ="%s"
		}
		resource "aci_logical_node_to_fabric_node" "test" {
			logical_node_profile_dn  = aci_logical_node_profile.test.id
			tdn               = "%s"
			rtr_id            = "%s"
		}
		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_logical_node_to_fabric_node.test.id
			addr           = "%s"
		}
		data "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_l3out_loopback_interface_profile.test.fabric_node_dn
			addr           = "%s"
		}
	`, rName, rName, rName, tdn, parent_addr, addr, addrother)
	return resource
}

func CreateAccLoopBackInterfaceProfileDSConfigRandomAttr(rName, tdn, parent_addr, addr, key, value string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		}
		resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name = "%s"
		}
		resource "aci_logical_node_profile" "test" {
			l3_outside_dn = aci_l3_outside.test.id
			name ="%s"
		}
		resource "aci_logical_node_to_fabric_node" "test" {
			logical_node_profile_dn  = aci_logical_node_profile.test.id
			tdn               = "%s"
			rtr_id            = "%s"
		}
		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_logical_node_to_fabric_node.test.id
			addr           = "%s"
		}
		data "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_l3out_loopback_interface_profile.test.fabric_node_dn
			addr           = aci_l3out_loopback_interface_profile.test.addr
			%s             = "%s"
		}
	`, rName, rName, rName, tdn, addr, addr, key, value)
	return resource
}

func CreateAccLoopBackInterfaceProfileDSConfig(rName, tdn, parent_addr, addr string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		}
		resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name = "%s"
		}
		resource "aci_logical_node_profile" "test" {
			l3_outside_dn = aci_l3_outside.test.id
			name ="%s"
		}
		resource "aci_logical_node_to_fabric_node" "test" {
			logical_node_profile_dn  = aci_logical_node_profile.test.id
			tdn               = "%s"
			rtr_id            = "%s"
		}
		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_logical_node_to_fabric_node.test.id
			addr           = "%s"
		}
		data "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_l3out_loopback_interface_profile.test.fabric_node_dn
			addr           = aci_l3out_loopback_interface_profile.test.addr
		}
	`, rName, rName, rName, tdn, addr, addr)
	return resource
}

func CreateLoopBackInterfaceProfileDSWithoutRequired(rName, tdn, parent_addr, addr, attrName string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		}
		resource "aci_l3_outside" "test" {
			tenant_dn      = aci_tenant.test.id
			name = "%s"
		}
		resource "aci_logical_node_profile" "test" {
			l3_outside_dn = aci_l3_outside.test.id
			name ="%s"
		}
		resource "aci_logical_node_to_fabric_node" "test" {
			logical_node_profile_dn  = aci_logical_node_profile.test.id
			tdn               = "%s"
			rtr_id  ="%s"
		}

		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn  = aci_logical_node_to_fabric_node.test.id
			addr  = "%s"
		}
	
	`
	switch attrName {
	case "fabric_node_dn":
		rBlock += `
	data "aci_l3out_loopback_interface_profile" "test" {
		#fabric_node_dn  = aci_l3out_loopback_interface_profile.test.fabric_node_dn
	    addr  = aci_l3out_loopback_interface_profile.test.addr
	}
		`
	case "addr":
		rBlock += `
	data "aci_l3out_loopback_interface_profile" "test" {
		fabric_node_dn  = aci_l3out_loopback_interface_profile.test.fabric_node_dn
	#	addr  = aci_l3out_loopback_interface_profile.test.addr
	}	`
	}

	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, parent_addr, addr)
}
