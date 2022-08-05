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

func TestAccAciVRFSnmpContext_Basic(t *testing.T) {
	var vrf_snmp_context_default models.SNMPContextProfile
	var vrf_snmp_context_updated models.SNMPContextProfile
	resourceName := "aci_vrf_snmp_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVRFSnmpContextWithoutRequired(fvTenantName, fvCtxName, "vrf_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextConfig(fvTenantName, fvCtxName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextExists(resourceName, &vrf_snmp_context_default),
					resource.TestCheckResourceAttr(resourceName, "vrf_dn", fmt.Sprintf("uni/tn-%s/ctx-%s", fvTenantName, fvCtxName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
				),
			},
			{
				Config: CreateAccVRFSnmpContextConfigWithOptionalValues(fvTenantName, fvCtxName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextExists(resourceName, &vrf_snmp_context_updated),
					resource.TestCheckResourceAttr(resourceName, "vrf_dn", fmt.Sprintf("uni/tn-%s/ctx-%s", fvTenantName, fvCtxName)),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vrf_snmp_context"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_vrf_snmp_context"),

					testAccCheckAciVRFSnmpContextIdEqual(&vrf_snmp_context_default, &vrf_snmp_context_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVRFSnmpContextRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextExists(resourceName, &vrf_snmp_context_updated),
					resource.TestCheckResourceAttr(resourceName, "vrf_dn", fmt.Sprintf("uni/tn-%s/ctx-%s", rName, rNameUpdated)),
					testAccCheckAciVRFSnmpContextIdNotEqual(&vrf_snmp_context_default, &vrf_snmp_context_updated),
				),
			},
			{
				Config: CreateAccVRFSnmpContextConfig(fvTenantName, fvCtxName),
			},
		},
	})
}

func TestAccAciVRFSnmpContext_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFSnmpContextConfig(fvTenantName, fvCtxName),
			},
			{
				Config:      CreateAccVRFSnmpContextWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVRFSnmpContextUpdatedAttr(fvTenantName, fvCtxName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFSnmpContextUpdatedAttr(fvTenantName, fvCtxName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVRFSnmpContextUpdatedAttr(fvTenantName, fvCtxName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVRFSnmpContextConfig(fvTenantName, fvCtxName),
			},
		},
	})
}

func testAccCheckAciVRFSnmpContextExists(name string, vrf_snmp_context *models.SNMPContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VRF Snmp Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VRF Snmp Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vrf_snmp_contextFound := models.SNMPContextProfileFromContainer(cont)
		if vrf_snmp_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VRF Snmp Context %s not found", rs.Primary.ID)
		}
		*vrf_snmp_context = *vrf_snmp_contextFound
		return nil
	}
}

func testAccCheckAciVRFSnmpContextDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vrf_snmp_context destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vrf_snmp_context" {
			cont, err := client.Get(rs.Primary.ID)
			vrf_snmp_context := models.SNMPContextProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VRF Snmp Context %s Still exists", vrf_snmp_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVRFSnmpContextIdEqual(m1, m2 *models.SNMPContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vrf_snmp_context DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVRFSnmpContextIdNotEqual(m1, m2 *models.SNMPContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vrf_snmp_context DNs are equal")
		}
		return nil
	}
}

func CreateVRFSnmpContextWithoutRequired(fvTenantName, fvCtxName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "vrf_dn":
		rBlock += `
	resource "aci_vrf_snmp_context" "test" {
	#	vrf_dn  = aci_vrf.test.id
	}
		`

	}
	return fmt.Sprintf(rBlock, fvTenantName, fvCtxName)
}

func CreateAccVRFSnmpContextConfigWithRequiredParams(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}
	`, fvTenantName, fvCtxName)
	return resource
}
func CreateAccVRFSnmpContextConfig(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
	}
	`, fvTenantName, fvCtxName)
	return resource
}

func CreateAccVRFSnmpContextConfigMultiple(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  testing multiple vrf_snmp_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		count = 5
	}
	`, fvTenantName, fvCtxName)
	return resource
}

func CreateAccVRFSnmpContextWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vrf_snmp_context creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_tenant.test.id	
	}
	`, rName)
	return resource
}

func CreateAccVRFSnmpContextConfigWithOptionalValues(fvTenantName, fvCtxName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = "${aci_vrf.test.id}"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vrf_snmp_context"
		name = "test_vrf_snmp_context"
		
	}
	`, fvTenantName, fvCtxName)

	return resource
}

func CreateAccVRFSnmpContextRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vrf_snmp_context" "test" {
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vrf_snmp_context"
		
	}
	`)

	return resource
}

func CreateAccVRFSnmpContextUpdatedAttr(fvTenantName, fvCtxName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vrf_snmp_context attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_vrf_snmp_context" "test" {
		vrf_dn  = aci_vrf.test.id
		%s = "%s"
	}
	`, fvTenantName, fvCtxName, attribute, value)
	return resource
}
