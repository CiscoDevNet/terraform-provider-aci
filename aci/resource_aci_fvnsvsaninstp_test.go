package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVSANPool_Basic(t *testing.T) {
	var vsan_pool models.VSANPool
	description := "vsan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVSANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists("aci_vsan_pool.foovsan_pool", &vsan_pool),
					testAccCheckAciVSANPoolAttributes(description, &vsan_pool),
				),
			},
		},
	})
}

func TestAccAciVSANPool_update(t *testing.T) {
	var vsan_pool models.VSANPool
	description := "vsan_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVSANPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVSANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists("aci_vsan_pool.foovsan_pool", &vsan_pool),
					testAccCheckAciVSANPoolAttributes(description, &vsan_pool),
				),
			},
			{
				Config: testAccCheckAciVSANPoolConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVSANPoolExists("aci_vsan_pool.foovsan_pool", &vsan_pool),
					testAccCheckAciVSANPoolAttributes(description, &vsan_pool),
				),
			},
		},
	})
}

func testAccCheckAciVSANPoolConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vsan_pool" "foovsan_pool" {
		description = "%s"
		
		name  = "example"
		
		  alloc_mode  = "static"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciVSANPoolExists(name string, vsan_pool *models.VSANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VSAN Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VSAN Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vsan_poolFound := models.VSANPoolFromContainer(cont)
		if vsan_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VSAN Pool %s not found", rs.Primary.ID)
		}
		*vsan_pool = *vsan_poolFound
		return nil
	}
}

func testAccCheckAciVSANPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_vsan_pool" {
			cont, err := client.Get(rs.Primary.ID)
			vsan_pool := models.VSANPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VSAN Pool %s Still exists", vsan_pool.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciVSANPoolAttributes(description string, vsan_pool *models.VSANPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != vsan_pool.Description {
			return fmt.Errorf("Bad vsan_pool Description %s", vsan_pool.Description)
		}

		if "example" != vsan_pool.Name {
			return fmt.Errorf("Bad vsan_pool name %s", vsan_pool.Name)
		}

		if "static" != vsan_pool.AllocMode {
			return fmt.Errorf("Bad vsan_pool alloc_mode %s", vsan_pool.AllocMode)
		}

		if "example" != vsan_pool.Annotation {
			return fmt.Errorf("Bad vsan_pool annotation %s", vsan_pool.Annotation)
		}

		if "example" != vsan_pool.NameAlias {
			return fmt.Errorf("Bad vsan_pool name_alias %s", vsan_pool.NameAlias)
		}

		return nil
	}
}
