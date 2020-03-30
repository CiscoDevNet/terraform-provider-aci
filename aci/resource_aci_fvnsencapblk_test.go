package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciRanges_Basic(t *testing.T) {
	var ranges models.Ranges
	description := "ranges created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRangesConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists("aci_ranges.fooranges", &ranges),
					testAccCheckAciRangesAttributes(description, &ranges),
				),
			},
			{
				ResourceName:      "aci_ranges",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciRanges_update(t *testing.T) {
	var ranges models.Ranges
	description := "ranges created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRangesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRangesConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists("aci_ranges.fooranges", &ranges),
					testAccCheckAciRangesAttributes(description, &ranges),
				),
			},
			{
				Config: testAccCheckAciRangesConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRangesExists("aci_ranges.fooranges", &ranges),
					testAccCheckAciRangesAttributes(description, &ranges),
				),
			},
		},
	})
}

func testAccCheckAciRangesConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_ranges" "fooranges" {
		  vlan_pool_dn  = "${aci_vlan_pool.example.id}"
		description = "%s"
		
		_from  = "example"
		
		to  = "example"
		  alloc_mode  = "example"
		  annotation  = "example"
		  from  = "example"
		  name_alias  = "example"
		  role  = "example"
		}
	`, description)
}

func testAccCheckAciRangesExists(name string, ranges *models.Ranges) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Ranges %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Ranges dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		rangesFound := models.RangesFromContainer(cont)
		if rangesFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Ranges %s not found", rs.Primary.ID)
		}
		*ranges = *rangesFound
		return nil
	}
}

func testAccCheckAciRangesDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_ranges" {
			cont, err := client.Get(rs.Primary.ID)
			ranges := models.RangesFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Ranges %s Still exists", ranges.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciRangesAttributes(description string, ranges *models.Ranges) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != ranges.Description {
			return fmt.Errorf("Bad ranges Description %s", ranges.Description)
		}

		if "example" != ranges.From {
			return fmt.Errorf("Bad ranges _from %s", ranges.From)
		}

		if "example" != ranges.To {
			return fmt.Errorf("Bad ranges to %s", ranges.To)
		}

		if "example" != ranges.AllocMode {
			return fmt.Errorf("Bad ranges alloc_mode %s", ranges.AllocMode)
		}

		if "example" != ranges.Annotation {
			return fmt.Errorf("Bad ranges annotation %s", ranges.Annotation)
		}

		if "example" != ranges.From {
			return fmt.Errorf("Bad ranges from %s", ranges.From)
		}

		if "example" != ranges.NameAlias {
			return fmt.Errorf("Bad ranges name_alias %s", ranges.NameAlias)
		}

		if "example" != ranges.Role {
			return fmt.Errorf("Bad ranges role %s", ranges.Role)
		}

		return nil
	}
}
