package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSPANSourcedestinationGroupMatchLabel_Basic(t *testing.T) {
	var span_sourcedestination_group_match_label models.SPANSourcedestinationGroupMatchLabel
	description := "span_source-destination_group_match_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANSourcedestinationGroupMatchLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists("aci_span_sourcedestination_group_match_label.foospan_sourcedestination_group_match_label", &span_sourcedestination_group_match_label),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelAttributes(description, &span_sourcedestination_group_match_label),
				),
			},
			{
				ResourceName:      "aci_span_sourcedestination_group_match_label",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciSPANSourcedestinationGroupMatchLabel_update(t *testing.T) {
	var span_sourcedestination_group_match_label models.SPANSourcedestinationGroupMatchLabel
	description := "span_source-destination_group_match_label created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANSourcedestinationGroupMatchLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists("aci_span_sourcedestination_group_match_label.foospan_sourcedestination_group_match_label", &span_sourcedestination_group_match_label),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelAttributes(description, &span_sourcedestination_group_match_label),
				),
			},
			{
				Config: testAccCheckAciSPANSourcedestinationGroupMatchLabelConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourcedestinationGroupMatchLabelExists("aci_span_sourcedestination_group_match_label.foospan_sourcedestination_group_match_label", &span_sourcedestination_group_match_label),
					testAccCheckAciSPANSourcedestinationGroupMatchLabelAttributes(description, &span_sourcedestination_group_match_label),
				),
			},
		},
	})
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_span_sourcedestination_group_match_label" "foospan_sourcedestination_group_match_label" {
		  span_source_group_dn  = "${aci_span_source_group.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		  tag  = "yellow"
		}
	`, description)
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelExists(name string, span_sourcedestination_group_match_label *models.SPANSourcedestinationGroupMatchLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPAN Source-destination Group Match Label %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPAN Source-destination Group Match Label dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_sourcedestination_group_match_labelFound := models.SPANSourcedestinationGroupMatchLabelFromContainer(cont)
		if span_sourcedestination_group_match_labelFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPAN Source-destination Group Match Label %s not found", rs.Primary.ID)
		}
		*span_sourcedestination_group_match_label = *span_sourcedestination_group_match_labelFound
		return nil
	}
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_span_source-destination_group_match_label" {
			cont, err := client.Get(rs.Primary.ID)
			span_sourcedestination_group_match_label := models.SPANSourcedestinationGroupMatchLabelFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPAN Source-destination Group Match Label %s Still exists", span_sourcedestination_group_match_label.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSPANSourcedestinationGroupMatchLabelAttributes(description string, span_sourcedestination_group_match_label *models.SPANSourcedestinationGroupMatchLabel) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != span_sourcedestination_group_match_label.Description {
			return fmt.Errorf("Bad span_sourcedestination_group_match_label Description %s", span_sourcedestination_group_match_label.Description)
		}

		if "example" != span_sourcedestination_group_match_label.Name {
			return fmt.Errorf("Bad span_sourcedestination_group_match_label name %s", span_sourcedestination_group_match_label.Name)
		}

		if "example" != span_sourcedestination_group_match_label.Annotation {
			return fmt.Errorf("Bad span_sourcedestination_group_match_label annotation %s", span_sourcedestination_group_match_label.Annotation)
		}

		if "example" != span_sourcedestination_group_match_label.NameAlias {
			return fmt.Errorf("Bad span_sourcedestination_group_match_label name_alias %s", span_sourcedestination_group_match_label.NameAlias)
		}

		if "yellow" != span_sourcedestination_group_match_label.Tag {
			return fmt.Errorf("Bad span_sourcedestination_group_match_label tag %s", span_sourcedestination_group_match_label.Tag)
		}

		return nil
	}
}
