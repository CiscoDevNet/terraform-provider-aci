package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL4L7ServiceGraphTemplate_Basic(t *testing.T) {
	var l4l7SGraphTemplate models.L4L7ServiceGraphTemplate
	description := "l4-l7_service_graph_template created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4l7SGraphTemplate),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4l7SGraphTemplate),
				),
			},
		},
	})
}

func TestAccAciL4L7ServiceGraphTemplate_update(t *testing.T) {
	var l4l7SGraphTemplate models.L4L7ServiceGraphTemplate
	description := "l4-l7_service_graph_template created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7ServiceGraphTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4l7SGraphTemplate),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4l7SGraphTemplate),
				),
			},
			{
				Config: testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7ServiceGraphTemplateExists("aci_l4_l7_service_graph_template.test", &l4l7SGraphTemplate),
					testAccCheckAciL4L7ServiceGraphTemplateAttributes(description, &l4l7SGraphTemplate),
				),
			},
		},
	})
}

func testAccCheckAciL4L7ServiceGraphTemplateConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "example" {
		name = "testacc"
	}
	resource "aci_l4_l7_service_graph_template" "test" {
		tenant_dn   = aci_tenant.example.id
		description = "%s"
		name        = "test"
		annotation  = "example"
		name_alias  = "example"
		l4_l7_service_graph_template_type = "legacy"
		ui_template_type                  = "ONE_NODE_ADC_ONE_ARM"
	}
	`, description)
}

func testAccCheckAciL4L7ServiceGraphTemplateExists(name string, l4l7SGraphTemplate *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
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

		l4l7SGraphTemplateFound := models.L4L7ServiceGraphTemplateFromContainer(cont)
		if l4l7SGraphTemplateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L4-L7 Service Graph Template %s not found", rs.Primary.ID)
		}
		*l4l7SGraphTemplate = *l4l7SGraphTemplateFound
		return nil
	}
}

func testAccCheckAciL4L7ServiceGraphTemplateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l4-l7_service_graph_template" {
			cont, err := client.Get(rs.Primary.ID)
			l4l7SGraphTemplate := models.L4L7ServiceGraphTemplateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L4-L7 Service Graph Template %s Still exists", l4l7SGraphTemplate.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL4L7ServiceGraphTemplateAttributes(description string, l4l7SGraphTemplate *models.L4L7ServiceGraphTemplate) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l4l7SGraphTemplate.Description {
			return fmt.Errorf("Bad l4 l7 service graph template Description %s", l4l7SGraphTemplate.Description)
		}

		if "test" != l4l7SGraphTemplate.Name {
			return fmt.Errorf("Bad l4 l7 service graph template name %s", l4l7SGraphTemplate.Name)
		}

		if "example" != l4l7SGraphTemplate.Annotation {
			return fmt.Errorf("Bad l4 l7 service graph template annotation %s", l4l7SGraphTemplate.Annotation)
		}

		if "example" != l4l7SGraphTemplate.NameAlias {
			return fmt.Errorf("Bad l4 l7 service graph template name_alias %s", l4l7SGraphTemplate.NameAlias)
		}

		if "legacy" != l4l7SGraphTemplate.L4L7ServiceGraphTemplate_type {
			return fmt.Errorf("Bad l4 l7 service graph template type %s", l4l7SGraphTemplate.L4L7ServiceGraphTemplate_type)
		}

		if "ONE_NODE_ADC_ONE_ARM" != l4l7SGraphTemplate.UiTemplateType {
			return fmt.Errorf("Bad l4 l7 service graph template ui_template_type %s", l4l7SGraphTemplate.UiTemplateType)
		}

		return nil
	}
}
