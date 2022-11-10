package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVXLANPool_Basic(t *testing.T) {
	var vxlan_pool models.VXLANPool
	description := "vxlan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVXLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists("aci_vxlan_pool.foovxlan_pool", &vxlan_pool),
					testAccCheckAciVXLANPoolAttributes(description, &vxlan_pool),
				),
			},
		},
	})
}

func TestAccAciVXLANPool_update(t *testing.T) {
	var vxlan_pool models.VXLANPool
	description := "vxlan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVXLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVXLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists("aci_vxlan_pool.foovxlan_pool", &vxlan_pool),
					testAccCheckAciVXLANPoolAttributes(description, &vxlan_pool),
				),
			},
			{
				Config: testAccCheckAciVXLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVXLANPoolExists("aci_vxlan_pool.foovxlan_pool", &vxlan_pool),
					testAccCheckAciVXLANPoolAttributes(description, &vxlan_pool),
				),
			},
		},
	})
}

func testAccCheckAciVXLANPoolConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vxlan_pool" "foovxlan_pool" {
			description = "%s"
			name  = "example"
			annotation  = "example"
			name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciVXLANPoolExists(name string, vxlan_pool *models.VXLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VXLAN Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VXLAN Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vxlan_poolFound := models.VXLANPoolFromContainer(cont)
		if vxlan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VXLAN Pool %s not found", rs.Primary.ID)
		}
		*vxlan_pool = *vxlan_poolFound
		return nil
	}
}

func testAccCheckAciVXLANPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vxlan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vxlan_pool := models.VXLANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VXLAN Pool %s Still exists", vxlan_pool.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVXLANPoolAttributes(description string, vxlan_pool *models.VXLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vxlan_pool.Description {
			return fmt.Errorf("Bad vxlan_pool Description %s", vxlan_pool.Description)
		}

		if "example" != vxlan_pool.Name {
			return fmt.Errorf("Bad vxlan_pool name %s", vxlan_pool.Name)
		}

		if "example" != vxlan_pool.Annotation {
			return fmt.Errorf("Bad vxlan_pool annotation %s", vxlan_pool.Annotation)
		}

		if "example" != vxlan_pool.NameAlias {
			return fmt.Errorf("Bad vxlan_pool name_alias %s", vxlan_pool.NameAlias)
		}

		return nil
	}
}
