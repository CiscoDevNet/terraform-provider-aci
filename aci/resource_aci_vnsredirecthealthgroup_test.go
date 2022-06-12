package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL4goL7RedirectHealthGroup_Basic(t *testing.T) {
	var l4l7_redirect_health_group models.L4L7RedirectHealthGroup
	fv_tenant_name := acctest.RandString(5)
	vns_redirect_health_group_name := acctest.RandString(5)
	description := "l4-l7_redirect_health_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7RedirectHealthGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7RedirectHealthGroupConfig_basic(fv_tenant_name, vns_redirect_health_group_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7RedirectHealthGroupExists("aci_l4_l7_redirect_health_group.foo_l4-l7_redirect_health_group", &L4L7_redirect_health_group),
					testAccCheckAciL4L7RedirectHealthGroupAttributes(fv_tenant_name, vns_redirect_health_group_name, description, &L4L7_redirect_health_group),
				),
			},
		},
	})
}

func TestAccAciL4L7RedirectHealthGroup_Update(t *testing.T) {
	var l4l7_redirect_health_group models.L4L7RedirectHealthGroup
	fv_tenant_name := acctest.RandString(5)
	vns_redirect_health_group_name := acctest.RandString(5)
	description := "l4l7_redirect_health_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL4L7RedirectHealthGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL4L7RedirectHealthGroupConfig_basic(fv_tenant_name, vns_redirect_health_group_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7RedirectHealthGroupExists("aci_l4_l7_redirect_health_group.foo_l4-l7_redirect_health_group", &L4L7_redirect_health_group),
					testAccCheckAciL4L7RedirectHealthGroupAttributes(fv_tenant_name, vns_redirect_health_group_name, description, &L4L7_redirect_health_group),
				),
			},
			{
				Config: testAccCheckAciL4L7RedirectHealthGroupConfig_basic(fv_tenant_name, vns_redirect_health_group_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL4L7RedirectHealthGroupExists("aci_l4_l7_redirect_health_group.foo_l4-l7_redirect_health_group", &L4L7_redirect_health_group),
					testAccCheckAciL4L7RedirectHealthGroupAttributes(fv_tenant_name, vns_redirect_health_group_name, description, &L4L7_redirect_health_group),
				),
			},
		},
	})
}

func testAccCheckAciL4L7RedirectHealthGroupConfig_basic(fv_tenant_name, vns_redirect_health_group_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l4_l7_redirect_health_group" "foo_l4-l7_redirect_health_group" {
		name 		= "%s"
		description = "l4-l7_redirect_health_group created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	`, fv_tenant_name, vns_redirect_health_group_name)
}

func testAccCheckAciL4L7RedirectHealthGroupExists(name string, l4l7_redirect_health_group *models.L4L7RedirectHealthGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L4-L7 Redirect Health Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L4-L7 Redirect Health Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l4l7_redirect_health_groupFound := models.L4 - L7RedirectHealthGroupFromContainer(cont)
		if l4l7_redirect_health_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L4-L7 Redirect Health Group %s not found", rs.Primary.ID)
		}
		*l4 - l7_redirect_health_group = *l4l7_redirect_health_groupFound
		return nil
	}
}

func testAccCheckAciL4L7RedirectHealthGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_l4_l7_redirect_health_group" {
			cont, err := client.Get(rs.Primary.ID)
			l4l7_redirect_health_group := models.L4L7RedirectHealthGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L4-L7 Redirect Health Group %s Still exists", l4l7_redirect_health_group.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciL4L7RedirectHealthGroupAttributes(fv_tenant_name, vns_redirect_health_group_name, description string, l4l7_redirect_health_group *models.L4L7RedirectHealthGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vns_redirect_health_group_name != GetMOName(l4l7_redirect_health_group.DistinguishedName) {
			return fmt.Errorf("Bad vns_redirect_health_group %s", GetMOName(l4l7_redirect_health_group.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(l4l7_redirect_health_group.DistinguishedName, l4l7_redirect_health_group.Rn)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(l4l7_redirect_health_group.DistinguishedName, l4l7_redirect_health_group.Rn)))
		}
		if description != l4l7_redirect_health_group.Description {
			return fmt.Errorf("Bad l4-l7_redirect_health_group Description %s", l4l7_redirect_health_group.Description)
		}
		return nil
	}
}
