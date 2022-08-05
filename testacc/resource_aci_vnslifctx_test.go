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

func TestAccAciLogicalInterfaceContext_Basic(t *testing.T) {
	var logical_interface_context_default models.LogicalInterfaceContext
	var logical_interface_context_updated models.LogicalInterfaceContext
	resourceName := "aci_logical_interface_context.test"
	rName := makeTestVariable(acctest.RandString(5))
	rNameUpdated := makeTestVariable(acctest.RandString(5))

	connNameOrLbl := makeTestVariable(acctest.RandString(5))
	connNameOrLblUpdated := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsLDevCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config:      CreateLogicalInterfaceContextWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, "logical_device_context_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateLogicalInterfaceContextWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, "conn_name_or_lbl"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfig(fvTenantName, vnsLDevCtxName, connNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists(resourceName, &logical_interface_context_default),
					resource.TestCheckResourceAttr(resourceName, "logical_device_context_dn", fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-any-n-any", fvTenantName, vnsLDevCtxName)),
					resource.TestCheckResourceAttr(resourceName, "conn_name_or_lbl", connNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name_alias", ""),
					resource.TestCheckResourceAttr(resourceName, "l3_dest", "yes"),
					resource.TestCheckResourceAttr(resourceName, "permit_log", "no"),
				),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfigWithOptionalValues(fvTenantName, vnsLDevCtxName, connNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists(resourceName, &logical_interface_context_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_device_context_dn", fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-any-n-any", fvTenantName, vnsLDevCtxName)),
					resource.TestCheckResourceAttr(resourceName, "conn_name_or_lbl", connNameOrLbl),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_logical_interface_context"),

					resource.TestCheckResourceAttr(resourceName, "l3_dest", "no"),

					resource.TestCheckResourceAttr(resourceName, "permit_log", "yes"),

					testAccCheckAciLogicalInterfaceContextIdEqual(&logical_interface_context_default, &logical_interface_context_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"logical_device_context_dn"},
			},
			{
				Config:      CreateAccLogicalInterfaceContextConfigWithInvalidName(fvTenantName, acctest.RandString(65)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceContextRemovingRequiredField(),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfigWithRequiredParams(rName, rNameUpdated, connNameOrLbl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists(resourceName, &logical_interface_context_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_device_context_dn", fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-any-n-any", rName, rNameUpdated)),
					resource.TestCheckResourceAttr(resourceName, "conn_name_or_lbl", connNameOrLbl),
					testAccCheckAciLogicalInterfaceContextIdNotEqual(&logical_interface_context_default, &logical_interface_context_updated),
				),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfig(fvTenantName, vnsLDevCtxName, connNameOrLbl),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfigWithRequiredParams(rName, rName, connNameOrLblUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists(resourceName, &logical_interface_context_updated),
					resource.TestCheckResourceAttr(resourceName, "logical_device_context_dn", fmt.Sprintf("uni/tn-%s/ldevCtx-c-%s-g-any-n-any", rName, rName)),
					resource.TestCheckResourceAttr(resourceName, "conn_name_or_lbl", connNameOrLblUpdated),
					testAccCheckAciLogicalInterfaceContextIdNotEqual(&logical_interface_context_default, &logical_interface_context_updated),
				),
			},
		},
	})
}

func TestAccAciLogicalInterfaceContext_Negative(t *testing.T) {
	rName := makeTestVariable(acctest.RandString(5))

	connNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsLDevCtxName := makeTestVariable(acctest.RandString(5))
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceContextConfig(fvTenantName, vnsLDevCtxName, connNameOrLbl),
			},
			{
				Config:      CreateAccLogicalInterfaceContextWithInValidParentDn(rName, connNameOrLbl),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, "description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, "annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, "name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, "l3_dest", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, "permit_log", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccLogicalInterfaceContextConfig(fvTenantName, vnsLDevCtxName, connNameOrLbl),
			},
		},
	})
}

func TestAccAciLogicalInterfaceContext_MultipleCreateDelete(t *testing.T) {
	connNameOrLbl := makeTestVariable(acctest.RandString(5))
	fvTenantName := makeTestVariable(acctest.RandString(5))
	vnsLDevCtxName := makeTestVariable(acctest.RandString(5))
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccLogicalInterfaceContextConfigMultiple(fvTenantName, vnsLDevCtxName, connNameOrLbl),
			},
		},
	})
}

func testAccCheckAciLogicalInterfaceContextExists(name string, logical_interface_context *models.LogicalInterfaceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Logical Interface Context %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Logical Interface Context dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		logical_interface_contextFound := models.LogicalInterfaceContextFromContainer(cont)
		if logical_interface_contextFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Logical Interface Context %s not found", rs.Primary.ID)
		}
		*logical_interface_context = *logical_interface_contextFound
		return nil
	}
}

func testAccCheckAciLogicalInterfaceContextDestroy(s *terraform.State) error {
	fmt.Println("=== STEP  testing logical_interface_context destroy")
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_logical_interface_context" {
			cont, err := client.Get(rs.Primary.ID)
			logical_interface_context := models.LogicalInterfaceContextFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Logical Interface Context %s Still exists", logical_interface_context.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLogicalInterfaceContextIdEqual(m1, m2 *models.LogicalInterfaceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("logical_interface_context DNs are not equal")
		}
		return nil
	}
}

func testAccCheckAciLogicalInterfaceContextIdNotEqual(m1, m2 *models.LogicalInterfaceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName == m2.DistinguishedName {
			return fmt.Errorf("logical_interface_context DNs are equal")
		}
		return nil
	}
}

func CreateLogicalInterfaceContextWithoutRequired(fvTenantName, vnsLDevCtxName, connNameOrLbl, attrName string) string {
	fmt.Println("=== STEP  Basic: testing logical_interface_context creation without ", attrName)
	rBlock := `
	
	resource "aci_tenant" "test" {
		name 		= "%s"
		
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	`
	switch attrName {
	case "logical_device_context_dn":
		rBlock += `
	resource "aci_logical_interface_context" "test" {
	#	logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}
		`
	case "conn_name_or_lbl":
		rBlock += `
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
	#	conn_name_or_lbl  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, fvTenantName, vnsLDevCtxName, connNameOrLbl)
}

func CreateAccLogicalInterfaceContextConfigWithRequiredParams(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_interface_context creation with updated naming arguments")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextConfigWithInvalidName(rName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_interface_context creation with invalid conn_name_or_lbl =", connNameOrLbl)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"

	}
	`, rName, rName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextConfig(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing logical_interface_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		tenant_dn = aci_tenant.test.id
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"

	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextConfigMultiple(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  testing multiple logical_interface_context creation with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s_${count.index}"
		count = 5
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextWithInValidParentDn(rName, connNameOrLbl string) string {
	fmt.Println("=== STEP  Negative Case: testing logical_interface_context creation with invalid parent Dn")
	resource := fmt.Sprintf(`
	resource "aci_tenant" "test"{
		name = "%s"
	}
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_tenant.test.id
		conn_name_or_lbl  = "%s"
	
	}
	`, rName, connNameOrLbl)
	return resource
}

func CreateAccLogicalInterfaceContextConfigWithOptionalValues(fvTenantName, vnsLDevCtxName, connNameOrLbl string) string {
	fmt.Println("=== STEP  Basic: testing logical_interface_context creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = "${aci_logical_device_context.test.id}"
		conn_name_or_lbl  = "%s"
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_logical_interface_context"
		l3_dest = "no"
		permit_log = "yes"
		
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl)

	return resource
}

func CreateAccLogicalInterfaceContextRemovingRequiredField() string {
	fmt.Println("=== STEP  Basic: testing logical_interface_context updation without required parameters")
	resource := fmt.Sprintf(`
	resource "aci_logical_interface_context" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_logical_interface_context"
		l3_dest = "no"
		permit_log = "yes"
		
	}
	`)

	return resource
}

func CreateAccLogicalInterfaceContextUpdatedAttr(fvTenantName, vnsLDevCtxName, connNameOrLbl, attribute, value string) string {
	fmt.Printf("=== STEP  testing logical_interface_context attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_tenant" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_logical_device_context" "test" {
		ctrct_name_or_lbl  = "%s"
		graph_name_or_lbl = "any"
		node_name_or_lbl  = "any"
		tenant_dn = aci_tenant.test.id
	}
	
	resource "aci_logical_interface_context" "test" {
		logical_device_context_dn  = aci_logical_device_context.test.id
		conn_name_or_lbl  = "%s"
		%s = "%s"
	}
	`, fvTenantName, vnsLDevCtxName, connNameOrLbl, attribute, value)
	return resource
}
