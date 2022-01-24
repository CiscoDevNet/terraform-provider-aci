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

func TestAccAciBgpBestPathPolicy_Basic(t *testing.T) {
	var bgp_best_path_policy_default models.BgpBestPathPolicy
	var bgp_best_path_policy_updated models.BgpBestPathPolicy
	resourceName := "aci_bgp_best_path_policy.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpBestPathPolicyWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpBestPathPolicyWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists(resourceName, &bgp_best_path_policy_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "0"),
				),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists(resourceName, &bgp_best_path_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bgp_best_path_policy"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "asPathMultipathRelax"),
					testAccCheckAciBgpBestPathPolicyIdEqual(&bgp_best_path_policy_default, &bgp_best_path_policy_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBgpBestPathPolicyConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccBgpBestPathPolicyRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists(resourceName, &bgp_best_path_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBgpBestPathPolicyIdNotEqual(&bgp_best_path_policy_default, &bgp_best_path_policy_updated),
				),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfig(rName, rName),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists(resourceName, &bgp_best_path_policy_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBgpBestPathPolicyIdNotEqual(&bgp_best_path_policy_default, &bgp_best_path_policy_updated),
				),
			},
		},
	})
}

func TestAccAciBgpBestPathPolicy_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpBestPathPolicyConfig(rName, rName),
			},
			{
				Config:      CreateAccBgpBestPathPolicyWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBgpBestPathPolicyUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpBestPathPolicyUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpBestPathPolicyUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpBestPathPolicyUpdatedAttr(rName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBgpBestPathPolicyUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBgpBestPathPolicyConfig(rName, rName),
			},
		},
	})
}

func TestAccAciBgpBestPathPolicy_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpBestPathPolicyConfigMultiple(rName),
			},
		},
	})
}

func testAccCheckAciBgpBestPathPolicyExists(name string, bgp_best_path_policy *models.BgpBestPathPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bgp Best Path Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bgp Best Path Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_best_path_policyFound := models.BgpBestPathPolicyFromContainer(cont)
		if bgp_best_path_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bgp Best Path Policy %s not found", rs.Primary.ID)
		}
		*bgp_best_path_policy = *bgp_best_path_policyFound
		return nil
	}
}

func testAccCheckAciBgpBestPathPolicyDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bgp_best_path_policy destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_best_path_policy" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_best_path_policy := models.BgpBestPathPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bgp Best Path Policy %s Still exists", bgp_best_path_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBgpBestPathPolicyIdEqual(m1, m2 *models.BgpBestPathPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_best_path_policy DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBgpBestPathPolicyIdNotEqual(m1, m2 *models.BgpBestPathPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_best_path_policy DNs are equal")
		}
		return nil
	}
}

func CreateBgpBestPathPolicyWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_best_path_policy creation without ", attrName)
	rBlock := `

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bgp_best_path_policy" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpBestPathPolicyConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing bgp_best_path_policy creation with tenant name %s and name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpBestPathPolicyConfigMultiple(rName string) string {
	fmt.Println("=== STEP  testing multiple bgp_best_path_policy creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, rName, rName)
	return resource
}

func CreateAccBgpBestPathPolicyConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_best_path_policy creation with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpBestPathPolicyWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bgp_best_path_policy creation with Invalid Parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccBgpBestPathPolicyConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_best_path_policy creation with optional parameters")
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_best_path_policy"
		ctrl = "asPathMultipathRelax"

	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBgpBestPathPolicyRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bgp_best_path_policy update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_bgp_best_path_policy" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_bgp_best_path_policy"
		ctrl = "asPathMultipathRelax"

	}
	`)

	return resource
}

func CreateAccBgpBestPathPolicyUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_best_path_policy attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`

	resource "aci_tenant" "test" {
		name 		= "%s"

	}

	resource "aci_bgp_best_path_policy" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
