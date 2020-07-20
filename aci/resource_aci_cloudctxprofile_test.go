package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciCloudContextProfile_Basic(t *testing.T) {
	var cloud_context_profile models.CloudContextProfile
	fv_tenant_name := acctest.RandString(5)
	cloud_ctx_profile_name := acctest.RandString(5)
	description := "cloud_context_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudContextProfileConfig_basic(fv_tenant_name, cloud_ctx_profile_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudContextProfileExists("aci_cloud_context_profile.foocloud_context_profile", &cloud_context_profile),
					testAccCheckAciCloudContextProfileAttributes(cloud_ctx_profile_name, description, &cloud_context_profile),
				),
			},
		},
	})
}

func testAccCheckAciCloudContextProfileConfig_basic(fv_tenant_name, cloud_ctx_profile_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_vrf" "vrf1" {
		tenant_dn = "${aci_tenant.footenant.id}"
		name      = "acc-vrf"
	}

	resource "aci_cloud_context_profile" "foocloud_context_profile" {
		name 		             = "%s"
		description              = "cloud_context_profile created while acceptance testing"
		tenant_dn                = "${aci_tenant.footenant.id}"
		primary_cidr             = "10.230.231.1/16"
		region                   = "us-west-1"
		relation_cloud_rs_to_ctx = "${aci_vrf.vrf1.id}"
	}

	`, fv_tenant_name, cloud_ctx_profile_name)
}

func testAccCheckAciCloudContextProfileExists(name string, cloud_context_profile *models.CloudContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Context Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Context Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_context_profileFound := models.CloudContextProfileFromContainer(cont)
		if cloud_context_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Context Profile %s not found", rs.Primary.ID)
		}
		*cloud_context_profile = *cloud_context_profileFound
		return nil
	}
}

func testAccCheckAciCloudContextProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_context_profile" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_context_profile := models.CloudContextProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Context Profile %s Still exists", cloud_context_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudContextProfileAttributes(cloud_ctx_profile_name, description string, cloud_context_profile *models.CloudContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if cloud_ctx_profile_name != GetMOName(cloud_context_profile.DistinguishedName) {
			return fmt.Errorf("Bad cloud_ctx_profile %s", GetMOName(cloud_context_profile.DistinguishedName))
		}

		if description != cloud_context_profile.Description {
			return fmt.Errorf("Bad cloud_context_profile Description %s", cloud_context_profile.Description)
		}

		if "us-west-1" != cloud_context_profile.Region {
			return fmt.Errorf("Bad cloud_context_profile region %s", cloud_context_profile.Region)
		}

		if "10.230.231.1/16" != cloud_context_profile.PrimaryCIDR {
			return fmt.Errorf("Bad cloud_context_profile Primary CIDR %s", cloud_context_profile.PrimaryCIDR)
		}

		return nil
	}
}
