package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciEndpointLoopProtectionDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_loop_protection.test"
	dataSourceName := "data.aci_endpoint_loop_protection.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	epLoopProtectPolicy, err := aci.GetRemoteEPLoopProtectionPolicy(sharedAciClient(), "uni/infra/epLoopProtectP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointLoopProtectionConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "action", resourceName, "action"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_intvl", resourceName, "loop_detect_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_mult", resourceName, "loop_detect_mult"),
				),
			},
			{
				Config:      CreateAccEndpointLoopProtectionDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointLoopProtectionDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: CreateAccEndpointLoopProtectionInitialConfig(epLoopProtectPolicy),
			},
		},
	})
}

func CreateAccEndpointLoopProtectionConfigDataSource() string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	}

	data "aci_endpoint_loop_protection" "test" {
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`)
	return resource
}

func CreateAccEndpointLoopProtectionDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
	}

	data "aci_endpoint_loop_protection" "test" {
		%s = "%s"
		depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, key, value)
	return resource
}

func CreateAccEndpointLoopProtectionDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_loop_protection Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_loop_protection" "test" {
		%s = "%s"
	}

	data "aci_endpoint_loop_protection" "test" {
			depends_on = [ aci_endpoint_loop_protection.test ]
	}
	`, key, value)
	return resource
}
