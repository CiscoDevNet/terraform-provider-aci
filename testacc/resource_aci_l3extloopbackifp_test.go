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

var fabricNodeDn1 string = "topology/pod-1/node-201"
var fabricNodeDn2 string = "topology/pod-1/node-101"
var fabricNodeDn3 string = "topology/pod-1/node-111"
var fabricNodeDn4 string = "topology/pod-1/node-1"

func TestAccAcil3extLoopBackIfP_Basic(t *testing.T) {
	var aci_l3out_loopback_interface_profile_default models.LoopBackInterfaceProfile
	var aci_l3out_loopback_interface_profile_updated models.LoopBackInterfaceProfile
	resourceName := "aci_l3out_loopback_interface_profile.test"
	rName := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("2.2.0.0/16")
	addrOther, _ := acctest.RandIpAddress("3.3.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoopBackInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLoopBackInterfaceProfileWithoutRequired(rName, fabricNodeDn1, addr, addr, "fabric_node_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLoopBackInterfaceProfileWithoutRequired(rName, fabricNodeDn1, addr, addr, "addr"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfig(rName, fabricNodeDn1, addr, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists(resourceName, &aci_l3out_loopback_interface_profile_default),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabricNodeDn1)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfigWithOptionalValues(rName, fabricNodeDn1, addr, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists(resourceName, &aci_l3out_loopback_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabricNodeDn1)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					resource.TestCheckResourceAttr(resourceName, "annotation", "example"),
					resource.TestCheckResourceAttr(resourceName, "description", "from terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "example"),
					testAccCheckAciLoopBackInterfaceProfileIdEqual(&aci_l3out_loopback_interface_profile_default, &aci_l3out_loopback_interface_profile_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileWithInvalidIP(rName, fabricNodeDn1, addr),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfigWithRequiredParams(rName, fabricNodeDn2, addr, addr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists(resourceName, &aci_l3out_loopback_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabricNodeDn2)),
					resource.TestCheckResourceAttr(resourceName, "addr", addr),
					testAccCheckAciLoopBackInterfaceProfileIdNotEqual(&aci_l3out_loopback_interface_profile_default, &aci_l3out_loopback_interface_profile_updated),
				),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfig(rName, fabricNodeDn1, addr, addr),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfigWithRequiredParams(rName, fabricNodeDn1, addr, addrOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLoopBackInterfaceProfileExists(resourceName, &aci_l3out_loopback_interface_profile_updated),
					resource.TestCheckResourceAttr(resourceName, "fabric_node_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s/rsnodeL3OutAtt-[%s]", rName, rName, rName, fabricNodeDn1)),
					resource.TestCheckResourceAttr(resourceName, "addr", addrOther),
					testAccCheckAciLoopBackInterfaceProfileIdNotEqual(&aci_l3out_loopback_interface_profile_default, &aci_l3out_loopback_interface_profile_updated),
				),
			},
		},
	})
}

func TestAccAcil3extLoopBackIfP_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	addr, _ := acctest.RandIpAddress("4.4.0.0/16")
	longDescAnnotation := acctest.RandString(129)
	longNameAlias := acctest.RandString(64)
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLoopBackInterfaceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLoopBackInterfaceProfileConfig(rName, fabricNodeDn3, addr, addr),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileConfigWithInvalidParentDn(rName, addr),
				ExpectError: regexp.MustCompile(`unknown property value (.)+, name dn, class l3extLoopBackIfP (.)+`),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileConfigUpdatedAttr(rName, fabricNodeDn3, addr, addr, "annotation", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property annotation of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileConfigUpdatedAttr(rName, fabricNodeDn3, addr, addr, "description", longDescAnnotation),
				ExpectError: regexp.MustCompile(`property descr of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileConfigUpdatedAttr(rName, fabricNodeDn3, addr, addr, "name_alias", longNameAlias),
				ExpectError: regexp.MustCompile(`property nameAlias of (.)+ failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLoopBackInterfaceProfileConfigUpdatedAttr(rName, fabricNodeDn3, addr, addr, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLoopBackInterfaceProfileConfig(rName, fabricNodeDn3, addr, addr),
			},
		},
	})
}

func CreateAccLoopBackInterfaceProfileRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing LoopBackInterfaceProfile updation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_l3out_loopback_interface_profile" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l3out_route_tag_policy"
	}
	`)

	return resource
}

func CreateAccLoopBackInterfaceProfileConfigWithOptionalValues(rName, tdn, parent_addr, addr string) string {
	fmt.Println("=== STEP  Basic: testing l3out_route_tag_policy creation with optional parameters")
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
			rtr_id ="%s"
		}
		resource "aci_l3out_loopback_interface_profile" "test" {
			fabric_node_dn = aci_logical_node_to_fabric_node.test.id
			addr           = "%s"
			description    = "from terraform"
 			annotation     = "example"
 			name_alias     = "example"
		}
	`, rName, rName, rName, tdn, addr, addr)

	return resource
}

func CreateAccLoopBackInterfaceProfileWithInvalidIP(rName, tdn, parent_addr string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile creation with invalid ip", rName)
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
	`, rName, rName, rName, tdn, parent_addr, rName)
	return resource
}

func CreateAccLoopBackInterfaceProfileConfigWithRequiredParams(rName, tdn, parent_addr, addr string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile creation with required parameters")
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
	`, rName, rName, rName, tdn, parent_addr, addr)
	return resource
}

func CreateAccLoopBackInterfaceProfileConfigUpdatedAttr(rName, tdn, parent_addr, addr, key, value string) string {
	fmt.Printf("=== STEP  testing LoopBackInterfaceProfile updation for %s = %s\n", key, value)
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
			%s = "%s"
		}
	`, rName, rName, rName, tdn, addr, addr, key, value)
	return resource
}

func CreateAccLoopBackInterfaceProfileConfigWithInvalidParentDn(rName, addr string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile creation with invalid fabric_node_dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	resource "aci_l3out_loopback_interface_profile" "test" {
		fabric_node_dn = aci_tenant.test.id
		addr           = "%s"
	}
	`, rName, addr)
	return resource
}

func CreateAccLoopBackInterfaceProfileConfig(rName, tdn, parent_addr, addr string) string {
	fmt.Println("=== STEP  testing LoopBackInterfaceProfile creation with required arguments only")
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
	`, rName, rName, rName, tdn, addr, addr)
	return resource
}

func CreateLoopBackInterfaceProfileWithoutRequired(rName, tdn, parent_addr, addr, attrName string) string {
	fmt.Println("=== STEP  Basic: testing LoopBackInterfaceProfile creation without ", attrName)
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
	
	`
	switch attrName {
	case "fabric_node_dn":
		rBlock += `
	resource "aci_l3out_loopback_interface_profile" "test" {
		#fabric_node_dn  = aci_logical_node_to_fabric_node.test.id
	    addr  = "%s"
	}
		`
	case "addr":
		rBlock += `
	resource "aci_l3out_loopback_interface_profile" "test" {
		fabric_node_dn  = aci_logical_node_to_fabric_node.test.id
	#	addr  = "%s"
	}	`
	}

	return fmt.Sprintf(rBlock, rName, rName, rName, tdn, parent_addr, addr)
}

func testAccCheckAciLoopBackInterfaceProfileExists(name string, loop_back_interface_profile *models.LoopBackInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Loopback Interface Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Loopback Interface Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		loop_back_interface_profileFound := models.LoopBackInterfaceProfileFromContainer(cont)
		if loop_back_interface_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Loopback Interface Profile %s not found", rs.Primary.ID)
		}
		*loop_back_interface_profile = *loop_back_interface_profileFound
		return nil
	}
}

func testAccCheckAciLoopBackInterfaceProfileDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing L3out Loopback Interface Profile destroy")
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3out_loopback_interface_profile" {
			cont, err := client.Get(rs.Primary.ID)
			loop_back_interface_profile := models.LoopBackInterfaceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Loopback Interface Profile %s Still exists", loop_back_interface_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciLoopBackInterfaceProfileIdEqual(m1, m2 *models.LoopBackInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("LoopBackInterfaceProfile DNs are not equal")
		}
		return nil
	}
}
func testAccCheckAciLoopBackInterfaceProfileIdNotEqual(m1, m2 *models.LoopBackInterfaceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("LoopBackInterfaceProfile DNs are equal")
		}
		return nil
	}
}
