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

func TestAccAciBgpRouteSummarization_Basic(t *testing.T) {
	var bgp_route_summarization_default models.BgpRouteSummarization
	var bgp_route_summarization_updated models.BgpRouteSummarization
	resourceName := "aci_bgp_route_summarization.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateBgpRouteSummarizationWithoutRequired(rName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateBgpRouteSummarizationWithoutRequired(rName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfig(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists(resourceName, &bgp_route_summarization_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "attrmap", ""),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "none"),
				),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfigWithOptionalValues(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists(resourceName, &bgp_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_bgp_route_summarization"),
					resource.TestCheckResourceAttr(resourceName, "attrmap", "attrmap_test"),
					resource.TestCheckResourceAttr(resourceName, "ctrl", "as-set"),
					testAccCheckAciBgpRouteSummarizationIdEqual(&bgp_route_summarization_default, &bgp_route_summarization_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:      CreateAccBgpRouteSummarizationConfigUpdatedName(rName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccBgpRouteSummarizationRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists(resourceName, &bgp_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciBgpRouteSummarizationIdNotEqual(&bgp_route_summarization_default, &bgp_route_summarization_updated),
				),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfig(rName, rName),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists(resourceName, &bgp_route_summarization_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciBgpRouteSummarizationIdNotEqual(&bgp_route_summarization_default, &bgp_route_summarization_updated),
				),
			},
		},
	})
}

func TestAccAciBgpRouteSummarization_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpRouteSummarizationConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccBgpRouteSummarizationWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, "attrmap", acctest.RandString(513)),
				ExpectError: regexp.MustCompile(`property (.)+ failed validation for value ''`),
			},
			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, "ctrl", randomValue),
				ExpectError: regexp.MustCompile(`expected (.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccBgpRouteSummarizationConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciBgpRouteSummarization_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccBgpRouteSummarizationConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciBgpRouteSummarizationExists(name string, bgp_route_summarization *models.BgpRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Bgp Route Summarization %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bgp Route Summarization dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_route_summarizationFound := models.BgpRouteSummarizationFromContainer(cont)
		if bgp_route_summarizationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Bgp Route Summarization %s not found", rs.Primary.ID)
		}
		*bgp_route_summarization = *bgp_route_summarizationFound
		return nil
	}
}

func testAccCheckAciBgpRouteSummarizationDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing bgp_route_summarization destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bgp_route_summarization" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_route_summarization := models.BgpRouteSummarizationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Bgp Route Summarization %s Still exists", bgp_route_summarization.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBgpRouteSummarizationIdEqual(m1, m2 *models.BgpRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("bgp_route_summarization DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciBgpRouteSummarizationIdNotEqual(m1, m2 *models.BgpRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("bgp_route_summarization DNs are equal")
		}
		return nil
	}
}

func CreateBgpRouteSummarizationWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_route_summarization creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_bgp_route_summarization" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccBgpRouteSummarizationConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccBgpRouteSummarizationConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteSummarizationConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing bgp_route_summarization creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteSummarizationConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple bgp_route_summarization creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccBgpRouteSummarizationWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing bgp_route_summarization creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		name = "%s"
		tenant_dn = aci_tenant.test.id
	}
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccBgpRouteSummarizationConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing bgp_route_summarization update without required parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_bgp_route_summarization"
		attrmap = "attrmap_test"
		ctrl = "as-set"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccBgpRouteSummarizationRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing bgp_route_summarization update without required parameters")
	resource := fmt.Sprintln(`
	resource "aci_bgp_route_summarization" "test" {
		description = "created while acceptance testing"
		annotation = "tag"
		name_alias = "test_bgp_route_summarization"
		attrmap = ""
		ctrl = "as-set"
	}
	`)

	return resource
}

func CreateAccBgpRouteSummarizationUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing bgp_route_summarization attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_bgp_route_summarization" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
