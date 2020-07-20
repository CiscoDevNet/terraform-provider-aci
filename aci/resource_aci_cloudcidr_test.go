package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciCloudCIDRPool_Basic(t *testing.T) {
	var cloud_cidr_pool models.CloudCIDRPool
	description := "cloud_cidr_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudCIDRPoolConfig_basic(description, "alias_cidr"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists("aci_cloud_cidr_pool.foocloud_cidr_pool", &cloud_cidr_pool),
					testAccCheckAciCloudCIDRPoolAttributes(description, "alias_cidr", &cloud_cidr_pool),
				),
			},
			{
				ResourceName:      "aci_cloud_cidr_pool",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudCIDRPool_update(t *testing.T) {
	var cloud_cidr_pool models.CloudCIDRPool
	description := "cloud_cidr_pool created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudCIDRPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudCIDRPoolConfig_basic(description, "alias_cidr"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists("aci_cloud_cidr_pool.foocloud_cidr_pool", &cloud_cidr_pool),
					testAccCheckAciCloudCIDRPoolAttributes(description, "alias_cidr", &cloud_cidr_pool),
				),
			},
			{
				Config: testAccCheckAciCloudCIDRPoolConfig_basic(description, "cidr_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudCIDRPoolExists("aci_cloud_cidr_pool.foocloud_cidr_pool", &cloud_cidr_pool),
					testAccCheckAciCloudCIDRPoolAttributes(description, "cidr_update", &cloud_cidr_pool),
				),
			},
		},
	})
}

func testAccCheckAciCloudCIDRPoolConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "ten_for_cidr"
		description = "tenant created while acceptance testing"

	}

	resource "aci_vrf" "vrf1" {
		tenant_dn = "${aci_tenant.footenant.id}"
		name      = "acc-vrf"
	}

	resource "aci_cloud_context_profile" "foocloud_context_profile" {
		name 		             = "ctx_prof_cidr"
		description              = "cloud_context_profile created while acceptance testing"
		tenant_dn                = "${aci_tenant.footenant.id}"
		primary_cidr             = "10.230.231.1/16"
		region                   = "us-west-1"
		relation_cloud_rs_to_ctx = "${aci_vrf.vrf1.id}"
	}

	resource "aci_cloud_cidr_pool" "foocloud_cidr_pool" {
		cloud_context_profile_dn = "${aci_cloud_context_profile.foocloud_context_profile.id}"
		description              = "%s"
		addr                     = "10.0.1.10/28"
		annotation               = "tag_cidr"
		name_alias               = "%s"
		primary                  = "yes"
	}
	  
	`, description, name_alias)
}

func testAccCheckAciCloudCIDRPoolExists(name string, cloud_cidr_pool *models.CloudCIDRPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud CIDR Pool %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud CIDR Pool dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_cidr_poolFound := models.CloudCIDRPoolFromContainer(cont)
		if cloud_cidr_poolFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud CIDR Pool %s not found", rs.Primary.ID)
		}
		*cloud_cidr_pool = *cloud_cidr_poolFound
		return nil
	}
}

func testAccCheckAciCloudCIDRPoolDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_cidr_pool" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_cidr_pool := models.CloudCIDRPoolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud CIDR Pool %s Still exists", cloud_cidr_pool.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudCIDRPoolAttributes(description, name_alias string, cloud_cidr_pool *models.CloudCIDRPool) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_cidr_pool.Description {
			return fmt.Errorf("Bad cloud_cidr_pool Description %s", cloud_cidr_pool.Description)
		}

		if "10.0.1.10/28" != cloud_cidr_pool.Addr {
			return fmt.Errorf("Bad cloud_cidr_pool addr %s", cloud_cidr_pool.Addr)
		}

		if "tag_cidr" != cloud_cidr_pool.Annotation {
			return fmt.Errorf("Bad cloud_cidr_pool annotation %s", cloud_cidr_pool.Annotation)
		}

		if name_alias != cloud_cidr_pool.NameAlias {
			return fmt.Errorf("Bad cloud_cidr_pool name_alias %s", cloud_cidr_pool.NameAlias)
		}

		if "yes" != cloud_cidr_pool.Primary {
			return fmt.Errorf("Bad cloud_cidr_pool primary %s", cloud_cidr_pool.Primary)
		}

		return nil
	}
}
