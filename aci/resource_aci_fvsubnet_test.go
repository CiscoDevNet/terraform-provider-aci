package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciSubnet_Basic(t *testing.T) {
	var subnet models.Subnet
	description := "subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "unspecified", &subnet),
				),
			},
			{
				ResourceName:      "aci_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciSubnet_update(t *testing.T) {
	var subnet models.Subnet
	description := "subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "unspecified", &subnet),
				),
			},
			{
				Config: testAccCheckAciSubnetConfig_basic(description, "nd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, "nd", &subnet),
				),
			},
		},
	})
}

func testAccCheckAciSubnetConfig_basic(description, Ctrl string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_subnet" {
		name        = "tenant_for_subnet"
		description = "This tenant is created by terraform ACI provider"
	}
	resource "aci_bridge_domain" "bd_for_subnet" {
		tenant_dn = "${aci_tenant.tenant_for_subnet.id}"
		name      = "bd_for_subnet
	}

	resource "aci_subnet" "foosubnet" {
		bridge_domain_dn = "${aci_bridge_domain.bd_for_subnet.id}"
		description      = "%s"
		ip               = "10.0.3.28/27"
		annotation       = "tag_subnet"
		ctrl             = "%s"
		name_alias       = "alias_subnet"
		preferred        = "no"
		scope            = "private"
		virtual          = "yes"
	} 
	`, description, Ctrl)
}

func testAccCheckAciSubnetExists(name string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		subnetFound := models.SubnetFromContainer(cont)
		if subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Subnet %s not found", rs.Primary.ID)
		}
		*subnet = *subnetFound
		return nil
	}
}

func testAccCheckAciSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_subnet" {
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

func testAccCheckAciSubnetAttributes(description, Ctrl string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != subnet.Description {
			return fmt.Errorf("Bad subnet Description %s", subnet.Description)
		}

		if "10.0.3.28/27" != subnet.Ip {
			return fmt.Errorf("Bad subnet ip %s", subnet.Ip)
		}

		if "tag_subnet" != subnet.Annotation {
			return fmt.Errorf("Bad subnet annotation %s", subnet.Annotation)
		}

		if Ctrl != subnet.Ctrl {
			return fmt.Errorf("Bad subnet ctrl %s", subnet.Ctrl)
		}

		if "alias_subnet" != subnet.NameAlias {
			return fmt.Errorf("Bad subnet name_alias %s", subnet.NameAlias)
		}

		if "no" != subnet.Preferred {
			return fmt.Errorf("Bad subnet preferred %s", subnet.Preferred)
		}

		if "private" != subnet.Scope {
			return fmt.Errorf("Bad subnet scope %s", subnet.Scope)
		}

		if "yes" != subnet.Virtual {
			return fmt.Errorf("Bad subnet virtual %s", subnet.Virtual)
		}

		return nil
	}
}
