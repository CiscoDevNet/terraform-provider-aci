package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudAvailabilityZoneDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_cloud_availability_zone.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudAvailabilityZoneDSWithoutRequired("cloud_providers_region_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudAvailabilityZoneDSWithoutRequired("name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudAvailabilityZoneConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "cloud_providers_region_dn", cloudProviderRegion),
					resource.TestCheckResourceAttr(dataSourceName, "name", zoneName),
				),
			},
			{
				Config:      CreateAccCloudAvailabilityZoneDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config:      CreateAccCloudAvailabilityZoneDSWithInvalidName(),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
		},
	})
}

func CreateAccCloudAvailabilityZoneConfigDataSource() string {
	fmt.Println("=== STEP  testing cloud_availability_zone Data Source with required arguments only")
	resource := fmt.Sprintf(`

	data "aci_cloud_availability_zone" "test" {
		cloud_providers_region_dn  = "%s"
		name  = "%s"
	}
	`, cloudProviderRegion, zoneName)
	return resource
}

func CreateCloudAvailabilityZoneDSWithoutRequired(attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_availability_zone Data Source without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "cloud_providers_region_dn":
		rBlock += `
	data "aci_cloud_availability_zone" "test" {
	#	cloud_providers_region_dn  = "%s"
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_availability_zone" "test" {
		cloud_providers_region_dn  = "%s"
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, cloudProviderRegion, zoneName)
}

func CreateAccCloudAvailabilityZoneDSWithInvalidName() string {
	fmt.Println("=== STEP  testing cloud_availability_zone Data Source with invalid name")
	resource := fmt.Sprintf(`

	data "aci_cloud_availability_zone" "test" {
		cloud_providers_region_dn  = "%s"
		name  = "%sxyz"
	}
	`, cloudProviderRegion, zoneName)
	return resource
}

func CreateAccCloudAvailabilityZoneDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing cloud_availability_zone Data Source with random attribute")
	resource := fmt.Sprintf(`
	data "aci_cloud_availability_zone" "test" {
		cloud_providers_region_dn  = "%s"
		name  = "%s"
		%s = "%s"
	}
	`, cloudProviderRegion, zoneName, key, value)
	return resource
}
