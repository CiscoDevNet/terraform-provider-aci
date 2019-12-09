package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciFilter_Basic(t *testing.T) {
	var filter models.Filter
	description := "filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterConfig_basic(description, "alias_filter"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(description, "alias_filter", &filter),
				),
			},
			{
				ResourceName:      "aci_filter",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFilter_update(t *testing.T) {
	var filter models.Filter
	description := "filter created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFilterConfig_basic(description, "alias_filter"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(description, "alias_filter", &filter),
				),
			},
			{
				Config: testAccCheckAciFilterConfig_basic(description, "update_filter"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFilterExists("aci_filter.foofilter", &filter),
					testAccCheckAciFilterAttributes(description, "update_filter", &filter),
				),
			},
		},
	})
}

func testAccCheckAciFilterConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_filter" "foofilter" {
		tenant_dn   = "${aci_tenant.example.id}"
		description = "%s"
		name        = "demo_filter"
		annotation  = "tag_filter"
		name_alias  = "%s"
	}
	  
	`, description, name_alias)
}

func testAccCheckAciFilterExists(name string, filter *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Filter %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Filter dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		filterFound := models.FilterFromContainer(cont)
		if filterFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Filter %s not found", rs.Primary.ID)
		}
		*filter = *filterFound
		return nil
	}
}

func testAccCheckAciFilterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_filter" {
			cont, err := client.Get(rs.Primary.ID)
			filter := models.FilterFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Filter %s Still exists", filter.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFilterAttributes(description, name_alias string, filter *models.Filter) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != filter.Description {
			return fmt.Errorf("Bad filter Description %s", filter.Description)
		}

		if "demo_filter" != filter.Name {
			return fmt.Errorf("Bad filter name %s", filter.Name)
		}

		if "tag_filter" != filter.Annotation {
			return fmt.Errorf("Bad filter annotation %s", filter.Annotation)
		}

		if name_alias != filter.NameAlias {
			return fmt.Errorf("Bad filter name_alias %s", filter.NameAlias)
		}

		return nil
	}
}
