package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciCloudProviderProfileDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_cloud_provider_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      CreateCloudProviderProfileDSWithoutRequired(cloudVendor, "vendor"),
				ExpectError: regexp.MustCompile(`Missing required argument`),
			},
			{
				Config: CreateAccCloudProviderProfileConfigDataSource(cloudVendor),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "vendor", cloudVendor),
				),
			},
			{
				Config:      CreateAccCloudProviderProfileDataSourceUpdate(cloudVendor, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCloudProviderProfileConfigDataSource(cloudVendor),
			},
		},
	})
}

func CreateAccCloudProviderProfileConfigDataSource(vendor string) string {
	fmt.Println("=== STEP  testing cloud_provider_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	data "aci_cloud_provider_profile" "test" {
		vendor  = "%s"
	}
	`, vendor)
	return resource
}

func CreateCloudProviderProfileDSWithoutRequired(vendor, attrName string) string {
	fmt.Println("=== STEP  Basic: testing cloud_provider_profile Data Source without ", attrName)
	rBlock := ``
	switch attrName {
	case "vendor":
		rBlock += `
	data "aci_cloud_provider_profile" "test" {
	
	#	vendor  = "%s"
	}
		`
	}
	return fmt.Sprintf(rBlock, vendor)
}

func CreateAccCloudProviderProfileDataSourceUpdate(vendor, key, value string) string {
	fmt.Println("=== STEP  testing cloud_provider_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	data "aci_cloud_provider_profile" "test" {
	
		vendor  = "%s"
		%s = "%s"
	}
	`, vendor, key, value)
	return resource
}
