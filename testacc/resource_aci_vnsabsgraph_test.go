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

func TestAccAciL4L7ServiceGraphTemplate_Basic(t *testing.T) {
	var l4_l7_service_graph_template_default models.L4L7ServiceGraphTemplate
	var l4_l7_service_graph_template_updated models.L4L7ServiceGraphTemplate
	resourceName := "aci_l4_l7_service_graph_template.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateL4L7ServiceGraphTemplateWithoutRequired(fvTenantName, rName, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateL4L7ServiceGraphTemplateWithoutRequired(fvTenantName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_type", "legacy"),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "UNSPECIFIED"),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfigWithOptionalValues(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_l4_l7_service_graph_template"),

					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_type", "cloud"),

					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "ONE_NODE_ADC_ONE_ARM"),

					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"term_cons_dn",
					"term_cons_name",
					"term_node_cons_dn",
					"term_node_prov_dn",
					"term_prov_dn",
					"term_prov_name",
				},
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateConfigUpdatedName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccL4L7ServiceGraphTemplateRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfigWithRequiredParams(rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciL4L7ServiceGraphTemplateIdNotEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfigWithRequiredParams(rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciL4L7ServiceGraphTemplateIdNotEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
		},
	})
}

func TestAccAciL4L7ServiceGraphTemplate_Update(t *testing.T) {
	var l4_l7_service_graph_template_default models.L4L7ServiceGraphTemplate
	var l4_l7_service_graph_template_updated models.L4L7ServiceGraphTemplate
	resourceName := "aci_l4_l7_service_graph_template.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_default),
				),
			},

			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "ONE_NODE_ADC_ONE_ARM_L3EXT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "ONE_NODE_ADC_ONE_ARM_L3EXT"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "ONE_NODE_ADC_TWO_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "ONE_NODE_ADC_TWO_ARM"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "ONE_NODE_FW_ROUTED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "ONE_NODE_FW_ROUTED"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "ONE_NODE_FW_TRANS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "ONE_NODE_FW_TRANS"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM_L3EXT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_ONE_ARM_L3EXT"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_TWO_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_ROUTED_ADC_TWO_ARM"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_ONE_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_ONE_ARM"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_ONE_ARM_L3EXT"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_ONE_ARM_L3EXT"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_TWO_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists(resourceName, &l4_l7_service_graph_template_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_template_type", "TWO_NODE_FW_TRANS_ADC_TWO_ARM"),
					testAccCheckAciL4L7ServiceGraphTemplateIdEqual(&l4_l7_service_graph_template_default, &l4_l7_service_graph_template_updated),
				),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciL4L7ServiceGraphTemplate_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "l4_l7_service_graph_template_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, "ui_template_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName),
			},
		},
	})
}

func TestAccAciL4L7ServiceGraphTemplate_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccL4L7ServiceGraphTemplateConfigMultiple(fvTenantName, rName),
			},
		},
	})
}

func testAccCheckAciL4L7ServiceGraphTemplateExists(name string, l4_l7_service_graph_template *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L4 L7 Service Graph Template %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L4 L7 Service Graph Template dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l4_l7_service_graph_templateFound := models.L4L7ServiceGraphTemplateFromContainer(cont)
		if l4_l7_service_graph_templateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L4 L7 Service Graph Template %s not found", rs.Primary.ID)
		}
		*l4_l7_service_graph_template = *l4_l7_service_graph_templateFound
		return nil
	}
}

func testAccCheckAciL4L7ServiceGraphTemplateDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l4_l7_service_graph_template" {
			cont, err := client.Get(rs.Primary.ID)
			l4_l7_service_graph_template := models.L4L7ServiceGraphTemplateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L4 L7 Service Graph Template %s Still exists", l4_l7_service_graph_template.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL4L7ServiceGraphTemplateIdEqual(m1, m2 *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("l4_l7_service_graph_template DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciL4L7ServiceGraphTemplateIdNotEqual(m1, m2 *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("l4_l7_service_graph_template DNs are equal")
		}
		return nil
	}
}

func CreateL4L7ServiceGraphTemplateWithoutRequired(fvTenantName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing l4_l7_service_graph_template creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_l4_l7_service_graph_template" "test" {
	#	tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, rName)
}

func CreateAccL4L7ServiceGraphTemplateConfigWithRequiredParams(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}
func CreateAccL4L7ServiceGraphTemplateConfigUpdatedName(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateConfig(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing l4_l7_service_graph_template creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateConfigMultiple(fvTenantName, rName string) string {
	fmt.Println("=== STEP  testing multiple l4_l7_service_graph_template creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, rName)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing l4_l7_service_graph_template creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_application_profile.test.id
		name  = "%s"	
	}
	`, rName, rName, rName)
	return resource
}

func CreateAccL4L7ServiceGraphTemplateConfigWithOptionalValues(fvTenantName, rName string) string {
	fmt.Println("=== STEP  Basic: testing l4_l7_service_graph_template creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l4_l7_service_graph_template"
		l4_l7_service_graph_template_type = "cloud"
		ui_template_type = "ONE_NODE_ADC_ONE_ARM"
		
	}
	`, fvTenantName, rName)

	return resource
}

func CreateAccL4L7ServiceGraphTemplateRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing l4_l7_service_graph_template updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_l4_l7_service_graph_template" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_l4_l7_service_graph_template"
		l4_l7_service_graph_template_type = "cloud"
		ui_template_type = "ONE_NODE_ADC_ONE_ARM"
		
	}
	`)

	return resource
}

func CreateAccL4L7ServiceGraphTemplateUpdatedAttr(fvTenantName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing l4_l7_service_graph_template attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = aci_tenant.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, rName, attribute, value)
	return resource
}
