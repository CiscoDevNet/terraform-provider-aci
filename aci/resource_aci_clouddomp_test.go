package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciCloudDomainProfile_Basic(t *testing.T) {
	var cloud_domain_profile models.CloudDomainProfile
	description := "cloud_domain_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudDomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudDomainProfileConfig_basic(description, "alias_domp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudDomainProfileExists("aci_cloud_domain_profile.foocloud_domain_profile", &cloud_domain_profile),
					testAccCheckAciCloudDomainProfileAttributes(description, "alias_domp", &cloud_domain_profile),
				),
			},
			{
				ResourceName:      "aci_cloud_domain_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudDomainProfile_update(t *testing.T) {
	var cloud_domain_profile models.CloudDomainProfile
	description := "cloud_domain_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudDomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudDomainProfileConfig_basic(description, "alias_domp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudDomainProfileExists("aci_cloud_domain_profile.foocloud_domain_profile", &cloud_domain_profile),
					testAccCheckAciCloudDomainProfileAttributes(description, "alias_domp", &cloud_domain_profile),
				),
			},
			{
				Config: testAccCheckAciCloudDomainProfileConfig_basic(description, "domp_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudDomainProfileExists("aci_cloud_domain_profile.foocloud_domain_profile", &cloud_domain_profile),
					testAccCheckAciCloudDomainProfileAttributes(description, "domp_update", &cloud_domain_profile),
				),
			},
		},
	})
}

func testAccCheckAciCloudDomainProfileConfig_basic(description, name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_cloud_domain_profile" "foocloud_domain_profile" {
		description = "%s"
		annotation  = "tag_domp"
		name_alias  = "%s"
		site_id     = "0"
	}
	  
	`, description, name_alias)
}

func testAccCheckAciCloudDomainProfileExists(name string, cloud_domain_profile *models.CloudDomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Domain Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Domain Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_domain_profileFound := models.CloudDomainProfileFromContainer(cont)
		if cloud_domain_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Domain Profile %s not found", rs.Primary.ID)
		}
		*cloud_domain_profile = *cloud_domain_profileFound
		return nil
	}
}

func testAccCheckAciCloudDomainProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_domain_profile" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_domain_profile := models.CloudDomainProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Domain Profile %s Still exists", cloud_domain_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudDomainProfileAttributes(description, name_alias string, cloud_domain_profile *models.CloudDomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_domain_profile.Description {
			return fmt.Errorf("Bad cloud_domain_profile Description %s", cloud_domain_profile.Description)
		}

		if "tag_domp" != cloud_domain_profile.Annotation {
			return fmt.Errorf("Bad cloud_domain_profile annotation %s", cloud_domain_profile.Annotation)
		}

		if name_alias != cloud_domain_profile.NameAlias {
			return fmt.Errorf("Bad cloud_domain_profile name_alias %s", cloud_domain_profile.NameAlias)
		}

		if "0" != cloud_domain_profile.SiteId {
			return fmt.Errorf("Bad cloud_domain_profile site_id %s", cloud_domain_profile.SiteId)
		}

		return nil
	}
}
