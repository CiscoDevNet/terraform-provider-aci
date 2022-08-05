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

func TestAccAciVRFSnmpContextCommunity_Basic(t *testing.T) {
	var vrf_snmp_context_community_default models.SNMPCommunity
	var vrf_snmp_context_community_updated models.SNMPCommunity
	resourceName := "aci_vrf_snmp_context_community.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateVRFSnmpContextCommunityWithoutRequired(fvTenantName, fvCtxName, rName, "vrf_snmp_context_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateVRFSnmpContextCommunityWithoutRequired(fvTenantName, fvCtxName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfig(fvTenantName, fvCtxName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextCommunityExists(resourceName, &vrf_snmp_context_community_default),
					resource.TestCheckResourceAttr(resourceName, "vrf_snmp_context_dn", fmt.Sprintf("uni/tn-%s/ctx-%s/snmpctx", fvTenantName, fvCtxName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
				),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfigWithOptionalValues(fvTenantName, fvCtxName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextCommunityExists(resourceName, &vrf_snmp_context_community_updated),
					resource.TestCheckResourceAttr(resourceName, "vrf_snmp_context_dn", fmt.Sprintf("uni/tn-%s/ctx-%s/snmpctx", fvTenantName, fvCtxName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_vrf_snmp_context_community"),

					testAccCheckAciVRFSnmpContextCommunityIdEqual(&vrf_snmp_context_community_default, &vrf_snmp_context_community_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityConfigUpdatedName(fvTenantName, fvCtxName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccVRFSnmpContextCommunityRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextCommunityExists(resourceName, &vrf_snmp_context_community_updated),
					resource.TestCheckResourceAttr(resourceName, "vrf_snmp_context_dn", fmt.Sprintf("uni/tn-%s/ctx-%s/snmpctx", rNameUpdated, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciVRFSnmpContextCommunityIdNotEqual(&vrf_snmp_context_community_default, &vrf_snmp_context_community_updated),
				),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfig(fvTenantName, fvCtxName, rName),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVRFSnmpContextCommunityExists(resourceName, &vrf_snmp_context_community_updated),
					resource.TestCheckResourceAttr(resourceName, "vrf_snmp_context_dn", fmt.Sprintf("uni/tn-%s/ctx-%s/snmpctx", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciVRFSnmpContextCommunityIdNotEqual(&vrf_snmp_context_community_default, &vrf_snmp_context_community_updated),
				),
			},
		},
	})
}

func TestAccAciVRFSnmpContextCommunity_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFSnmpContextCommunityConfig(fvTenantName, fvCtxName, rName),
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityUpdatedAttr(fvTenantName, fvCtxName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityUpdatedAttr(fvTenantName, fvCtxName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccVRFSnmpContextCommunityUpdatedAttr(fvTenantName, fvCtxName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccVRFSnmpContextCommunityUpdatedAttr(fvTenantName, fvCtxName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccVRFSnmpContextCommunityConfig(fvTenantName, fvCtxName, rName),
			},
		},
	})
}

func TestAccAciVRFSnmpContextCommunity_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	fvCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciVRFSnmpContextCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccVRFSnmpContextCommunityConfigMultiple(fvTenantName, fvCtxName, rName),
			},
		},
	})
}

func testAccCheckAciVRFSnmpContextCommunityExists(name string, vrf_snmp_context_community *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VRF Snmp Context Community %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VRF Snmp Context Community dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vrf_snmp_context_communityFound := models.SNMPCommunityFromContainer(cont)
		if vrf_snmp_context_communityFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VRF Snmp Context Community %s not found", rs.Primary.ID)
		}
		*vrf_snmp_context_community = *vrf_snmp_context_communityFound
		return nil
	}
}

func testAccCheckAciVRFSnmpContextCommunityDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing vrf_snmp_context_community destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vrf_snmp_context_community" {
			cont, err := client.Get(rs.Primary.ID)
			vrf_snmp_context_community := models.SNMPCommunityFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VRF Snmp Context Community %s Still exists", vrf_snmp_context_community.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVRFSnmpContextCommunityIdEqual(m1, m2 *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("vrf_snmp_context_community DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciVRFSnmpContextCommunityIdNotEqual(m1, m2 *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("vrf_snmp_context_community DNs are equal")
		}
		return nil
	}
}

func CreateVRFSnmpContextCommunityWithoutRequired(fvTenantName, fvCtxName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context_community creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	`
	switch attrName {
	case "vrf_snmp_context_dn":
		rBlock += `
	resource "aci_vrf_snmp_context_community" "test" {
	#	vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, fvCtxName, rName)
}

func CreateAccVRFSnmpContextCommunityConfigWithRequiredParams(prName, rName string) string {
	fmt.Printf("=== STEP  testing vrf_snmp_context_community creation with parent resource name %s and resource name %s\n", prName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "%s"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s"
	}
	`, prName, prName, prName, rName)
	return resource
}
func CreateAccVRFSnmpContextCommunityConfigUpdatedName(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s"
	}
	`, fvTenantName, fvCtxName, rName)
	return resource
}

func CreateAccVRFSnmpContextCommunityConfig(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  testing vrf_snmp_context_community creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s"
	}
	`, fvTenantName, fvCtxName, rName)
	return resource
}

func CreateAccVRFSnmpContextCommunityConfigMultiple(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  testing multiple vrf_snmp_context_community creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, fvCtxName, rName)
	return resource
}

func CreateAccVRFSnmpContextCommunityWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing vrf_snmp_context_community creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccVRFSnmpContextCommunityConfigWithOptionalValues(fvTenantName, fvCtxName, rName string) string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context_community creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = "${aci_vrf_snmp_context.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vrf_snmp_context_community"
		
	}
	`, fvTenantName, fvCtxName, rName)

	return resource
}

func CreateAccVRFSnmpContextCommunityRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing vrf_snmp_context_community updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_vrf_snmp_context_community" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_vrf_snmp_context_community"
		
	}
	`)

	return resource
}

func CreateAccVRFSnmpContextCommunityUpdatedAttr(fvTenantName, fvCtxName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing vrf_snmp_context_community attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_vrf" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}

	resource "aci_vrf_snmp_context" "test" {
		vrf_dn = aci_vrf.test.id
		name = "example"
	}
	
	resource "aci_vrf_snmp_context_community" "test" {
		vrf_snmp_context_dn  = aci_vrf_snmp_context.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, fvCtxName, rName, attribute, value)
	return resource
}
