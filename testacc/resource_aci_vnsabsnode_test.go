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

func TestAccAciFunctionNode_Basic(t *testing.T) {
	var function_node_default models.FunctionNode
	var function_node_updated models.FunctionNode
	resourceName := "aci_function_node.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateFunctionNodeWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "l4_l7_service_graph_template_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateFunctionNodeWithoutRequired(fvTenantName, vnsAbsGraphName, rName, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_default),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", fvTenantName, vnsAbsGraphName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "OTHER"),
					resource.TestCheckResourceAttr(resourceName, "func_type", "GoTo"),
					resource.TestCheckResourceAttr(resourceName, "is_copy", "no"),
					resource.TestCheckResourceAttr(resourceName, "managed", "yes"),
					resource.TestCheckResourceAttr(resourceName, "routing_mode", "unspecified"),
					resource.TestCheckResourceAttr(resourceName, "sequence_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "share_encap", "no"),
				),
			},
			{
				Config: CreateAccFunctionNodeConfigWithOptionalValues(rName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, vnsAbsGraphName)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_function_node"),

					resource.TestCheckResourceAttr(resourceName, "func_template_type", "ADC_ONE_ARM"),

					resource.TestCheckResourceAttr(resourceName, "func_type", "GoThrough"),

					resource.TestCheckResourceAttr(resourceName, "is_copy", "yes"),

					resource.TestCheckResourceAttr(resourceName, "managed", "no"),

					resource.TestCheckResourceAttr(resourceName, "routing_mode", "Redirect"),

					resource.TestCheckResourceAttr(resourceName, "share_encap", "yes"),
					// need to recreate the resource to update "is_copy" attribute.
					testAccCheckAciFunctionNodeIdNotEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"conn_consumer_dn",
					"conn_provider_dn",
					"l4_l7_service_graph_template_dn",
				},
			},
			{
				Config:      CreateAccFunctionNodeConfigUpdatedName(fvTenantName, vnsAbsGraphName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`property name of (.)+ failed validation`),
			},

			{
				Config:      CreateAccFunctionNodeRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccFunctionNodeConfigWithRequiredParams(rName, rNameUpdated, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAciFunctionNodeIdNotEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
			},
			{
				Config: CreateAccFunctionNodeConfigWithRequiredParams(rName, rName, rNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "l4_l7_service_graph_template_dn", fmt.Sprintf("uni/tn-%s/AbsGraph-%s", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdated),
					testAccCheckAciFunctionNodeIdNotEqual(&function_node_default, &function_node_updated),
				),
			},
		},
	})
}

func TestAccAciFunctionNode_Update(t *testing.T) {
	var function_node_default models.FunctionNode
	var function_node_updated models.FunctionNode
	resourceName := "aci_function_node.test"
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_default),
				),
			},

			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "ADC_TWO_ARM"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "ADC_TWO_ARM"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "CLOUD_NATIVE_FW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "CLOUD_NATIVE_FW"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "CLOUD_NATIVE_LB"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "CLOUD_NATIVE_LB"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "CLOUD_VENDOR_FW"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "CLOUD_VENDOR_FW"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "CLOUD_VENDOR_LB"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "CLOUD_VENDOR_LB"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "FW_ROUTED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "FW_ROUTED"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", "FW_TRANS"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_template_type", "FW_TRANS"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_type", "L1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_type", "L1"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_type", "L2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_type", "L2"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_type", "None"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists(resourceName, &function_node_updated),
					resource.TestCheckResourceAttr(resourceName, "func_type", "None"),
					testAccCheckAciFunctionNodeIdEqual(&function_node_default, &function_node_updated),
				),
			},
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
			},
		},
	})
}

func TestAccAciFunctionNode_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
			},
			{
				Config:      CreateAccFunctionNodeWithInValidParentDn(rName),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_template_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "func_type", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "is_copy", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "managed", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "routing_mode", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "sequence_number", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, "share_encap", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName),
			},
		},
	})
}

func TestAccAciFunctionNode_MultipleCreateDelete(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsAbsGraphName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFunctionNodeConfigMultiple(fvTenantName, vnsAbsGraphName, rName),
			},
		},
	})
}

func testAccCheckAciFunctionNodeExists(name string, function_node *models.FunctionNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Function Node %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Function Node dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		function_nodeFound := models.FunctionNodeFromContainer(cont)
		if function_nodeFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Function Node %s not found", rs.Primary.ID)
		}
		*function_node = *function_nodeFound
		return nil
	}
}

func testAccCheckAciFunctionNodeDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing function_node destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_function_node" {
			cont, err := client.Get(rs.Primary.ID)
			function_node := models.FunctionNodeFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Function Node %s Still exists", function_node.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciFunctionNodeIdEqual(m1, m2 *models.FunctionNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("function_node DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciFunctionNodeIdNotEqual(m1, m2 *models.FunctionNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("function_node DNs are equal")
		}
		return nil
	}
}

func CreateFunctionNodeWithoutRequired(fvTenantName, vnsAbsGraphName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing function_node creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "l4_l7_service_graph_template_dn":
		rBlock += `
	resource "aci_function_node" "test" {
	#	l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsAbsGraphName, rName)
}

func CreateAccFunctionNodeConfigWithRequiredParams(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing function_node creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}
func CreateAccFunctionNodeConfigUpdatedName(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing function_node creation with invalid name = ", rName)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccFunctionNodeConfig(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing function_node creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccFunctionNodeConfigMultiple(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  testing multiple function_node creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, vnsAbsGraphName, rName)
	return resource
}

func CreateAccFunctionNodeWithInValidParentDn(rName string) string {
	fmt.Println("=== STEP  Negative Case: testing function_node creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_tenant.test.id
		name  = "%s"	
	}
	`, rName, rName)
	return resource
}

func CreateAccFunctionNodeConfigWithOptionalValues(fvTenantName, vnsAbsGraphName, rName string) string {
	fmt.Println("=== STEP  Basic: testing function_node creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = "${aci_l4_l7_service_graph_template.test.id}"
		name  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_function_node"
		func_template_type = "ADC_ONE_ARM"
		func_type = "GoThrough"
		is_copy = "yes"
		managed = "no"
		routing_mode = "Redirect"
		share_encap = "yes"
		
	}
	`, fvTenantName, vnsAbsGraphName, rName)

	return resource
}

func CreateAccFunctionNodeRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing function_node updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_function_node" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_function_node"
		func_template_type = "ADC_ONE_ARM"
		func_type = "GoThrough"
		is_copy = "yes"
		managed = "no"
		routing_mode = "Redirect"
		share_encap = "yes"
		
	}
	`)

	return resource
}

func CreateAccFunctionNodeUpdatedAttr(fvTenantName, vnsAbsGraphName, rName, attribute, value string) string {
	fmt.Printf("=== STEP  testing function_node attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_l4_l7_service_graph_template" "test" {
		name 		= "%s"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_function_node" "test" {
		l4_l7_service_graph_template_dn  = aci_l4_l7_service_graph_template.test.id
		name  = "%s"
		%s = "%s"
	}
	`, fvTenantName, vnsAbsGraphName, rName, attribute, value)
	return resource
}
