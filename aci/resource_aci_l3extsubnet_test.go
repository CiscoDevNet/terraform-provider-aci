package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
				Config: testAccCheckAciSubnetConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, &subnet),
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
				Config: testAccCheckAciSubnetConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, &subnet),
				),
			},
			{
				Config: testAccCheckAciSubnetConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSubnetExists("aci_subnet.foosubnet", &subnet),
					testAccCheckAciSubnetAttributes(description, &subnet),
				),
			},
		},
	})
}

func testAccCheckAciSubnetConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_subnet" "foosubnet" {
		  external_network_instance_profile_dn  = "${aci_external_network_instance_profile.example.id}"
		description = "%s"
		
		ip  = "example"
		  aggregate  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		  scope  = "example"
		}
	`, description)
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

func testAccCheckAciSubnetAttributes(description string, subnet *models.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != subnet.Description {
			return fmt.Errorf("Bad subnet Description %s", subnet.Description)
		}

		if "example" != subnet.Ip {
			return fmt.Errorf("Bad subnet ip %s", subnet.Ip)
		}

		if "example" != subnet.Aggregate {
			return fmt.Errorf("Bad subnet aggregate %s", subnet.Aggregate)
		}

		if "example" != subnet.Annotation {
			return fmt.Errorf("Bad subnet annotation %s", subnet.Annotation)
		}

		if "example" != subnet.NameAlias {
			return fmt.Errorf("Bad subnet name_alias %s", subnet.NameAlias)
		}

		if "example" != subnet.Scope {
			return fmt.Errorf("Bad subnet scope %s", subnet.Scope)
		}

		return nil
	}
}
