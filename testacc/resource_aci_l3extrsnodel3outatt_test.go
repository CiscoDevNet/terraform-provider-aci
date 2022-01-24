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

func TestAccAciFabricNode_Basic(t *testing.T) {
	var fabric_node_default models.FabricNode
	var fabric_node_updated models.FabricNode
	resourceName := "aci_logical_node_to_fabric_node.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	rtrid, _ := acctest.RandIpAddress("10.2.0.0/16")
	rtridOther, _ := acctest.RandIpAddress("10.3.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFabricNodeWithoutRequired(rName, rName, rName, fabDn2, "logical_node_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFabricNodeWithoutRequired(rName, rName, rName, fabDn2, "tdn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeConfig(rName, rName, rName, fabDn2, rtrid),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_default),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fabDn2),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "none"),
					resource.TestCheckResourceAttr(resourceName, "rtr_id", rtrid),
					resource.TestCheckResourceAttr(resourceName, "rtr_id_loop_back", "yes"),
				),
			},
			{
				Config: CreateAccFabricNodeConfigWithOptionalValues(rName, rName, rName, fabDn2, rtridOther),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fabDn2),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "anchor-node-mismatch"),
					resource.TestCheckResourceAttr(resourceName, "rtr_id", rtridOther),
					resource.TestCheckResourceAttr(resourceName, "rtr_id_loop_back", "no"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccFabricNodeRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "bd-profile-missmatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "bd-profile-missmatch"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "loopback-ip-missing"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "loopback-ip-missing"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "missing-mpls-infra-l3out"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "missing-mpls-infra-l3out"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "missing-rs-export-route-profile"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "missing-rs-export-route-profile"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "node-path-misconfig"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "node-path-misconfig"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "node-vlif-misconfig"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "node-vlif-misconfig"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "routerid-not-changable-with-mcast"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "routerid-not-changable-with-mcast"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn2, rtrid, "config_issues", "subnet-mismatch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "config_issues", "subnet-mismatch"),
					testAccCheckAciFabricNodeIdEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeConfigWithRequiredParams(rNameUpdated, fabDn2, rtrid),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rNameUpdated, rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fabDn2),
					testAccCheckAciFabricNodeIdNotEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
			{
				Config: CreateAccFabricNodeConfig(rName, rName, rName, fabDn2, rtrid),
			},
			{
				Config: CreateAccFabricNodeConfigWithRequiredParams(rName, fabDn3, rtrid),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFabricNodeExists(resourceName, &fabric_node_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_node_profile_dn", fmt.Sprintf("uni/tn-%s/out-%s/lnodep-%s", rName, rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "tdn", fabDn3),
					testAccCheckAciFabricNodeIdNotEqual(&fabric_node_default, &fabric_node_updated),
				),
			},
		},
	})
}

func TestAccAciFabricNode_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	rtrid, _ := acctest.RandIpAddress("10.5.0.0/16")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFabricNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricNodeConfig(rName, rName, rName, fabDn4, rtrid),
			},
			{
				Config:      CreateAccFabricNodeWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn4, rtrid, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn4, rtrid, "config_issues", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccFabricNodeWithInvalidRtr(rName, rName, rName, fabDn4, randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn4, rtrid, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccFabricNodeUpdatedAttr(rName, rName, rName, fabDn4, rtrid, "rtr_id_loop_back", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config: CreateAccFabricNodeConfig(rName, rName, rName, fabDn4, rtrid),
			},
		},
	})
}

func testAccCheckAciFabricNodeExists(name string, logical_node_to_fabric_node *models.FabricNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fabric Node %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fabric Node dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fabric_nodeFound := models.FabricNodeFromContainer(cont)
		if fabric_nodeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fabric Node %s not found", rs.Primary.ID)
		}
		*logical_node_to_fabric_node = *fabric_nodeFound
		return nil
	}
}

func testAccCheckAciFabricNodeDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_logical_node_to_fabric_node" {
			cont, err := client.Get(rs.Primary.ID)
			logical_node_to_fabric_node := models.FabricNodeFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fabric Node %s Still exists", logical_node_to_fabric_node.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFabricNodeIdEqual(m1, m2 *models.FabricNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("logical_node_to_fabric_node DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFabricNodeIdNotEqual(m1, m2 *models.FabricNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("logical_node_to_fabric_node DNs are equal")
		}
		return nil
	}
}

func CreateFabricNodeWithoutRequired(fvTenantName, l3extOutName, l3extLNodePName, tDn, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_node_to_fabric_node creation without ", attrName)
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
	
	`
	switch attrName {
	case "logical_node_profile_dn":
		rBlock += `
	resource "aci_logical_node_to_fabric_node" "test" {
	#	logical_node_profile_dn  = aci_logical_node_profile.test.id
		tdn  = "%s"
	}
		`
	case "tdn":
		rBlock += `
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_logical_node_profile.test.id
	#	tdn  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, l3extOutName, l3extLNodePName, tDn)
}

func CreateAccFabricNodeConfigWithRequiredParams(prName, tDn, ip string) string {
	fmt.Printf("=== STEP  testing logical_node_to_fabric_node creation with parent resource name %s and tdn %s\n", prName, tDn)
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
	`, prName, prName, prName, tDn, ip)
	return resource
}

func CreateAccFabricNodeConfig(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node creation with required arguments and rtr_id")
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
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip)
	return resource
}

func CreateAccFabricNodeWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing logical_node_to_fabric_node creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_logical_node_to_fabric_node" "test" {
		logical_node_profile_dn  = aci_tenant.test.id
		tdn  = aci_tenant.test.id
	}
	`, rName)
	return resource
}

func CreateAccFabricNodeConfigWithOptionalValues(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip string) string {
	fmt.Println("=== STEP  Basic: testing logical_node_to_fabric_node creation with optional parameters")
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
		logical_node_profile_dn  = "${aci_logical_node_profile.test.id}"
		tdn  = "%s"
		annotation = "orchestrator:terraform_testacc"
		config_issues = "anchor-node-mismatch"
		rtr_id = "%s"
		rtr_id_loop_back = "no"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip)

	return resource
}

func CreateAccFabricNodeRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing logical_node_to_fabric_node updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_logical_node_to_fabric_node" "test" {
		annotation = "orchestrator:terraform_testacc"
		config_issues = "anchor-node-mismatch"
		rtr_id = ""
		rtr_id_loop_back = "no"
		
	}
	`)

	return resource
}

func CreateAccFabricNodeWithInvalidRtr(fvTenantName, l3extOutName, l3extLNodePName, tDn, value string) string {
	fmt.Println("=== STEP  testing logical_node_to_fabric_node attribute: rtr_id =", value)
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
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, value)
	return resource
}

func CreateAccFabricNodeUpdatedAttr(fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, attribute, value string) string {
	fmt.Printf("=== STEP  testing logical_node_to_fabric_node attribute: %s = %s \n", attribute, value)
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
		%s = "%s"
	}
	`, fvTenantName, l3extOutName, l3extLNodePName, tDn, ip, attribute, value)
	return resource
}
