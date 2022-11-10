package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTag_Basic(t *testing.T) {
	var tag models.Tag
	planner_tenant_tmpl_name := acctest.RandString(5)
	planner_match_tenant_name := acctest.RandString(5)
	fv_tenant_name := acctest.RandString(5)
	tag_tag_name := acctest.RandString(5)
	description := "tag created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTagConfig_basic(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists("aci_tag.footag", &tag),
					testAccCheckAciTagAttributes(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name, description, &tag),
				),
			},
		},
	})
}

func TestAccAciTag_Update(t *testing.T) {
	var tag models.Tag
	planner_tenant_tmpl_name := acctest.RandString(5)
	planner_match_tenant_name := acctest.RandString(5)
	fv_tenant_name := acctest.RandString(5)
	tag_tag_name := acctest.RandString(5)
	description := "tag created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTagConfig_basic(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists("aci_tag.footag", &tag),
					testAccCheckAciTagAttributes(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name, description, &tag),
				),
			},
			{
				Config: testAccCheckAciTagConfig_basic(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTagExists("aci_tag.footag", &tag),
					testAccCheckAciTagAttributes(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name, description, &tag),
				),
			},
		},
	})
}

func testAccCheckAciTagConfig_basic(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name string) string {
	return fmt.Sprintf(`

	resource "aci_optimizer_tenant_template" "foooptimizer_tenant_template" {
		name 		= "%s"
		description = "optimizer_tenant_template created while acceptance testing"

	}

	resource "aci_tenantaffinity" "footenantaffinity" {
		name 		= "%s"
		description = "tenantaffinity created while acceptance testing"
		optimizer_tenant_template_dn = aci_optimizer_tenant_template.foooptimizer_tenant_template.id
	}

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"
		tenantaffinity_dn = aci_tenantaffinity.footenantaffinity.id
	}

	resource "aci_tag" "footag" {
		name 		= "%s"
		description = "tag created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name)
}

func testAccCheckAciTagExists(name string, tag *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Tag %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Tag dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		tagFound := models.TagFromContainer(cont)
		if tagFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Tag %s not found", rs.Primary.ID)
		}
		*tag = *tagFound
		return nil
	}
}

func testAccCheckAciTagDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tag" {
			cont, err := client.Get(rs.Primary.ID)
			tag := models.TagFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Tag %s Still exists", tag.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTagAttributes(planner_tenant_tmpl_name, planner_match_tenant_name, fv_tenant_name, tag_tag_name, description string, tag *models.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if tag_tag_name != GetMOName(tag.DistinguishedName) {
			return fmt.Errorf("Bad tag_tag %s", GetMOName(tag.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(tag.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(tag.DistinguishedName)))
		}
		if description != tag.Description {
			return fmt.Errorf("Bad tag Description %s", tag.Description)
		}
		return nil
	}
}
