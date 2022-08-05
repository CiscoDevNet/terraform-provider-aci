package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outStaticRouteNextHopDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_static_route_next_hop.test"
	dataSourceName := "data.aci_l3out_static_route_next_hop.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	rtrId, _ := acctest.RandIpAddress("20.0.0.0/16")
	nhAddr, _ := acctest.RandIpAddress("20.1.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outStaticRouteNextHopDSWithoutRequired(rName, fabDn1, rtrId, nhAddr, "nh_addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outStaticRouteNextHopDSWithoutRequired(rName, fabDn1, rtrId, nhAddr, "static_route_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfigDataSource(rName, fabDn1, rtrId, nhAddr),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "static_route_dn", resourceName, "static_route_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "nh_addr", resourceName, "nh_addr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref", resourceName, "pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "nexthop_profile_type", resourceName, "nexthop_profile_type"),
				),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopDataSourceUpdate(rName, fabDn1, rtrId, nhAddr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outStaticRouteNextHopDSWithInvalidIP(rName, fabDn1, rtrId, nhAddr),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopDataSourceUpdatedResource(rName, fabDn1, rtrId, nhAddr, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outStaticRouteNextHopConfigDataSource(rName, tdn, rtrId, nhAddr string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop Data Source with required arguments only")
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

	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
	}

	data "aci_l3out_static_route_next_hop" "test" {
		static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
		nh_addr  = aci_l3out_static_route_next_hop.test.nh_addr
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)
	return resource
}

func CreateL3outStaticRouteNextHopDSWithoutRequired(rName, tdn, rtrId, nhAddr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route_next_hop Data Source without ", attrName)
	rBlock := `
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

	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
	}
	`
	switch attrName {
	case "nh_addr":
		rBlock += `
	data "aci_l3out_static_route_next_hop" "test" {
		static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
	#	nh_addr  = aci_l3out_static_route_next_hop.test.nh_addr
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
		`
	case "static_route_dn":
		rBlock += `
	data "aci_l3out_static_route_next_hop" "test" {
	#	static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
		nh_addr  = aci_l3out_static_route_next_hop.test.nh_addr
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)
}

func CreateAccL3outStaticRouteNextHopDSWithInvalidIP(rName, tdn, rtrId, nhAddr string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop Data Source with required arguments only")
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

	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
	}

	data "aci_l3out_static_route_next_hop" "test" {
		static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
		nh_addr  = "%s"
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr, rtrId)
	return resource
}

func CreateAccL3outStaticRouteNextHopDataSourceUpdate(rName, tdn, rtrId, nhAddr, key, value string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop Data Source with random attribute")
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

	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
	}

	data "aci_l3out_static_route_next_hop" "test" {
		static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
		nh_addr  = aci_l3out_static_route_next_hop.test.nh_addr
		%s = "%s"
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr, key, value)
	return resource
}

func CreateAccL3outStaticRouteNextHopDataSourceUpdatedResource(rName, tdn, rtrId, nhAddr, key, value string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop Data Source with updated resource")
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

	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
		rtr_id = "%s"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
		%s = "%s"
	}

	data "aci_l3out_static_route_next_hop" "test" {
		static_route_dn = aci_l3out_static_route_next_hop.test.static_route_dn
		nh_addr  = aci_l3out_static_route_next_hop.test.nh_addr
		depends_on = [ aci_l3out_static_route_next_hop.test ]
	}
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr, key, value)
	return resource
}
