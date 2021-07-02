package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSPANSourceGroup_Basic(t *testing.T) {
	var span_source_group models.SPANSourceGroup
	description := "span_source_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANSourceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists("aci_span_source_group.foospan_source_group", &span_source_group),
					testAccCheckAciSPANSourceGroupAttributes(description, &span_source_group),
				),
			},
		},
	})
}

func TestAccAciSPANSourceGroup_update(t *testing.T) {
	var span_source_group models.SPANSourceGroup
	description := "span_source_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSPANSourceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSPANSourceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists("aci_span_source_group.foospan_source_group", &span_source_group),
					testAccCheckAciSPANSourceGroupAttributes(description, &span_source_group),
				),
			},
			{
				Config: testAccCheckAciSPANSourceGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSPANSourceGroupExists("aci_span_source_group.foospan_source_group", &span_source_group),
					testAccCheckAciSPANSourceGroupAttributes(description, &span_source_group),
				),
			},
		},
	})
}

func testAccCheckAciSPANSourceGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_span_source_group" "foospan_source_group" {
		tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		name  = "example"
		admin_st  = "enabled"
	    annotation  = "example"
    	name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciSPANSourceGroupExists(name string, span_source_group *models.SPANSourceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SPAN Source Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SPAN Source Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		span_source_groupFound := models.SPANSourceGroupFromContainer(cont)
		if span_source_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SPAN Source Group %s not found", rs.Primary.ID)
		}
		*span_source_group = *span_source_groupFound
		return nil
	}
}

func testAccCheckAciSPANSourceGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_span_source_group" {
			cont, err := client.Get(rs.Primary.ID)
			span_source_group := models.SPANSourceGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SPAN Source Group %s Still exists", span_source_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciSPANSourceGroupAttributes(description string, span_source_group *models.SPANSourceGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != span_source_group.Description {
			return fmt.Errorf("Bad span_source_group Description %s", span_source_group.Description)
		}

		if "example" != span_source_group.Name {
			return fmt.Errorf("Bad span_source_group name %s", span_source_group.Name)
		}

		if "enabled" != span_source_group.AdminSt {
			return fmt.Errorf("Bad span_source_group admin_st %s", span_source_group.AdminSt)
		}

		if "example" != span_source_group.Annotation {
			return fmt.Errorf("Bad span_source_group annotation %s", span_source_group.Annotation)
		}

		if "example" != span_source_group.NameAlias {
			return fmt.Errorf("Bad span_source_group name_alias %s", span_source_group.NameAlias)
		}

		return nil
	}
}
