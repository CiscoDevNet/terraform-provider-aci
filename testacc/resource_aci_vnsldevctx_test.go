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

func TestAccAciLogicalDeviceContext_Basic(t *testing.T) {
	var logical_device_context_default models.LogicalDeviceContext
	var logical_device_context_updated models.LogicalDeviceContext
	resourceName := "aci_logical_device_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	ctrctNameOrLbl := makeTestVariable(acctest.RandString(5))
	ctrctNameOrLblUpdated := makeTestVariable(acctest.RandString(5))

	graphNameOrLbl := makeTestVariable(acctest.RandString(5))
	graphNameOrLblUpdated := makeTestVariable(acctest.RandString(5))

	nodeNameOrLbl := makeTestVariable(acctest.RandString(5))
	nodeNameOrLblUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLogicalDeviceContextWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "tenant_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLogicalDeviceContextWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "ctrct_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateLogicalDeviceContextWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "graph_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},

			{
				Config:      CreateLogicalDeviceContextWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "node_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalDeviceContextConfig(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_default),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "context", ""),
				),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigWithOptionalValues(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", fvTenantName)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_logical_device_context"),

					resource.TestCheckResourceAttr(resourceName, "context", "ctx"),

					testAccCheckAciLogicalDeviceContextIdEqual(&logical_device_context_default, &logical_device_context_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config:      CreateAccLogicalDeviceContextRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigWithRequiredParams(rNameUpdated, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLbl),
					testAccCheckAciLogicalDeviceContextIdNotEqual(&logical_device_context_default, &logical_device_context_updated),
				),
			},
			{
				Config: CreateAccLogicalDeviceContextConfig(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigWithRequiredParams(rName, ctrctNameOrLblUpdated, graphNameOrLbl, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLblUpdated),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLbl),
					testAccCheckAciLogicalDeviceContextIdNotEqual(&logical_device_context_default, &logical_device_context_updated),
				),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigWithRequiredParams(rName, ctrctNameOrLbl, graphNameOrLblUpdated, nodeNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLblUpdated),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLbl),
					testAccCheckAciLogicalDeviceContextIdNotEqual(&logical_device_context_default, &logical_device_context_updated),
				),
			},
			{
				Config: CreateAccLogicalDeviceContextConfigWithRequiredParams(rName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLblUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists(resourceName, &logical_device_context_updated),
					resource.TestCheckResourceAttr(resourceName, "tenant_dn", fmt.Sprintf("uni/tn-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "ctrct_name_or_lbl", ctrctNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "graph_name_or_lbl", graphNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "node_name_or_lbl", nodeNameOrLblUpdated),
					testAccCheckAciLogicalDeviceContextIdNotEqual(&logical_device_context_default, &logical_device_context_updated),
				),
			},
		},
	})
}

func TestAccAciLogicalDeviceContext_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	ctrctNameOrLbl := makeTestVariable(acctest.RandString(5))
	graphNameOrLbl := makeTestVariable(acctest.RandString(5))
	nodeNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalDeviceContextConfig(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
			},
			{
				Config:      CreateAccLogicalDeviceContextWithInValidParentDn(rName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLogicalDeviceContextUpdatedAttr(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalDeviceContextUpdatedAttr(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalDeviceContextUpdatedAttr(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLogicalDeviceContextUpdatedAttr(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLogicalDeviceContextConfig(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
			},
		},
	})
}

func TestAccAciLogicalDeviceContext_MultipleCreateDelete(t *testing.T) {

	ctrctNameOrLbl := makeTestVariable(acctest.RandString(5))

	graphNameOrLbl := makeTestVariable(acctest.RandString(5))

	nodeNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalDeviceContextConfigMultiple(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl),
			},
		},
	})
}

func testAccCheckAciLogicalDeviceContextExists(name string, logical_device_context *models.LogicalDeviceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Device Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Device Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_device_contextFound := models.LogicalDeviceContextFromContainer(cont)
		if logical_device_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Device Context %s not found", rs.Primary.ID)
		}
		*logical_device_context = *logical_device_contextFound
		return nil
	}
}

func testAccCheckAciLogicalDeviceContextDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing logical_device_context destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_logical_device_context" {
			cont, err := client.Get(rs.Primary.ID)
			logical_device_context := models.LogicalDeviceContextFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Device Context %s Still exists", logical_device_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLogicalDeviceContextIdEqual(m1, m2 *models.LogicalDeviceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("logical_device_context DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLogicalDeviceContextIdNotEqual(m1, m2 *models.LogicalDeviceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("logical_device_context DNs are equal")
		}
		return nil
	}
}

func CreateLogicalDeviceContextWithoutRequired(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_device_context creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	`
	switch attrName {
	case "tenant_dn":
		rBlock += `
	resource "aci_logical_device_context" "test" {
	#	tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"	
		graph_name_or_lbl  = "%s"	
		node_name_or_lbl  = "%s"
	}
		`
	case "ctrct_name_or_lbl":
		rBlock += `
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
	#	ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}
		`
	case "graph_name_or_lbl":
		rBlock += `
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
	#	graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}
		`
	case "node_name_or_lbl":
		rBlock += `
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
	#	node_name_or_lbl  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
}

func CreateAccLogicalDeviceContextConfigWithRequiredParams(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_device_context creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateAccLogicalDeviceContextConfig(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_device_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateAccLogicalDeviceContextConfigMultiple(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  testing multiple logical_device_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s_${count.index}"
		graph_name_or_lbl  = "%s_${count.index}"
		node_name_or_lbl  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateAccLogicalDeviceContextWithInValidParentDn(rName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  Negative Case: testing logical_device_context creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_application_profile" "test"{
		tenant_dn  = aci_tenant.test.id
		name = "%s"
	}
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_application_profile.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"	
	}
	`, rName, rName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)
	return resource
}

func CreateAccLogicalDeviceContextConfigWithOptionalValues(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl string) string {
	fmt.Println("=== STEP  Basic: testing logical_device_context creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = "${aci_tenant.test.id}"
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_logical_device_context"
		context = "ctx"
		
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl)

	return resource
}

func CreateAccLogicalDeviceContextRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing logical_device_context updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_logical_device_context" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_logical_device_context"
		context = "ctx"
		
	}
	`)

	return resource
}

func CreateAccLogicalDeviceContextUpdatedAttr(fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, attribute, value string) string {
	fmt.Printf("=== STEP  testing logical_device_context attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn  = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl  = "%s"
		node_name_or_lbl  = "%s"
		%s = "%s"
	}
	`, fvTenantName, ctrctNameOrLbl, graphNameOrLbl, nodeNameOrLbl, attribute, value)
	return resource
}
