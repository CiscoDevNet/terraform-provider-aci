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

func TestAccAciHsrpInterfacePolicy_Basic(t *testing.T) {
	var hsrp_interface_policy_default models.HSRPInterfacePolicy
	var hsrp_interface_policy_updated models.HSRPInterfacePolicy
	resourceName := "aci_hsrp_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHsrpInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateHsrpInterfacePolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateHsrpInterfacePolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl", ""),
					resource.TestCheckResourceAttr(resourceName, "delay", "0"),
					resource.TestCheckResourceAttr(resourceName, "reload_delay", "0"),
				),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_hsrp_interface_policy"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "bfd"),
					resource.TestCheckResourceAttr(resourceName, "delay", "1"),
					resource.TestCheckResourceAttr(resourceName, "reload_delay", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccHsrpInterfacePolicyConfigInvalidName(acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)* failed validation`),
			},

			{
				Config: CreateAccHsrpInterfacePolicyConfigUpdatedRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciHsrpInterfacePolicyIdNotEqual(&hsrp_interface_policy_default, &hsrp_interface_policy_updated),
				),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfigUpdatedRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciHsrpInterfacePolicyIdNotEqual(&hsrp_interface_policy_default, &hsrp_interface_policy_updated),
				),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyConfigUpdateWithoutRequiredArguments(rName, "description", acctest.RandString(5)),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciHsrpInterfacePolicy_Update(t *testing.T) {
	var hsrp_interface_policy_default models.HSRPInterfacePolicy
	var hsrp_interface_policy_updated models.HSRPInterfacePolicy
	resourceName := "aci_hsrp_interface_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHsrpInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_default),
				),
			},

			{

				Config: CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "ctrl", "bia"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "bia"),
					testAccCheckAciHsrpInterfacePolicyIdEqual(&hsrp_interface_policy_default, &hsrp_interface_policy_updated),
				),
			},
			{

				Config: CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "ctrl", "bfd,bia"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "bfd,bia"),
					testAccCheckAciHsrpInterfacePolicyIdEqual(&hsrp_interface_policy_default, &hsrp_interface_policy_updated),
				),
			},
			{
				Config: CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "ctrl", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHsrpInterfacePolicyExists(resourceName, &hsrp_interface_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "ctrl", ""),
					testAccCheckAciHsrpInterfacePolicyIdEqual(&hsrp_interface_policy_default, &hsrp_interface_policy_updated),
				),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciHsrpInterfacePolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHsrpInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyWithInvalidParentDn(rName, rName),
				ExpectError: regexp.MustCompile(`configured object (.)+ not found (.)+,`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected(.*)to be one of(.)*, got(.)*`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "ctrl", "bfd,bfd"),
				ExpectError: regexp.MustCompile(`unexpected duplicate values in ctrl`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "delay", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, "reload_delay", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccHsrpInterfacePolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named(.)*is not expected here.`),
			},
			{
				Config: CreateAccHsrpInterfacePolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccHsrpInterfacePolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHsrpInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAciHsrpInterfaceProfilesConfig(rName),
			},
		},
	})
}

func testAccCheckAciHsrpInterfacePolicyExists(name string, hsrp_interface_policy *models.HSRPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Hsrp Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Hsrp Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		hsrp_interface_policyFound := models.HSRPInterfacePolicyFromContainer(cont)
		if hsrp_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Hsrp Interface Policy %s not found", rs.Primary.ID)
		}
		*hsrp_interface_policy = *hsrp_interface_policyFound
		return nil
	}
}

func testAccCheckAciHsrpInterfacePolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing hsrp_interface_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_hsrp_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			hsrp_interface_policy := models.HSRPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Hsrp Interface Policy %s Still exists", hsrp_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciHsrpInterfacePolicyIdEqual(m1, m2 *models.HSRPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("hsrp_interface_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciHsrpInterfacePolicyIdNotEqual(m1, m2 *models.HSRPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("hsrp_interface_policy DNs are equal")
		}
		return nil
	}
}

func CreateAccAciHsrpInterfaceProfilesConfig(rName string) string {
	fmt.Println("=== STEP  creating multiple Hrsp Interface Profiles")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
		resource "aci_hsrp_interface_policy" "test1" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	
		resource "aci_hsrp_interface_policy" "test2" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}

	resource "aci_hsrp_interface_policy" "test3" {
		tenant_dn = aci_tenant.test.id
		name = "%s"
	}
	`, rName, rName, rName+"1", rName+"2", rName+"3")
	return resource

}
func CreateHsrpInterfacePolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing hsrp_interface_policy creation without Requierd Parameters ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"	
	}
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_hsrp_interface_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
		description = "created while acceptance testing"
	}
		`
	case "name":
		rBlock += `
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
		description = "created while acceptance testing"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccHsrpInterfacePolicyConfigUpdatedRequiredParams(rName, rName2 string) string {
	fmt.Println("=== STEP  testing hsrp_interface_policy updation using required arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName2)
	return resource
}

func CreateAccHsrpInterfacePolicyConfigUpdateWithoutRequiredArguments(rName, attribute, value string) string {
	fmt.Println("=== STEP  testing hsrp_interface_policy updation without required arguments")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}

	resource "aci_hsrp_interface_policy" "test" {
		%s = "%s"
	}
	`, rName, attribute, value)
	return resource
}

func CreateAccHsrpInterfacePolicyConfigInvalidName(rName string) string {
	fmt.Println("=== STEP  testing hsrp_interface_policy creation with Invalid Name")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing hsrp_interface_policy creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyWithInvalidParentDn(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Negative Case: testing hsrp_interface_policy creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}invalid"
		name  = "%s"	
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccHsrpInterfacePolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing hsrp_interface_policy creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_hsrp_interface_policy"
		ctrl = "bfd"
		delay = "1"
		reload_delay = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccHsrpInterfacePolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing hsrp_interface_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
	}
	
	resource "aci_hsrp_interface_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
