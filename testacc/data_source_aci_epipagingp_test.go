package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciIPAgingPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_ip_aging_profile.test"
	dataSourceName := "data.aci_endpoint_ip_aging_profile.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciIPAgingPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointIpAgingProfileConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
				),
			},
			{
				Config:      CreateAccEndpointIpAgingProfileDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointIpAgingProfileDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccEndpointIpAgingProfileConfigDataSource() string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
	}

	data "aci_endpoint_ip_aging_profile" "test" {
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`)
	return resource
}

func CreateAccEndpointIpAgingProfileDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {

	}

	data "aci_endpoint_ip_aging_profile" "test" {
		%s = "%s"
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, key, value)
	return resource
}

func CreateAccEndpointIpAgingProfileDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_ip_aging_profile Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_ip_aging_profile" "test" {
		%s = "%s"
	}

	data "aci_endpoint_ip_aging_profile" "test" {
		depends_on = [ aci_endpoint_ip_aging_profile.test ]
	}
	`, key, value)
	return resource
}
