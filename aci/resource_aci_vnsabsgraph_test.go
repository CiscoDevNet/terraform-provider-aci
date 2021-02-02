package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL4L7ServiceGraphTemplate_Basic(t *testing.T) {
	var l4_l7_service_graph_template models.L4L7ServiceGraphTemplate
	description := "l4-l7_service_graph_template created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4_l7_service_graph_template),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4_l7_service_graph_template),
				),
			},
		},
	})
}

func TestAccAciL4L7ServiceGraphTemplate_update(t *testing.T) {
	var l4_l7_service_graph_template models.L4L7ServiceGraphTemplate
	description := "l4-l7_service_graph_template created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4_l7_service_graph_template),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4_l7_service_graph_template),
				),
			},
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4_l7_service_graph_template),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4_l7_service_graph_template),
				),
			},
		},
	})
}

func testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
		l4-l7_service_graph_template_type  = "example"
		ui_template_type  = "example"
	}
	`, description)
}

func testAccCheckAciL4L7ServiceGraphTemplateExists(name string, l4_l7_service_graph_template *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L4-L7 Service Graph Template %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L4-L7 Service Graph Template dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l4_l7_service_graph_templateFound := models.L4L7ServiceGraphTemplateFromContainer(cont)
		if l4_l7_service_graph_templateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L4-L7 Service Graph Template %s not found", rs.Primary.ID)
		}
		*l4_l7_service_graph_template = *l4_l7_service_graph_templateFound
		return nil
	}
}

func testAccCheckAciL4L7ServiceGraphTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l4-l7_service_graph_template" {
			cont, err := client.Get(rs.Primary.ID)
			l4_l7_service_graph_template := models.L4L7ServiceGraphTemplateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L4-L7 Service Graph Template %s Still exists", l4_l7_service_graph_template.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL4L7ServiceGraphTemplateAttributes(description string, l4_l7_service_graph_template *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l4_l7_service_graph_template.Description {
			return fmt.Errorf("Bad l4_l7_service_graph_template Description %s", l4_l7_service_graph_template.Description)
		}

		if "example" != l4_l7_service_graph_template.Name {
			return fmt.Errorf("Bad l4_l7_service_graph_template name %s", l4_l7_service_graph_template.Name)
		}

		if "example" != l4_l7_service_graph_template.Annotation {
			return fmt.Errorf("Bad l4_l7_service_graph_template annotation %s", l4_l7_service_graph_template.Annotation)
		}

		if "example" != l4_l7_service_graph_template.NameAlias {
			return fmt.Errorf("Bad l4_l7_service_graph_template name_alias %s", l4_l7_service_graph_template.NameAlias)
		}

		if "example" != l4_l7_service_graph_template.L4L7ServiceGraphTemplate_type {
			return fmt.Errorf("Bad l4_l7_service_graph_template l4_l7_service_graph_template_type %s", l4_l7_service_graph_template.L4L7ServiceGraphTemplate_type)
		}

		if "example" != l4_l7_service_graph_template.UiTemplateType {
			return fmt.Errorf("Bad l4_l7_service_graph_template ui_template_type %s", l4_l7_service_graph_template.UiTemplateType)
		}

		return nil
	}
}
