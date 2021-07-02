package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBgpRouteSummarization_Basic(t *testing.T) {
	var bgp_route_summarization models.BgpRouteSummarization
	description := "bg_proutesummarizationpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists("aci_bgp_route_summarization.foobgp_route_summarization", &bgp_route_summarization),
					testAccCheckAciBgpRouteSummarizationAttributes(description, &bgp_route_summarization),
				),
			},
		},
	})
}

func TestAccAciBgpRouteSummarization_update(t *testing.T) {
	var bgp_route_summarization models.BgpRouteSummarization
	description := "bg_proutesummarizationpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists("aci_bgp_route_summarization.foobgp_route_summarization", &bgp_route_summarization),
					testAccCheckAciBgpRouteSummarizationAttributes(description, &bgp_route_summarization),
				),
			},
			{
				Config: testAccCheckAciBgpRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpRouteSummarizationExists("aci_bgp_route_summarization.foobgp_route_summarization", &bgp_route_summarization),
					testAccCheckAciBgpRouteSummarizationAttributes(description, &bgp_route_summarization),
				),
			},
		},
	})
}

func testAccCheckAciBgpRouteSummarizationConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_bgp_route_summarization" "foobgp_route_summarization" {
		tenant_dn  = "${aci_tenant.demo_dev_tenant_test.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		attrmap  = "example"
  		ctrl = "as-set"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciBgpRouteSummarizationExists(name string, bgp_route_summarization *models.BgpRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BgpRouteSummarization %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BgpRouteSummarization dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_route_summarizationFound := models.BgpRouteSummarizationFromContainer(cont)
		if bgp_route_summarizationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BgpRouteSummarization %s not found", rs.Primary.ID)
		}
		*bgp_route_summarization = *bgp_route_summarizationFound
		return nil
	}
}

func testAccCheckAciBgpRouteSummarizationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bg_proutesummarizationpolicy" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_route_summarization := models.BgpRouteSummarizationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BgpRouteSummarization %s Still exists", bgp_route_summarization.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBgpRouteSummarizationAttributes(description string, bgp_route_summarization *models.BgpRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_route_summarization.Description {
			return fmt.Errorf("Bad bgp_route_summarization Description %s", bgp_route_summarization.Description)
		}

		if "example" != bgp_route_summarization.Name {
			return fmt.Errorf("Bad bgp_route_summarization name %s", bgp_route_summarization.Name)
		}

		if "example" != bgp_route_summarization.Annotation {
			return fmt.Errorf("Bad bgp_route_summarization annotation %s", bgp_route_summarization.Annotation)
		}

		if "example" != bgp_route_summarization.Attrmap {
			return fmt.Errorf("Bad bgp_route_summarization attrmap %s", bgp_route_summarization.Attrmap)
		}

		if "as-set" != bgp_route_summarization.Ctrl {
			return fmt.Errorf("Bad bgp_route_summarization ctrl %s", bgp_route_summarization.Ctrl)
		}

		if "example" != bgp_route_summarization.NameAlias {
			return fmt.Errorf("Bad bgp_route_summarization name_alias %s", bgp_route_summarization.NameAlias)
		}

		return nil
	}
}
