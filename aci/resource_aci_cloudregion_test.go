package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudProvidersRegion_Basic(t *testing.T) {
	var cloud_providers_region models.CloudProvidersRegion
	description := "cloud_providers_region created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudProvidersRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudProvidersRegionConfig_basic(description, "unmanaged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProvidersRegionExists("aci_cloud_providers_region.foocloud_providers_region", &cloud_providers_region),
					testAccCheckAciCloudProvidersRegionAttributes(description, "unmanaged", &cloud_providers_region),
				),
			},
			{
				ResourceName:      "aci_cloud_providers_region",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudProvidersRegion_update(t *testing.T) {
	var cloud_providers_region models.CloudProvidersRegion
	description := "cloud_providers_region created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudProvidersRegionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudProvidersRegionConfig_basic(description, "unmanaged"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProvidersRegionExists("aci_cloud_providers_region.foocloud_providers_region", &cloud_providers_region),
					testAccCheckAciCloudProvidersRegionAttributes(description, "unmanaged", &cloud_providers_region),
				),
			},
			{
				Config: testAccCheckAciCloudProvidersRegionConfig_basic(description, "managed"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudProvidersRegionExists("aci_cloud_providers_region.foocloud_providers_region", &cloud_providers_region),
					testAccCheckAciCloudProvidersRegionAttributes(description, "managed", &cloud_providers_region),
				),
			},
		},
	})
}

func testAccCheckAciCloudProvidersRegionConfig_basic(description, admin_st string) string {
	return fmt.Sprintf(`

	resource "aci_cloud_providers_region" "foocloud_providers_region" {
		cloud_provider_profile_dn = "${aci_cloud_provider_profile.example.id}"
		description               = "%s"
		name                      = "us-east-1"
		admin_st                  = "%s"
		annotation                = "tag_region"
		name_alias                = "default_reg"
	}
	  
	`, description, admin_st)
}

func testAccCheckAciCloudProvidersRegionExists(name string, cloud_providers_region *models.CloudProvidersRegion) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Providers Region %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Providers Region dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_providers_regionFound := models.CloudProvidersRegionFromContainer(cont)
		if cloud_providers_regionFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Providers Region %s not found", rs.Primary.ID)
		}
		*cloud_providers_region = *cloud_providers_regionFound
		return nil
	}
}

func testAccCheckAciCloudProvidersRegionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_providers_region" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_providers_region := models.CloudProvidersRegionFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Providers Region %s Still exists", cloud_providers_region.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudProvidersRegionAttributes(description, admin_st string, cloud_providers_region *models.CloudProvidersRegion) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_providers_region.Description {
			return fmt.Errorf("Bad cloud_providers_region Description %s", cloud_providers_region.Description)
		}

		if "us-east-1" != cloud_providers_region.Name {
			return fmt.Errorf("Bad cloud_providers_region name %s", cloud_providers_region.Name)
		}

		if admin_st != cloud_providers_region.AdminSt {
			return fmt.Errorf("Bad cloud_providers_region admin_st %s", cloud_providers_region.AdminSt)
		}

		if "tag_region" != cloud_providers_region.Annotation {
			return fmt.Errorf("Bad cloud_providers_region annotation %s", cloud_providers_region.Annotation)
		}

		if "default_reg" != cloud_providers_region.NameAlias {
			return fmt.Errorf("Bad cloud_providers_region name_alias %s", cloud_providers_region.NameAlias)
		}

		return nil
	}
}
