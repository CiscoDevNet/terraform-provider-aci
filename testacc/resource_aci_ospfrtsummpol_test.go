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

func TestAccAciOspfRouteSummarization_Basic(t *testing.T) {
	var ospf_route_summarization_default models.OspfRouteSummarization
	var ospf_route_summarization_updated models.OspfRouteSummarization
	resourceName := "aci_ospf_route_summarization.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateOspfRouteSummarizationWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateOspfRouteSummarizationWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "cost", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "inter_area_enabled", "no"),
					resource.TestCheckResourceAttr(resourceName, "tag", "0"),
				),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_ospf_route_summarization"),
					resource.TestCheckResourceAttr(resourceName, "cost", "1"),
					resource.TestCheckResourceAttr(resourceName, "inter_area_enabled", "yes"),
					resource.TestCheckResourceAttr(resourceName, "tag", "1"),
					testAccCheckAciOspfRouteSummarizationIdEqual(&ospf_route_summarization_default, &ospf_route_summarization_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccOspfRouteSummarizationConfigWithRequiredParams(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccOspfRouteSummarizationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciOspfRouteSummarizationIdNotEqual(&ospf_route_summarization_default, &ospf_route_summarization_updated),
				),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfig(rName, rName),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciOspfRouteSummarizationIdNotEqual(&ospf_route_summarization_default, &ospf_route_summarization_updated),
				),
			},
		},
	})
}

func TestAccAciOspfRouteSummarization_Update(t *testing.T) {
	var ospf_route_summarization_default models.OspfRouteSummarization
	var ospf_route_summarization_updated models.OspfRouteSummarization
	resourceName := "aci_ospf_route_summarization.test"
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOspfRouteSummarizationConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_default),
				),
			},
			{
				Config: CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "cost", "16777215"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "cost", "16777215"),
					testAccCheckAciOspfRouteSummarizationIdEqual(&ospf_route_summarization_default, &ospf_route_summarization_updated),
				),
			},
			{
				Config: CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "cost", "8000000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists(resourceName, &ospf_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "cost", "8000000"),
					testAccCheckAciOspfRouteSummarizationIdEqual(&ospf_route_summarization_default, &ospf_route_summarization_updated),
				),
			},
		},
	})
}

func TestAccAciOspfRouteSummarization_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOspfRouteSummarizationConfig(rName, rName),
			},
			{
				Config:      CreateAccOspfRouteSummarizationWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "cost", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "cost", "16777216"),
				ExpectError: regexp.MustCompile(`Property (.)+ is out of range`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "inter_area_enabled", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, "tag", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccOspfRouteSummarizationUpdatedAttr(rName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccOspfRouteSummarizationConfig(rName, rName),
			},
		},
	})
}

func TestAccAciOspfRouteSummarization_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccOspfRouteSummarizationConfigs(rName),
			},
		},
	})
}

func testAccCheckAciOspfRouteSummarizationExists(name string, ospf_route_summarization *models.OspfRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Ospf Route Summarization %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ospf Route Summarization dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_route_summarizationFound := models.OspfRouteSummarizationFromContainer(cont)
		if ospf_route_summarizationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Ospf Route Summarization %s not found", rs.Primary.ID)
		}
		*ospf_route_summarization = *ospf_route_summarizationFound
		return nil
	}
}

func testAccCheckAciOspfRouteSummarizationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing ospf_route_summarization destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ospf_route_summarization" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_route_summarization := models.OspfRouteSummarizationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Ospf Route Summarization %s Still exists", ospf_route_summarization.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciOspfRouteSummarizationIdEqual(m1, m2 *models.OspfRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("ospf_route_summarization DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciOspfRouteSummarizationIdNotEqual(m1, m2 *models.OspfRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("ospf_route_summarization DNs are equal")
		}
		return nil
	}
}

func CreateOspfRouteSummarizationWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_route_summarization Data Source without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_ospf_route_summarization" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccOspfRouteSummarizationConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing ospf_route_summarization creation with parent resource name %s and resource name %s", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOspfRouteSummarizationConfigs(rName string) string {
	fmt.Println("=== STEP  testing multiple ospf_route_summarization creation")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_ospf_route_summarization" "test1" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}

	resource "aci_ospf_route_summarization" "test2" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, rName, rName, rName+"1", rName+"2")
	return resource
}

func CreateAccOspfRouteSummarizationConfig(fvTenantName, rName string) string {
	fmt.Printf("=== STEP  testing ospf_route_summarization creation with parent resource name %s and name %s\n", fvTenantName, rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccOspfRouteSummarizationWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing ospf_route_summarization creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn  = aci_tenant.test.id
	}
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccOspfRouteSummarizationConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing ospf_route_summarization creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_ospf_route_summarization"
		cost = "1"
		inter_area_enabled = "yes"
		tag = "1"
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccOspfRouteSummarizationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing ospf_route_summarization updadation without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_ospf_route_summarization" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_ospf_route_summarization"
		cost = "1"
		inter_area_enabled = "yes"
		tag = ""
	}
	`)

	return resource
}

func CreateAccOspfRouteSummarizationUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing ospf_route_summarization attribute: %s=%s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_ospf_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
