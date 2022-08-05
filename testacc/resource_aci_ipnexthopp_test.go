package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outStaticRouteNextHop_Basic(t *testing.T) {
	var l3out_static_route_next_hop_default models.L3outStaticRouteNextHop
	var l3out_static_route_next_hop_updated models.L3outStaticRouteNextHop
	resourceName := "aci_l3out_static_route_next_hop.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	rtrId, _ := acctest.RandIpAddress("20.2.0.0/16")
	nhAddr, _ := acctest.RandIpAddress("20.3.0.0/16")
	nhAddrUpdated, _ := acctest.RandIpAddress("20.4.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outStaticRouteNextHopWithoutRequired(rName, fabDn2, rtrId, nhAddr, "nh_addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outStaticRouteNextHopWithoutRequired(rName, fabDn2, rtrId, nhAddr, "static_route_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfig(rName, fabDn2, rtrId, nhAddr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_default),
					resource.TestCheckResourceAttr(resourceName, "static_route_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", rName, rName, rName, fabDn2, rtrId)),
					resource.TestCheckResourceAttr(resourceName, "nh_addr", nhAddr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "pref", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "nexthop_profile_type", "prefix"),
				),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfigWithOptionalValues(rName, fabDn2, rtrId, nhAddr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_updated),
					resource.TestCheckResourceAttr(resourceName, "static_route_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", rName, rName, rName, fabDn2, rtrId)),
					resource.TestCheckResourceAttr(resourceName, "nh_addr", nhAddr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_static_route_next_hop"),
					resource.TestCheckResourceAttr(resourceName, "pref", "1"),
					testAccCheckAciL3outStaticRouteNextHopIdEqual(&l3out_static_route_next_hop_default, &l3out_static_route_next_hop_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn2, rtrId, nhAddr, "pref", "128"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_updated),
					resource.TestCheckResourceAttr(resourceName, "pref", "128"),
					testAccCheckAciL3outStaticRouteNextHopIdEqual(&l3out_static_route_next_hop_default, &l3out_static_route_next_hop_updated),
				),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn2, rtrId, nhAddr, "pref", "255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_updated),
					resource.TestCheckResourceAttr(resourceName, "pref", "255"),
					testAccCheckAciL3outStaticRouteNextHopIdEqual(&l3out_static_route_next_hop_default, &l3out_static_route_next_hop_updated),
				),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopWithInvalidIP(rName, fabDn2, rtrId),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccL3outStaticRouteNextHopRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL3outStaticRouteNextHopConfigWithRequiredParams(rName, fabDn2, rtrId, nhAddrUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_updated),
					resource.TestCheckResourceAttr(resourceName, "static_route_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", rName, rName, rName, fabDn2, rtrId)),
					resource.TestCheckResourceAttr(resourceName, "nh_addr", nhAddrUpdated),
					testAccCheckAciL3outStaticRouteNextHopIdNotEqual(&l3out_static_route_next_hop_default, &l3out_static_route_next_hop_updated),
				),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfig(rName, fabDn2, rtrId, nhAddr),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfigWithRequiredParams(rNameUpdated, fabDn2, rtrId, nhAddrUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteNextHopExists(resourceName, &l3out_static_route_next_hop_updated),
					resource.TestCheckResourceAttr(resourceName, "static_route_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]/rt-[%s]", rNameUpdated, rNameUpdated, rNameUpdated, fabDn2, rtrId)),
					resource.TestCheckResourceAttr(resourceName, "nh_addr", nhAddrUpdated),
					testAccCheckAciL3outStaticRouteNextHopIdNotEqual(&l3out_static_route_next_hop_default, &l3out_static_route_next_hop_updated),
				),
			},
		},
	})
}

func TestAccAciL3outStaticRouteNextHop_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	rtrId, _ := acctest.RandIpAddress("20.5.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	nhAddr, _ := acctest.RandIpAddress("20.6.0.0/16")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outStaticRouteNextHopConfig(rName, fabDn3, rtrId, nhAddr),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopConfigInvalidParentDn(rName, nhAddr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "pref", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "pref", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "pref", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, "nexthop_profile_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, fabDn3, rtrId, nhAddr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outStaticRouteNextHopConfig(rName, fabDn3, rtrId, nhAddr),
			},
		},
	})
}

func TestAccAciL3outStaticRouteNextHop_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteNextHopDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outStaticRouteNextHopConfigMultiple(rName, fabDn4),
			},
		},
	})
}

func testAccCheckAciL3outStaticRouteNextHopExists(name string, l3out_static_route_next_hop *models.L3outStaticRouteNextHop) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Static Route Next Hop %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Static Route Next Hop dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_static_route_next_hopFound := models.L3outStaticRouteNextHopFromContainer(cont)
		if l3out_static_route_next_hopFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Static Route Next Hop %s not found", rs.Primary.ID)
		}
		*l3out_static_route_next_hop = *l3out_static_route_next_hopFound
		return nil
	}
}

func testAccCheckAciL3outStaticRouteNextHopDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_static_route_next_hop" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_static_route_next_hop := models.L3outStaticRouteNextHopFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Static Route Next Hop %s Still exists", l3out_static_route_next_hop.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outStaticRouteNextHopIdEqual(m1, m2 *models.L3outStaticRouteNextHop) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_static_route_next_hop DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outStaticRouteNextHopIdNotEqual(m1, m2 *models.L3outStaticRouteNextHop) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_static_route_next_hop DNs are equal")
		}
		return nil
	}
}

func CreateL3outStaticRouteNextHopWithoutRequired(rName, tdn, rtrId, nhAddr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route_next_hop creation without ", attrName)
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
	case "nh_addr":
		rBlock += `
	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
	#	nh_addr  = "%s"
	}
		`
	case "static_route_dn":
		rBlock += `
	resource "aci_l3out_static_route_next_hop" "test" {
	#	static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "%s"
	}
		`
	}

	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)
}

func CreateAccL3outStaticRouteNextHopConfigWithRequiredParams(rName, tdn, rtrId, nhAddr string) string {
	fmt.Printf("=== STEP  testing l3out_static_route_next_hop creation with parent resource name %s,tdn %s, rtr_id %s and nh_addr %s,\n", rName, tdn, rtrId, nhAddr)
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
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)
	return resource
}

func CreateAccL3outStaticRouteNextHopWithInvalidIP(rName, tdn, rtrId string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop creation with invalid IP")
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
	`, rName, rName, rName, tdn, rtrId, rtrId, rName)
	return resource
}

func CreateAccL3outStaticRouteNextHopConfigInvalidParentDn(rName, nhAddr string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop creation with invalid parent dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_tenant.test.id
		nh_addr  = "%s"
	}
	`, rName, nhAddr)
	return resource
}

func CreateAccL3outStaticRouteNextHopConfig(rName, tdn, rtrId, nhAddr string) string {
	fmt.Println("=== STEP  testing l3out_static_route_next_hop creation with required arguments only")
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
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)
	return resource
}

func CreateAccL3outStaticRouteNextHopConfigMultiple(rName, tdn string) string {
	fmt.Println("=== STEP  testing multiple l3out_static_route_next_hop creation with required arguments only")
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
		rtr_id = "20.7.0.0"
	}

	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "20.7.0.0"
	}

	resource "aci_l3out_static_route_next_hop" "test" {
		static_route_dn  = aci_l3out_static_route.test.id
		nh_addr  = "20.8.0.${count.index}"
		count = 5
	}
	`, rName, rName, rName, tdn)
	return resource
}

func CreateAccL3outStaticRouteNextHopConfigWithOptionalValues(rName, tdn, rtrId, nhAddr string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route_next_hop creation with optional parameters")
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
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_static_route_next_hop"
		pref = "1"
	}
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr)

	return resource
}

func CreateAccL3outStaticRouteNextHopRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route_next_hop updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_static_route_next_hop" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_static_route_next_hop"
		pref = "2"
		nexthop_profile_type = "none"
	}
	`)
	return resource
}

func CreateAccL3outStaticRouteNextHopUpdatedAttr(rName, tdn, rtrId, nhAddr, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_static_route_next_hop attribute: %s = %s \n", attribute, value)
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
	`, rName, rName, rName, tdn, rtrId, rtrId, nhAddr, attribute, value)
	return resource
}
