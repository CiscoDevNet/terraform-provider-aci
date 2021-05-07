package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciOspfRouteSummarization_Basic(t *testing.T) {
	var ospf_route_summarization models.OspfRouteSummarization
	description := "osp_froutesummarizationpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOspfRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists("aci_ospf_route_summarization.fooospf_route_summarization", &ospf_route_summarization),
					testAccCheckAciOspfRouteSummarizationAttributes(description, &ospf_route_summarization),
				),
			},
		},
	})
}

func TestAccAciOspfRouteSummarization_update(t *testing.T) {
	var ospf_route_summarization models.OspfRouteSummarization
	description := "osp_froutesummarizationpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOspfRouteSummarizationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOspfRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists("aci_ospf_route_summarization.fooospf_route_summarization", &ospf_route_summarization),
					testAccCheckAciOspfRouteSummarizationAttributes(description, &ospf_route_summarization),
				),
			},
			{
				Config: testAccCheckAciOspfRouteSummarizationConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOspfRouteSummarizationExists("aci_ospf_route_summarization.fooospf_route_summarization", &ospf_route_summarization),
					testAccCheckAciOspfRouteSummarizationAttributes(description, &ospf_route_summarization),
				),
			},
		},
	})
}

func testAccCheckAciOspfRouteSummarizationConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_ospf_route_summarization" "fooospf_route_summarization" {
		tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		cost = "1"
  		inter_area_enabled = "no"
  		name_alias  = "example"
  		tag  = "1"
	}
	`, description)
}

func testAccCheckAciOspfRouteSummarizationExists(name string, ospf_route_summarization *models.OspfRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("OspfRouteSummarization %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OspfRouteSummarization dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_route_summarizationFound := models.OspfRouteSummarizationFromContainer(cont)
		if ospf_route_summarizationFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("OspfRouteSummarization %s not found", rs.Primary.ID)
		}
		*ospf_route_summarization = *ospf_route_summarizationFound
		return nil
	}
}

func testAccCheckAciOspfRouteSummarizationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_osp_froutesummarizationpolicy" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_route_summarization := models.OspfRouteSummarizationFromContainer(cont)
			if err == nil {
				return fmt.Errorf("OspfRouteSummarization %s Still exists", ospf_route_summarization.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOspfRouteSummarizationAttributes(description string, ospf_route_summarization *models.OspfRouteSummarization) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != ospf_route_summarization.Description {
			return fmt.Errorf("Bad ospf_route_summarization Description %s", ospf_route_summarization.Description)
		}

		if "example" != ospf_route_summarization.Name {
			return fmt.Errorf("Bad ospf_route_summarization name %s", ospf_route_summarization.Name)
		}

		if "example" != ospf_route_summarization.Annotation {
			return fmt.Errorf("Bad ospf_route_summarization annotation %s", ospf_route_summarization.Annotation)
		}

		if "1" != ospf_route_summarization.Cost {
			return fmt.Errorf("Bad ospf_route_summarization cost %s", ospf_route_summarization.Cost)
		}

		if "no" != ospf_route_summarization.InterAreaEnabled {
			return fmt.Errorf("Bad ospf_route_summarization inter_area_enabled %s", ospf_route_summarization.InterAreaEnabled)
		}

		if "example" != ospf_route_summarization.NameAlias {
			return fmt.Errorf("Bad ospf_route_summarization name_alias %s", ospf_route_summarization.NameAlias)
		}

		if "1" != ospf_route_summarization.Tag {
			return fmt.Errorf("Bad ospf_route_summarization tag %s", ospf_route_summarization.Tag)
		}

		return nil
	}
}
