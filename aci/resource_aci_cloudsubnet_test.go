package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciCloudSubnet_Basic(t *testing.T) {
	var cloud_subnet models.CloudSubnet
	description := "cloud_subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudSubnetConfig_basic(description, "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists("aci_cloud_subnet.foocloud_subnet", &cloud_subnet),
					testAccCheckAciCloudSubnetAttributes(description, "private", &cloud_subnet),
				),
			},
			{
				ResourceName:      "aci_cloud_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudSubnet_update(t *testing.T) {
	var cloud_subnet models.CloudSubnet
	description := "cloud_subnet created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudSubnetConfig_basic(description, "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists("aci_cloud_subnet.foocloud_subnet", &cloud_subnet),
					testAccCheckAciCloudSubnetAttributes(description, "private", &cloud_subnet),
				),
			},
			{
				Config: testAccCheckAciCloudSubnetConfig_basic(description, "public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudSubnetExists("aci_cloud_subnet.foocloud_subnet", &cloud_subnet),
					testAccCheckAciCloudSubnetAttributes(description, "public", &cloud_subnet),
				),
			},
		},
	})
}

func testAccCheckAciCloudSubnetConfig_basic(description, scope string) string {
	return fmt.Sprintf(`

	resource "aci_cloud_subnet" "foocloud_subnet" {
		cloud_cidr_pool_dn = "${aci_cloud_cidr_pool.example.id}"
		description        = "%s"
		ip                 = "14.12.0.0/28"
		annotation         = "tag_subnet"
		name_alias         = "alias_subnet"
		scope              = "%s"
		usage              = "user"
		zone 			   = "uni/clouddomp/provp-aws/region-us-west-1/zone-us-west-1b"
	}
	  
	`, description, scope)
}

func testAccCheckAciCloudSubnetExists(name string, cloud_subnet *models.CloudSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Subnet %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Subnet dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_subnetFound := models.CloudSubnetFromContainer(cont)
		if cloud_subnetFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Subnet %s not found", rs.Primary.ID)
		}
		*cloud_subnet = *cloud_subnetFound
		return nil
	}
}

func testAccCheckAciCloudSubnetDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_subnet" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_subnet := models.CloudSubnetFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Subnet %s Still exists", cloud_subnet.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudSubnetAttributes(description, scope string, cloud_subnet *models.CloudSubnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_subnet.Description {
			return fmt.Errorf("Bad cloud_subnet Description %s", cloud_subnet.Description)
		}

		if "14.12.0.0/28" != cloud_subnet.Ip {
			return fmt.Errorf("Bad cloud_subnet ip %s", cloud_subnet.Ip)
		}

		if "tag_subnet" != cloud_subnet.Annotation {
			return fmt.Errorf("Bad cloud_subnet annotation %s", cloud_subnet.Annotation)
		}

		if "alias_subnet" != cloud_subnet.NameAlias {
			return fmt.Errorf("Bad cloud_subnet name_alias %s", cloud_subnet.NameAlias)
		}

		if scope != cloud_subnet.Scope {
			return fmt.Errorf("Bad cloud_subnet scope %s", cloud_subnet.Scope)
		}

		if "user" != cloud_subnet.Usage {
			return fmt.Errorf("Bad cloud_subnet usage %s", cloud_subnet.Usage)
		}

		return nil
	}
}
