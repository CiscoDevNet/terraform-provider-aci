package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciEndpointControlsDataSource_Basic(t *testing.T) {
	resourceName := "aci_endpoint_controls.test"
	dataSourceName := "data.aci_endpoint_controls.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	epControls, err := aci.GetRemoteEndpointControlPolicy(sharedAciClient(), "uni/infra/epCtrlP-default")
	if err != nil {
		t.Errorf("reading initial config of EP Loop Protection Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccEndpointControlsConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hold_intvl", resourceName, "hold_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rogue_ep_detect_intvl", resourceName, "rogue_ep_detect_intvl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "rogue_ep_detect_mult", resourceName, "rogue_ep_detect_mult"),
				),
			},
			{
				Config:      CreateAccEndpointControlsDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccEndpointControlsDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: CreateAccEndpointControlsInitialConfig(epControls),
			},
		},
	})
}

func CreateAccEndpointControlsConfigDataSource() string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	}

	data "aci_endpoint_controls" "test" {
		depends_on = [ aci_endpoint_controls.test ]
	}
	`)
	return resource
}

func CreateAccEndpointControlsDSWithInvalidName() string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	}

	data "aci_endpoint_controls" "test" {
		depends_on = [ aci_endpoint_controls.test ]
	}
	`)
	return resource
}

func CreateAccEndpointControlsDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
	}

	data "aci_endpoint_controls" "test" {
		%s = "%s"
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, key, value)
	return resource
}

func CreateAccEndpointControlsDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing endpoint_controls Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_endpoint_controls" "test" {
		%s = "%s"
	}

	data "aci_endpoint_controls" "test" {
		depends_on = [ aci_endpoint_controls.test ]
	}
	`, key, value)
	return resource
}
