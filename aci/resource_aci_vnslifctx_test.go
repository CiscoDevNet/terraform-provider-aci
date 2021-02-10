package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciLogicalInterfaceContext_Basic(t *testing.T) {
	var logical_interface_context models.LogicalInterfaceContext
	description := "logical_interface_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists("aci_logical_interface_context.foological_interface_context", &logical_interface_context),
					testAccCheckAciLogicalInterfaceContextAttributes(description, &logical_interface_context),
				),
			},
		},
	})
}

func TestAccAciLogicalInterfaceContext_update(t *testing.T) {
	var logical_interface_context models.LogicalInterfaceContext
	description := "logical_interface_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalInterfaceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalInterfaceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists("aci_logical_interface_context.foological_interface_context", &logical_interface_context),
					testAccCheckAciLogicalInterfaceContextAttributes(description, &logical_interface_context),
				),
			},
			{
				Config: testAccCheckAciLogicalInterfaceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalInterfaceContextExists("aci_logical_interface_context.foological_interface_context", &logical_interface_context),
					testAccCheckAciLogicalInterfaceContextAttributes(description, &logical_interface_context),
				),
			},
		},
	})
}

func testAccCheckAciLogicalInterfaceContextConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_logical_interface_context" "foological_interface_context" {
		logical_device_context_dn  = "uni/tn-check_context_tenant/ldevCtx-c-x-g-y-n-z"
		description = "%s"
		annotation  = "anno"
		conn_name_or_lbl  = "lbl"
		l3_dest  = "no"
		name_alias  = "alias"
		permit_log  = "no"
	}
	`, description)
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

func testAccCheckAciLogicalInterfaceContextAttributes(description string, logical_interface_context *models.LogicalInterfaceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != logical_interface_context.Description {
			return fmt.Errorf("Bad logical_interface_context Description %s", logical_interface_context.Description)
		}

		if "anno" != logical_interface_context.Annotation {
			return fmt.Errorf("Bad logical_interface_context annotation %s", logical_interface_context.Annotation)
		}

		if "lbl" != logical_interface_context.ConnNameOrLbl {
			return fmt.Errorf("Bad logical_interface_context conn_name_or_lbl %s", logical_interface_context.ConnNameOrLbl)
		}

		if "no" != logical_interface_context.L3Dest {
			return fmt.Errorf("Bad logical_interface_context l3_dest %s", logical_interface_context.L3Dest)
		}

		if "alias" != logical_interface_context.NameAlias {
			return fmt.Errorf("Bad logical_interface_context name_alias %s", logical_interface_context.NameAlias)
		}

		if "no" != logical_interface_context.PermitLog {
			return fmt.Errorf("Bad logical_interface_context permit_log %s", logical_interface_context.PermitLog)
		}

		return nil
	}
}
