package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciLogicalDeviceContext_Basic(t *testing.T) {
	var logical_device_context models.LogicalDeviceContext
	description := "logical_device_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalDeviceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists("aci_logical_device_context.test", &logical_device_context),
					testAccCheckAciLogicalDeviceContextAttributes(description, &logical_device_context),
				),
			},
		},
	})
}

func TestAccAciLogicalDeviceContext_update(t *testing.T) {
	var logical_device_context models.LogicalDeviceContext
	description := "logical_device_context created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLogicalDeviceContextDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLogicalDeviceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists("aci_logical_device_context.test", &logical_device_context),
					testAccCheckAciLogicalDeviceContextAttributes(description, &logical_device_context),
				),
			},
			{
				Config: testAccCheckAciLogicalDeviceContextConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLogicalDeviceContextExists("aci_logical_device_context.test", &logical_device_context),
					testAccCheckAciLogicalDeviceContextAttributes(description, &logical_device_context),
				),
			},
		},
	})
}

func testAccCheckAciLogicalDeviceContextConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_logical_device_context" "test" {
		tenant_dn  = "${aci_tenant.tenentcheck.id}"
		description = "%s"
		annotation  = "test"
		context  = "test"
		ctrct_name_or_lbl  = "any"
		graph_name_or_lbl  = "any"
		name_alias  = "test"
		node_name_or_lbl  = "any"
	}
	`, description)
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

func testAccCheckAciLogicalDeviceContextAttributes(description string, logical_device_context *models.LogicalDeviceContext) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != logical_device_context.Description {
			return fmt.Errorf("Bad logical_device_context Description %s", logical_device_context.Description)
		}

		if "any" != logical_device_context.CtrctNameOrLbl {
			return fmt.Errorf("Bad logical_device_context ctrct_name_or_lbl %s", logical_device_context.CtrctNameOrLbl)
		}

		if "any" != logical_device_context.GraphNameOrLbl {
			return fmt.Errorf("Bad logical_device_context graph_name_or_lbl %s", logical_device_context.GraphNameOrLbl)
		}

		if "any" != logical_device_context.NodeNameOrLbl {
			return fmt.Errorf("Bad logical_device_context node_name_or_lbl %s", logical_device_context.NodeNameOrLbl)
		}

		if "test" != logical_device_context.Annotation {
			return fmt.Errorf("Bad logical_device_context annotation %s", logical_device_context.Annotation)
		}

		if "test" != logical_device_context.Context {
			return fmt.Errorf("Bad logical_device_context context %s", logical_device_context.Context)
		}

		if "test" != logical_device_context.NameAlias {
			return fmt.Errorf("Bad logical_device_context name_alias %s", logical_device_context.NameAlias)
		}

		return nil
	}
}
