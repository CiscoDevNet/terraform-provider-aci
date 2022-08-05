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

func TestAccAciL3outStaticRoute_Basic(t *testing.T) {
	var l3out_static_route_default models.L3outStaticRoute
	var l3out_static_route_updated models.L3outStaticRoute
	resourceName := "aci_l3out_static_route.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("20.2.0.0/16")
	ipUpdated, _ := acctest.RandIpAddress("20.3.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL3outStaticRouteWithoutRequired(rName, fabDn2, ip, "ip"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL3outStaticRouteWithoutRequired(rName, fabDn2, ip, "fabric_node_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL3outStaticRouteConfig(rName, fabDn2, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists(resourceName, &l3out_static_route_default),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabDn2)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "no"),
					resource.TestCheckResourceAttr(resourceName, "pref", "1"),
					resource.TestCheckResourceAttr(resourceName, "rt_ctrl", "unspecified"),
				),
			},
			{
				Config: CreateAccL3outStaticRouteConfigWithOptionalValues(rName, fabDn2, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists(resourceName, &l3out_static_route_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabDn2)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l3out_static_route"),
					resource.TestCheckResourceAttr(resourceName, "aggregate", "yes"),
					resource.TestCheckResourceAttr(resourceName, "pref", "255"),
					resource.TestCheckResourceAttr(resourceName, "rt_ctrl", "bfd"),
					testAccCheckAciL3outStaticRouteIdEqual(&l3out_static_route_default, &l3out_static_route_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn2, ip, "pref", "125"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists(resourceName, &l3out_static_route_updated),
					resource.TestCheckResourceAttr(resourceName, "pref", "125"),
					testAccCheckAciL3outStaticRouteIdEqual(&l3out_static_route_default, &l3out_static_route_updated),
				),
			},
			{
				Config:      CreateAccL3outStaticRouteWithInavalidIP(rName, fabDn2, ip),
				ExpectError: regexp.MustCompile(`unknown property value (.)+`),
			},

			{
				Config:      CreateAccL3outStaticRouteRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config: CreateAccL3outStaticRouteConfigWithRequiredParams(rName, fabDn2, ip, ipUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists(resourceName, &l3out_static_route_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabDn2)),
					resource.TestCheckResourceAttr(resourceName, "ip", ipUpdated),
					testAccCheckAciL3outStaticRouteIdNotEqual(&l3out_static_route_default, &l3out_static_route_updated),
				),
			},
			{
				Config: CreateAccL3outStaticRouteConfig(rName, fabDn2, ip),
			},
			{
				Config: CreateAccL3outStaticRouteConfigWithRequiredParams(rNameUpdated, fabDn2, ip, ip),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outStaticRouteExists(resourceName, &l3out_static_route_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rNameUpdated, rNameUpdated, rNameUpdated, fabDn2)),
					resource.TestCheckResourceAttr(resourceName, "ip", ip),
					testAccCheckAciL3outStaticRouteIdNotEqual(&l3out_static_route_default, &l3out_static_route_updated),
				),
			},
		},
	})
}

func TestAccAciL3outStaticRoute_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("20.4.0.0/16")
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outStaticRouteConfig(rName, fabDn3, ip),
			},
			{
				Config:      CreateAccL3outStaticRouteWithInValidParentDn(rName, ip),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "aggregate", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "pref", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "pref", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "pref", "256"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, "rt_ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccL3outStaticRouteUpdatedAttr(rName, fabDn3, ip, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL3outStaticRouteConfig(rName, fabDn3, ip),
			},
		},
	})
}

func TestAccAciL3outStaticRoute_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	ip, _ := acctest.RandIpAddress("20.5.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL3outStaticRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL3outStaticRouteConfigMultiple(rName, fabDn4, ip[:len(ip)-1]),
			},
		},
	})
}

func testAccCheckAciL3outStaticRouteExists(name string, l3out_static_route *models.L3outStaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Static Route %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Static Route dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_static_routeFound := models.L3outStaticRouteFromContainer(cont)
		if l3out_static_routeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Static Route %s not found", rs.Primary.ID)
		}
		*l3out_static_route = *l3out_static_routeFound
		return nil
	}
}

func testAccCheckAciL3outStaticRouteDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l3out_static_route destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l3out_static_route" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_static_route := models.L3outStaticRouteFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Static Route %s Still exists", l3out_static_route.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL3outStaticRouteIdEqual(m1, m2 *models.L3outStaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l3out_static_route DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL3outStaticRouteIdNotEqual(m1, m2 *models.L3outStaticRoute) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l3out_static_route DNs are equal")
		}
		return nil
	}
}

func CreateL3outStaticRouteWithoutRequired(rName, tdn, ip, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route creation without ", attrName)
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
	`
	switch attrName {
	case "ip":
		rBlock += `
	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_logical_node_to_fabric_node.test.id
	#	ip  = "%s"
	}
		`
	case "fabric_node_dn":
		rBlock += `
	resource "aci_l3out_static_route" "test" {
	#	fabric_node_dn = aci_logical_node_to_fabric_node.test.id
		ip  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, ip, ip)
}

func CreateAccL3outStaticRouteConfigWithRequiredParams(rName, tdn, rtr, ip string) string {
	fmt.Printf("=== STEP  testing l3out_static_route creation with parent resource name %s,tdn %s, rtr_id %s and ip %s\n", rName, tdn, rtr, ip)
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
	`, rName, rName, rName, tdn, rtr, ip)
	return resource
}

func CreateAccL3outStaticRouteWithInavalidIP(rName, tdn, ip string) string {
	fmt.Println("=== STEP  testing l3out_static_route creation with invalid IP")
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
	`, rName, rName, rName, tdn, ip, rName)
	return resource
}

func CreateAccL3outStaticRouteWithInValidParentDn(rName, ip string) string {
	fmt.Println("=== STEP  testing l3out_static_route creation with invalid parent dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l3out_static_route" "test" {
		fabric_node_dn = aci_tenant.test.id
		ip  = "%s"
	}
	`, rName, ip)
	return resource
}

func CreateAccL3outStaticRouteConfig(rName, tdn, ip string) string {
	fmt.Println("=== STEP  testing l3out_static_route creation with required arguments only")
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
	`, rName, rName, rName, tdn, ip, ip)
	return resource
}

func CreateAccL3outStaticRouteConfigMultiple(rName, tdn, ip string) string {
	fmt.Println("=== STEP  testing multiple l3out_static_route creation with required arguments only")
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
		ip  = "%s${count.index}"
		count = 5
	}
	`, rName, rName, rName, tdn, ip, ip)
	return resource
}

func CreateAccL3outStaticRouteConfigWithOptionalValues(rName, tdn, ip string) string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route creation with optional parameters")
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
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_static_route"
		aggregate = "yes"
		pref = "255"
		rt_ctrl = "bfd"
	}
	`, rName, rName, rName, tdn, ip, ip)

	return resource
}

func CreateAccL3outStaticRouteRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l3out_static_route updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l3out_static_route" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_static_route"
		aggregate = "yes"
		pref = "2"
		rt_ctrl = "bfd"
	}
	`)

	return resource
}

func CreateAccL3outStaticRouteUpdatedAttr(rName, tdn, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing l3out_static_route attribute: %s = %s \n", attribute, value)
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
	`, rName, rName, rName, tdn, ip, ip, attribute, value)
	return resource
}
