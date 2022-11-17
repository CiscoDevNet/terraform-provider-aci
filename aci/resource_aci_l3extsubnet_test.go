package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3ExtSubnet_Basic(t *testing.T) {
	var subnet models.L3ExtSubnet
	description := "L3 subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3ExtSubnetConfig_basic(description, "import-rtctrl"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists("aci_l3_ext_subnet.foosubnet", &subnet),
					testAccCheckAciL3ExtSubnetAttributes(description, "import-rtctrl", &subnet),
				),
			},
		},
	})
}

func TestAccAciL3ExtSubnet_update(t *testing.T) {
	var subnet models.L3ExtSubnet
	description := "subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3ExtSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3ExtSubnetConfig_basic(description, "import-rtctrl"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists("aci_l3_ext_subnet.foosubnet", &subnet),
					testAccCheckAciL3ExtSubnetAttributes(description, "import-rtctrl", &subnet),
				),
			},
			{
				Config: testAccCheckAciL3ExtSubnetConfig_basic(description, "export-rtctrl"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3ExtSubnetExists("aci_l3_ext_subnet.foosubnet", &subnet),
					testAccCheckAciL3ExtSubnetAttributes(description, "export-rtctrl", &subnet),
				),
			},
		},
	})
}

func testAccCheckAciL3ExtSubnetConfig_basic(description, aggregate string) string {
	return fmt.Sprintf(`

	resource "aci_l3_ext_subnet" "foosubnet" {
		# external_network_instance_profile_dn = aci_external_network_instance_profile.example.id
		external_network_instance_profile_dn  = "uni/tn-check_tenantnk/out-demo_l3out/instP-demo_inst_prof"
		description                           = "%s"
		ip                                    = "0.0.0.0/0"
		aggregate                             = "%s"
		annotation                            = "tag_ext_subnet"
		name_alias                            = "alias_ext_subnet"
		scope                                 = ["import-security"]
	}
	`, description, aggregate)
}

func testAccCheckAciL3ExtSubnetExists(name string, subnet *models.L3ExtSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Ext Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Ext Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.L3ExtSubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Ext Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound
		return nil
	}
}

func testAccCheckAciL3ExtSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3_ext_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			subnet := models.SubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Subnet %s Still exists", subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3ExtSubnetAttributes(description, aggregate string, subnet *models.L3ExtSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != subnet.Description {
			return fmt.Errorf("Bad subnet Description %s", subnet.Description)
		}

		if "0.0.0.0/0" != subnet.Ip {
			return fmt.Errorf("Bad subnet ip %s", subnet.Ip)
		}

		if aggregate != subnet.Aggregate {
			return fmt.Errorf("Bad subnet aggregate %s", subnet.Aggregate)
		}

		if "tag_ext_subnet" != subnet.Annotation {
			return fmt.Errorf("Bad subnet annotation %s", subnet.Annotation)
		}

		if "alias_ext_subnet" != subnet.NameAlias {
			return fmt.Errorf("Bad subnet name_alias %s", subnet.NameAlias)
		}

		return nil
	}
}
