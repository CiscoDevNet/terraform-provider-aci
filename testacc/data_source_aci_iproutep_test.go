package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciL3outStaticRouteDataSource_Basic(t *testing.T) {
	resourceName := "aci_l3out_static_route.test"
	dataSourceName := "data.aci_l3out_static_route.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("20.0.0.0/16")
	ipother, _ := acctest.RandIpAddress("20.1.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{

			{
				Config:      CreateL3outStaticRouteDSWithoutRequired(rName, fabDn1, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outStaticRouteDSWithoutRequired(rName, fabDn1, ip, "fabric_node_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outStaticRouteConfigDataSource(rName, fabDn1, ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "ip", resourceName, "ip"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_node_dn", resourceName, "fabric_node_dn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "aggregate", resourceName, "aggregate"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pref", resourceName, "pref"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rt_ctrl", resourceName, "rt_ctrl"),
				),
			},
			{
				Config:      CreateAccL3outStaticRouteDataSourceUpdate(rName, fabDn1, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccL3outStaticRouteDSWithInvalidIP(rName, fabDn1, ip, ipother),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccL3outStaticRouteDataSourceUpdatedResource(rName, fabDn1, ip, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccL3outStaticRouteConfigDataSource(rName, tdn, ip string) string {
	fmt.Println("=== STEP  testing l3out_static_route Data Source with required arguments only")
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

	data "aci_l3out_static_route" "test" {
		ip  = aci_l3out_static_route.test.ip
		fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		depends_on = [ aci_l3out_static_route.test ]
	}
	`, rName, rName, rName, tdn, ip, ip)
	return resource
}

func CreateL3outStaticRouteDSWithoutRequired(rName, tdn, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route Data Source without ", attrName)
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
	`
	switch attrName {
	case "ip":
		rBlock += `
	data "aci_l3out_static_route" "test" {
	#	ip  = aci_l3out_static_route.test.ip
		fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		depends_on = [ aci_l3out_static_route.test ]
	}
		`
	case "fabric_node_dn":
		rBlock += `
	data "aci_l3out_static_route" "test" {
		ip  = aci_l3out_static_route.test.ip
	#	fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		depends_on = [ aci_l3out_static_route.test ]
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, ip, ip)
}

func CreateAccL3outStaticRouteDSWithInvalidIP(rName, tdn, ip, ipother string) string {
	fmt.Println("=== STEP  testing l3out_static_route Data Source with invalid IP")
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

	data "aci_l3out_static_route" "test" {
		ip  = "%s"
		fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		depends_on = [ aci_l3out_static_route.test ]
	}
	`, rName, rName, rName, tdn, ip, ip, ipother)
	return resource
}

func CreateAccL3outStaticRouteDataSourceUpdate(rName, tdn, ip, key, value string) string {
	fmt.Println("=== STEP  testing l3out_static_route Data Source with random attribute")
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

	data "aci_l3out_static_route" "test" {
		ip  = aci_l3out_static_route.test.ip
		fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		%s = "%s"
		depends_on = [ aci_l3out_static_route.test ]
	}
	`, rName, rName, rName, tdn, ip, ip, key, value)
	return resource
}

func CreateAccL3outStaticRouteDataSourceUpdatedResource(rName, tdn, ip, key, value string) string {
	fmt.Println("=== STEP  testing l3out_static_route Data Source with updated resource")
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
		%s = "%s"
	}

	data "aci_l3out_static_route" "test" {
		ip  = aci_l3out_static_route.test.ip
		fabric_node_dn = aci_l3out_static_route.test.fabric_node_dn
		depends_on = [ aci_l3out_static_route.test ]
	}
	`, rName, rName, rName, tdn, ip, ip, key, value)
	return resource
}
