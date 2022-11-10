package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVLANPool_Basic(t *testing.T) {
	var vlan_pool models.VLANPool
	description := "vlan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists("aci_vlan_pool.foovlan_pool", &vlan_pool),
					testAccCheckAciVLANPoolAttributes(description, &vlan_pool),
				),
			},
		},
	})
}

func TestAccAciVLANPool_update(t *testing.T) {
	var vlan_pool models.VLANPool
	description := "vlan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVLANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists("aci_vlan_pool.foovlan_pool", &vlan_pool),
					testAccCheckAciVLANPoolAttributes(description, &vlan_pool),
				),
			},
			{
				Config: testAccCheckAciVLANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVLANPoolExists("aci_vlan_pool.foovlan_pool", &vlan_pool),
					testAccCheckAciVLANPoolAttributes(description, &vlan_pool),
				),
			},
		},
	})
}

func testAccCheckAciVLANPoolConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vlan_pool" "foovlan_pool" {
		description = "%s"
		name  = "example"
		alloc_mode  = "static"
		annotation  = "example"
		name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciVLANPoolExists(name string, vlan_pool *models.VLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VLAN Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VLAN Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vlan_poolFound := models.VLANPoolFromContainer(cont)
		if vlan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VLAN Pool %s not found", rs.Primary.ID)
		}
		*vlan_pool = *vlan_poolFound
		return nil
	}
}

func testAccCheckAciVLANPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vlan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vlan_pool := models.VLANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VLAN Pool %s Still exists", vlan_pool.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVLANPoolAttributes(description string, vlan_pool *models.VLANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vlan_pool.Description {
			return fmt.Errorf("Bad vlan_pool Description %s", vlan_pool.Description)
		}

		if "example" != vlan_pool.Name {
			return fmt.Errorf("Bad vlan_pool name %s", vlan_pool.Name)
		}

		if "static" != vlan_pool.AllocMode {
			return fmt.Errorf("Bad vlan_pool alloc_mode %s", vlan_pool.AllocMode)
		}

		// if "example" != vlan_pool.AllocMode {
		// 	return fmt.Errorf("Bad vlan_pool alloc_mode %s", vlan_pool.AllocMode)
		// }

		if "example" != vlan_pool.Annotation {
			return fmt.Errorf("Bad vlan_pool annotation %s", vlan_pool.Annotation)
		}

		if "example" != vlan_pool.NameAlias {
			return fmt.Errorf("Bad vlan_pool name_alias %s", vlan_pool.NameAlias)
		}

		return nil
	}
}
