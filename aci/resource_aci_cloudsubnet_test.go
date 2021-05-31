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
	

	resource "aci_tenant" "example" {
		name       = "test_tenant"
		annotation = "atag"
		name_alias = "alias_tenant"
	  }
	  
	  resource "aci_cloud_provider_profile" "example" {
		# description = "cloud provider profile1"
		vendor      = "aws"
		annotation  = "tag_aws_prof1"
	  }
	  
	  resource "aci_vrf" "example" {
		tenant_dn              = aci_tenant.example.id
		name                   = "demo_vrf"
		annotation             = "tag_vrf"
		bd_enforced_enable     = "no"
		ip_data_plane_learning = "enabled"
		knw_mcast_act          = "permit"
		name_alias             = "alias_vrf"
		pc_enf_dir             = "egress"
		pc_enf_pref            = "unenforced"
	  }
	  
	  resource "aci_cloud_context_profile" "example" {
		name                     = "s"
		# description              = "cloud_context_profile created while acceptance testing"
		tenant_dn                = aci_tenant.example.id
		primary_cidr             = "10.235.231.1/24"
		region                   = "us-east-1"
		cloud_vendor             = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.example.id
	  }
	  resource "aci_cloud_cidr_pool" "example" {
		cloud_context_profile_dn = aci_cloud_context_profile.example.id
		# description              = "cloud CIDR"
		addr                     = "10.230.0.0/16"
		annotation               = "tag_cidr"
		name_alias               = "s"
		primary                  = "no"
	  }
	  
	  data "aci_cloud_providers_region" "region_aws" {
		cloud_provider_profile_dn = aci_cloud_provider_profile.example.id
		name                      = "us-east-1"
	  }
	  data "aci_cloud_availability_zone" "az_us_east_1_aws" {
		cloud_providers_region_dn = data.aci_cloud_providers_region.region_aws.id
		name                      = "us-east-1a"
	  }
	  
	  resource "aci_cloud_subnet" "foocloud_subnet" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.example.id
		ip                 = "10.230.0.1/24"
		description        = "%s"
		annotation         = "tag_subnet"
		name_alias         = "alias_subnet"
		scope              = "%s"
		usage              = "user"
		zone               = data.aci_cloud_availability_zone.az_us_east_1_aws.id
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

		if "10.230.0.1/24" != cloud_subnet.Ip {
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
