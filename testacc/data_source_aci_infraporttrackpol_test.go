package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciPortTrackingDataSource_Basic(t *testing.T) {
	resourceName := "aci_port_tracking.test"
	dataSourceName := "data.aci_port_tracking.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	infraPortTrackPol, err := aci.GetRemotePortTracking(sharedAciClient(), "uni/infra/trackEqptFabP-default")
	if err != nil {
		t.Errorf("reading initial config of infraPortTrackPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccPortTrackingConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "delay", resourceName, "delay"),
					resource.TestCheckResourceAttrPair(dataSourceName, "include_apic_ports", resourceName, "include_apic_ports"),
					resource.TestCheckResourceAttrPair(dataSourceName, "minlinks", resourceName, "minlinks"),
				),
			},
			{
				Config:      CreateAccPortTrackingDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccPortTrackingDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restorePortTrackingConfig(infraPortTrackPol),
			},
		},
	})
}

func CreateAccPortTrackingConfigDataSource() string {
	fmt.Println("=== STEP  testing port_tracking Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	}

	data "aci_port_tracking" "test" {
		depends_on = [ aci_port_tracking.test ]
	}
	`)
	return resource
}

func CreateAccPortTrackingDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
	}

	data "aci_port_tracking" "test" {
		%s = "%s"
		depends_on = [ aci_port_tracking.test ]
	}
	`, key, value)
	return resource
}

func CreateAccPortTrackingDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing port_tracking Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_port_tracking" "test" {
		%s = "%s"
	}

	data "aci_port_tracking" "test" {
		depends_on = [ aci_port_tracking.test ]
	}
	`, key, value)
	return resource
}
