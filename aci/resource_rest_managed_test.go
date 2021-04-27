package aci

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciRestManaged_basic(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCDPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAciRestManagedConfigTenant(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRestManagedCheckObject("aci_rest_managed.tenant"),
				),
			},
		},
	})
}

func testAccAciRestManagedConfigTenant(name string) string {
	return fmt.Sprintf(`
	resource "aci_rest_managed" "tenant" {
		dn = "uni/tn-%[1]s"
		class_name = "fvTenant"
		content = {
			name = "%[1]s"
			descr = "Test description"
		}
	}
	`, name)
}

func testAccCheckAciRestManagedCheckObject(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Resource aci_rest_managed %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No aci_rest_managed dn attribute was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		className := rs.Primary.Attributes["class_name"]
		dn := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", className, "attributes", "dn").String()))

		if dn != rs.Primary.ID {
			return fmt.Errorf("APIC object %s not found", rs.Primary.ID)
		}

		for key, value := range rs.Primary.Attributes {
			if strings.Contains(key, "content.") && key != "content.%" {
				attr := key[8:]
				v := models.StripQuotes(models.StripSquareBrackets(cont.Search("imdata", className, "attributes", attr).String()))
				if v != value {
					return fmt.Errorf("APIC object %s, expected: %s, got: %s", rs.Primary.ID, value, v)
				}
			}
		}

		return nil
	}
}

func testAccCheckAciRestManagedDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_rest_managed" {
			_, err := client.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Resource aci_rest_managed %s still exists", rs.Primary.ID)
			}

		} else {
			continue
		}
	}

	return nil
}
