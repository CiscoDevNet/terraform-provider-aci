package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciSPANDestinationGroup_Basic(t *testing.T) {
	var span_destination_group models.SPANDestinationGroup
	description := "span_destination_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists("aci_span_destination_group.foospan_destination_group", &span_destination_group),
					testAccCheckAciSPANDestinationGroupAttributes(description, &span_destination_group),
				),
			},
			{
				ResourceName:      "aci_span_destination_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciSPANDestinationGroup_update(t *testing.T) {
	var span_destination_group models.SPANDestinationGroup
	description := "span_destination_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANDestinationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists("aci_span_destination_group.foospan_destination_group", &span_destination_group),
					testAccCheckAciSPANDestinationGroupAttributes(description, &span_destination_group),
				),
			},
			{
				Config: testAccCheckAciSPANDestinationGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANDestinationGroupExists("aci_span_destination_group.foospan_destination_group", &span_destination_group),
					testAccCheckAciSPANDestinationGroupAttributes(description, &span_destination_group),
				),
			},
		},
	})
}

func testAccCheckAciSPANDestinationGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_span_destination_group" "foospan_destination_group" {
		  tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciSPANDestinationGroupExists(name string, span_destination_group *models.SPANDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPAN Destination Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPAN Destination Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_destination_groupFound := models.SPANDestinationGroupFromContainer(cont)
		if span_destination_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPAN Destination Group %s not found", rs.Primary.ID)
		}
		*span_destination_group = *span_destination_groupFound
		return nil
	}
}

func testAccCheckAciSPANDestinationGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_span_destination_group" {
			cont, err := client.Get(rs.Primary.ID)
			span_destination_group := models.SPANDestinationGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPAN Destination Group %s Still exists", span_destination_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSPANDestinationGroupAttributes(description string, span_destination_group *models.SPANDestinationGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != span_destination_group.Description {
			return fmt.Errorf("Bad span_destination_group Description %s", span_destination_group.Description)
		}

		if "example" != span_destination_group.Name {
			return fmt.Errorf("Bad span_destination_group name %s", span_destination_group.Name)
		}

		if "example" != span_destination_group.Annotation {
			return fmt.Errorf("Bad span_destination_group annotation %s", span_destination_group.Annotation)
		}

		if "example" != span_destination_group.NameAlias {
			return fmt.Errorf("Bad span_destination_group name_alias %s", span_destination_group.NameAlias)
		}

		return nil
	}
}
