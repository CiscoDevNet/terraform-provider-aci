package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciFunctionNode_Basic(t *testing.T) {
	var function_node models.FunctionNode
	description := "function_node created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFunctionNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists("aci_function_node.foofunction_node", &function_node),
					testAccCheckAciFunctionNodeAttributes(description, &function_node),
				),
			},
		},
	})
}

func TestAccAciFunctionNode_update(t *testing.T) {
	var function_node models.FunctionNode
	description := "function_node created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFunctionNodeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFunctionNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists("aci_function_node.foofunction_node", &function_node),
					testAccCheckAciFunctionNodeAttributes(description, &function_node),
				),
			},
			{
				Config: testAccCheckAciFunctionNodeConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFunctionNodeExists("aci_function_node.foofunction_node", &function_node),
					testAccCheckAciFunctionNodeAttributes(description, &function_node),
				),
			},
		},
	})
}

func testAccCheckAciFunctionNodeConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_function_node" "foofunction_node" {
		  l4_l7_service_graph_template_dn  = "uni/tn-check_tenantnk/AbsGraph-one"
		  
		description = "%s"
		
		name  = "xyz"
		  annotation  = "example"
		  func_template_type  = "OTHER"
		  func_type  = "None"
		  is_copy  = "no"
		  managed  = "yes"
		  name_alias  = "example_alias"
		  routing_mode  = "unspecified"
		  sequence_number  = "3"
		  share_encap  = "yes"
		}
	`, description)
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

func testAccCheckAciFunctionNodeAttributes(description string, function_node *models.FunctionNode) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != function_node.Description {
			return fmt.Errorf("Bad function_node Description %s", function_node.Description)
		}

		if "xyz" != function_node.Name {
			return fmt.Errorf("Bad function_node name %s", function_node.Name)
		}

		if "example" != function_node.Annotation {
			return fmt.Errorf("Bad function_node annotation %s", function_node.Annotation)
		}

		if "OTHER" != function_node.FuncTemplateType {
			return fmt.Errorf("Bad function_node func_template_type %s", function_node.FuncTemplateType)
		}

		if "None" != function_node.FuncType {
			return fmt.Errorf("Bad function_node func_type %s", function_node.FuncType)
		}

		if "no" != function_node.IsCopy {
			return fmt.Errorf("Bad function_node is_copy %s", function_node.IsCopy)
		}

		if "yes" != function_node.Managed {
			return fmt.Errorf("Bad function_node managed %s", function_node.Managed)
		}

		if "example_alias" != function_node.NameAlias {
			return fmt.Errorf("Bad function_node name_alias %s", function_node.NameAlias)
		}

		if "unspecified" != function_node.RoutingMode {
			return fmt.Errorf("Bad function_node routing_mode %s", function_node.RoutingMode)
		}

		if "3" != function_node.SequenceNumber {
			return fmt.Errorf("Bad function_node sequence_number %s", function_node.SequenceNumber)
		}

		if "yes" != function_node.ShareEncap {
			return fmt.Errorf("Bad function_node share_encap %s", function_node.ShareEncap)
		}

		return nil
	}
}
