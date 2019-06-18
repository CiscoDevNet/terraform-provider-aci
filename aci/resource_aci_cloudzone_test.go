package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAciCloudAvailabilityZone_Basic(t *testing.T) {
	var cloud_availability_zone models.CloudAvailabilityZone
	description := "cloud_availability_zone created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudAvailabilityZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudAvailabilityZoneConfig_basic(description, "tag_zone"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAvailabilityZoneExists("aci_cloud_availability_zone.foocloud_availability_zone", &cloud_availability_zone),
					testAccCheckAciCloudAvailabilityZoneAttributes(description, "tag_zone", &cloud_availability_zone),
				),
			},
			{
				ResourceName:      "aci_cloud_availability_zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciCloudAvailabilityZone_update(t *testing.T) {
	var cloud_availability_zone models.CloudAvailabilityZone
	description := "cloud_availability_zone created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciCloudAvailabilityZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciCloudAvailabilityZoneConfig_basic(description, "tag_zone"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAvailabilityZoneExists("aci_cloud_availability_zone.foocloud_availability_zone", &cloud_availability_zone),
					testAccCheckAciCloudAvailabilityZoneAttributes(description, "tag_zone", &cloud_availability_zone),
				),
			},
			{
				Config: testAccCheckAciCloudAvailabilityZoneConfig_basic(description, "zone_update"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciCloudAvailabilityZoneExists("aci_cloud_availability_zone.foocloud_availability_zone", &cloud_availability_zone),
					testAccCheckAciCloudAvailabilityZoneAttributes(description, "zone_update", &cloud_availability_zone),
				),
			},
		},
	})
}

func testAccCheckAciCloudAvailabilityZoneConfig_basic(description, annotation string) string {
	return fmt.Sprintf(`

	resource "aci_cloud_availability_zone" "foocloud_availability_zone" {
		cloud_providers_region_dn = "${aci_cloud_providers_region.example.id}"
		description               = "%s"
		name                      = "us-east-1a"
		annotation                = "%s"
		name_alias                = "alias_zone"
	}
	  
	`, description, annotation)
}

func testAccCheckAciCloudAvailabilityZoneExists(name string, cloud_availability_zone *models.CloudAvailabilityZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Cloud Availability Zone %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cloud Availability Zone dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		cloud_availability_zoneFound := models.CloudAvailabilityZoneFromContainer(cont)
		if cloud_availability_zoneFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Cloud Availability Zone %s not found", rs.Primary.ID)
		}
		*cloud_availability_zone = *cloud_availability_zoneFound
		return nil
	}
}

func testAccCheckAciCloudAvailabilityZoneDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_cloud_availability_zone" {
			cont, err := client.Get(rs.Primary.ID)
			cloud_availability_zone := models.CloudAvailabilityZoneFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Cloud Availability Zone %s Still exists", cloud_availability_zone.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciCloudAvailabilityZoneAttributes(description, annotation string, cloud_availability_zone *models.CloudAvailabilityZone) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != cloud_availability_zone.Description {
			return fmt.Errorf("Bad cloud_availability_zone Description %s", cloud_availability_zone.Description)
		}

		if "us-east-1a" != cloud_availability_zone.Name {
			return fmt.Errorf("Bad cloud_availability_zone name %s", cloud_availability_zone.Name)
		}

		if annotation != cloud_availability_zone.Annotation {
			return fmt.Errorf("Bad cloud_availability_zone annotation %s", cloud_availability_zone.Annotation)
		}

		if "alias_zone" != cloud_availability_zone.NameAlias {
			return fmt.Errorf("Bad cloud_availability_zone name_alias %s", cloud_availability_zone.NameAlias)
		}

		return nil
	}
}
