package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudProviderProfile_Basic(t *testing.T) {
	var cloud_provider_profile models.CloudProviderProfile
	description := "cloud_provider_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudProviderProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudProviderProfileConfig_basic(description, "tag_provp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProviderProfileExists("aci_cloud_provider_profile.foocloud_provider_profile", &cloud_provider_profile),
					testAccCheckAciCloudProviderProfileAttributes(description, "tag_provp", &cloud_provider_profile),
				),
			},
			{
				ResourceName:      "aci_cloud_provider_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudProviderProfile_update(t *testing.T) {
	var cloud_provider_profile models.CloudProviderProfile
	description := "cloud_provider_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudProviderProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudProviderProfileConfig_basic(description, "tag_provp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProviderProfileExists("aci_cloud_provider_profile.foocloud_provider_profile", &cloud_provider_profile),
					testAccCheckAciCloudProviderProfileAttributes(description, "tag_provp", &cloud_provider_profile),
				),
			},
			{
				Config: testAccCheckAciCloudProviderProfileConfig_basic(description, "provp_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProviderProfileExists("aci_cloud_provider_profile.foocloud_provider_profile", &cloud_provider_profile),
					testAccCheckAciCloudProviderProfileAttributes(description, "provp_update", &cloud_provider_profile),
				),
			},
		},
	})
}

func testAccCheckAciCloudProviderProfileConfig_basic(description, annotation string) string {
	return fmt.Sprintf(`
	resource "aci_cloud_provider_profile" "foocloud_provider_profile" {
		description = "%s"
		vendor      = "aws"
		annotation  = "%s"
	}
	  
	`, description, annotation)
}

func testAccCheckAciCloudProviderProfileExists(name string, cloud_provider_profile *models.CloudProviderProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Provider Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Provider Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_provider_profileFound := models.CloudProviderProfileFromContainer(cont)
		if cloud_provider_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Provider Profile %s not found", rs.Primary.ID)
		}
		*cloud_provider_profile = *cloud_provider_profileFound
		return nil
	}
}

func testAccCheckAciCloudProviderProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_provider_profile" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_provider_profile := models.CloudProviderProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Provider Profile %s Still exists", cloud_provider_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudProviderProfileAttributes(description, annotation string, cloud_provider_profile *models.CloudProviderProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_provider_profile.Description {
			return fmt.Errorf("Bad cloud_provider_profile Description %s", cloud_provider_profile.Description)
		}

		if "aws" != cloud_provider_profile.Vendor {
			return fmt.Errorf("Bad cloud_provider_profile vendor %s", cloud_provider_profile.Vendor)
		}

		if annotation != cloud_provider_profile.Annotation {
			return fmt.Errorf("Bad cloud_provider_profile annotation %s", cloud_provider_profile.Annotation)
		}

		return nil
	}
}
