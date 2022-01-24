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

func TestAccAciVlanEncapsulationforVxlanTraffic_Basic(t *testing.T) {
	var vlan_encapsulationfor_vxlan_traffic_default models.VlanEncapsulationforVxlanTraffic
	var vlan_encapsulationfor_vxlan_traffic_updated models.VlanEncapsulationforVxlanTraffic
	resourceName := "aci_vlan_encapsulationfor_vxlan_traffic.test"
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVlanEncapsulationforVxlanTrafficWithoutRequired(infraAttEntityPName, "attachable_access_entity_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfig(infraAttEntityPName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists(resourceName, &vlan_encapsulationfor_vxlan_traffic_default),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", infraAttEntityPName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfigWithOptionalValues(infraAttEntityPName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists(resourceName, &vlan_encapsulationfor_vxlan_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", infraAttEntityPName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vlan_encapsulationfor_vxlan_traffic"),

					testAccCheckAciVlanEncapsulationforVxlanTrafficIdEqual(&vlan_encapsulationfor_vxlan_traffic_default, &vlan_encapsulationfor_vxlan_traffic_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfigWithRequiredParams(rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVlanEncapsulationforVxlanTrafficExists(resourceName, &vlan_encapsulationfor_vxlan_traffic_updated),
					resource.TestCheckResourceAttr(resourceName, "attachable_access_entity_profile_dn", fmt.Sprintf("uni/infra/attentp-%s", rNameUpdated)),
					testAccCheckAciVlanEncapsulationforVxlanTrafficIdNotEqual(&vlan_encapsulationfor_vxlan_traffic_default, &vlan_encapsulationfor_vxlan_traffic_updated),
				),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfig(infraAttEntityPName),
			},
		},
	})
}

func TestAccAciVlanEncapsulationforVxlanTraffic_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	infraAttEntityPName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfig(infraAttEntityPName),
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficUpdatedAttr(infraAttEntityPName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficUpdatedAttr(infraAttEntityPName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficUpdatedAttr(infraAttEntityPName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVlanEncapsulationforVxlanTrafficUpdatedAttr(infraAttEntityPName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVlanEncapsulationforVxlanTrafficConfig(infraAttEntityPName),
			},
		},
	})
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficExists(name string, vlan_encapsulationfor_vxlan_traffic *models.VlanEncapsulationforVxlanTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Vlan Encapsulation for Vxlan Traffic dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vlan_encapsulationfor_vxlan_trafficFound := models.VlanEncapsulationforVxlanTrafficFromContainer(cont)
		if vlan_encapsulationfor_vxlan_trafficFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s not found", rs.Primary.ID)
		}
		*vlan_encapsulationfor_vxlan_traffic = *vlan_encapsulationfor_vxlan_trafficFound
		return nil
	}
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vlan_encapsulationfor_vxlan_traffic" {
			cont, err := client.Get(rs.Primary.ID)
			vlan_encapsulationfor_vxlan_traffic := models.VlanEncapsulationforVxlanTrafficFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Vlan Encapsulation for Vxlan Traffic %s Still exists", vlan_encapsulationfor_vxlan_traffic.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficIdEqual(m1, m2 *models.VlanEncapsulationforVxlanTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vlan_encapsulationfor_vxlan_traffic DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVlanEncapsulationforVxlanTrafficIdNotEqual(m1, m2 *models.VlanEncapsulationforVxlanTraffic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vlan_encapsulationfor_vxlan_traffic DNs are equal")
		}
		return nil
	}
}

func CreateVlanEncapsulationforVxlanTrafficWithoutRequired(infraAttEntityPName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vlan_encapsulationfor_vxlan_traffic creation without ", attrName)
	rBlock := `
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "attachable_access_entity_profile_dn":
		rBlock += `
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
	#	attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	
	}
		`

	}
	return fmt.Sprintf(rBlock, infraAttEntityPName)
}

func CreateAccVlanEncapsulationforVxlanTrafficConfigWithRequiredParams(infraAttEntityPName string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic creation with Updated required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficConfig(infraAttEntityPName string) string {
	fmt.Println("=== STEP  testing vlan_encapsulationfor_vxlan_traffic creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
	}
	`, infraAttEntityPName)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vlan_encapsulationfor_vxlan_traffic creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_tenant.test.id	
	}
	`, rName)
	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficConfigWithOptionalValues(infraAttEntityPName string) string {
	fmt.Println("=== STEP  Basic: testing vlan_encapsulationfor_vxlan_traffic creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = "${aci_attachable_access_entity_profile.test.id}"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vlan_encapsulationfor_vxlan_traffic"
		
	}
	`, infraAttEntityPName)

	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vlan_encapsulationfor_vxlan_traffic updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vlan_encapsulationfor_vxlan_traffic"
		
	}
	`)

	return resource
}

func CreateAccVlanEncapsulationforVxlanTrafficUpdatedAttr(infraAttEntityPName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vlan_encapsulationfor_vxlan_traffic attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_attachable_access_entity_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vlan_encapsulationfor_vxlan_traffic" "test" {
		attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.test.id
		%s = "%s"
	}
	`, infraAttEntityPName, attribute, value)
	return resource
}
