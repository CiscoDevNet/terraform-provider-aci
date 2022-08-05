package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudProvidersRegionDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_cloud_providers_region.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudProvidersRegionDSWithoutRequired(cloudProvPName, name, "cloud_provider_profile_dn"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config:      CreateCloudProvidersRegionDSWithoutRequired(cloudProvPName, name, "name"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudProvidersRegionConfigDataSource(cloudProvPName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "cloud_provider_profile_dn", fmt.Sprintf("uni/clouddomp/provp-%s", cloudProvPName)),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "admin_st"),
				),
			},
			{
				Config:      CreateAccCloudProvidersRegionDataSourceUpdate(cloudProvPName, name, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config:      CreateAccCloudProvidersRegionDSWithInvalidName(cloudProvPName, randomParameter),
				ExpectError: regexp.MustCompile(`(.)+ Object may not exists`),
			},
			{
				Config: CreateAccCloudProvidersRegionConfigDataSource(cloudProvPName, name),
			},
		},
	})
}

func CreateAccCloudProvidersRegionConfigDataSource(cloudProvPName, rName string) string {
	fmt.Println("=== STEP  testing cloud_providers_region Data Source with required arguments only")
	resource := fmt.Sprintf(`

	data "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = "uni/clouddomp/provp-%s"
		name  = "%s"
	}
	`, cloudProvPName, rName)
	return resource
}

func CreateCloudProvidersRegionDSWithoutRequired(cloudProvPName, rName, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_providers_region Data Source without ", attrName)
	rBlock := `
	`
	switch attrName {
	case "cloud_provider_profile_dn":
		rBlock += `
	data "aci_cloud_providers_region" "test" {
	#	cloud_provider_profile_dn  = "%s"
		name  = "%s"
	}
		`
	case "name":
		rBlock += `
	data "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = "%s"
	#	name  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, cloudProvPName, rName)
}

func CreateAccCloudProvidersRegionDSWithInvalidName(cloudProvPName, rName string) string {
	fmt.Println("=== STEP  testing cloud_providers_region Data Source with invalid name")
	resource := fmt.Sprintf(`

	data "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = "uni/clouddomp/provp-%s"
		name  = "%s"
	}
	`, cloudProvPName, rName)
	return resource
}

func CreateAccCloudProvidersRegionDataSourceUpdate(cloudProvPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_providers_region Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	data "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = "uni/clouddomp/provp-%s"
		name  = "%s"
		%s = "%s"
	}
	`, cloudProvPName, rName, key, value)
	return resource
}

func CreateAccCloudProvidersRegionDataSourceUpdatedResource(cloudProvPName, rName, key, value string) string {
	fmt.Println("=== STEP  testing cloud_providers_region Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_cloud_provider_profile" "test" {
		name 		= "%s"
	
	}
	
	resource "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = aci_cloud_provider_profile.test.id
		name  = "%s"
		%s = "%s"
	}

	data "aci_cloud_providers_region" "test" {
		cloud_provider_profile_dn  = aci_cloud_provider_profile.test.id
		name  = aci_cloud_providers_region.test.name
		depends_on = [ aci_cloud_providers_region.test ]
	}
	`, cloudProvPName, rName, key, value)
	return resource
}
