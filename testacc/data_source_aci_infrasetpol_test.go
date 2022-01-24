package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciFabricWideSettingsPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_fabric_wide_settings.test"
	dataSourceName := "data.aci_fabric_wide_settings.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	fabricWideSettings, err := aci.GetRemoteFabricWideSettingsPolicy(sharedAciClient(), "uni/infra/settings")
	if err != nil {
		t.Errorf("reading initial config of fabricWideSettings")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccFabricWideSettingsPolicyConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "disable_ep_dampening", resourceName, "disable_ep_dampening"),
					resource.TestCheckResourceAttrPair(dataSourceName, "domain_validation", resourceName, "domain_validation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_mo_streaming", resourceName, "enable_mo_streaming"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_remote_leaf_direct", resourceName, "enable_remote_leaf_direct"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enforce_subnet_check", resourceName, "enforce_subnet_check"),
					resource.TestCheckResourceAttrPair(dataSourceName, "opflexp_authenticate_clients", resourceName, "opflexp_authenticate_clients"),
					resource.TestCheckResourceAttrPair(dataSourceName, "opflexp_use_ssl", resourceName, "opflexp_use_ssl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "reallocate_gipo", resourceName, "reallocate_gipo"),
					resource.TestCheckResourceAttrPair(dataSourceName, "restrict_infra_vlan_traffic", resourceName, "restrict_infra_vlan_traffic"),
					resource.TestCheckResourceAttrPair(dataSourceName, "unicast_xr_ep_learn_disable", resourceName, "unicast_xr_ep_learn_disable"),
					resource.TestCheckResourceAttrPair(dataSourceName, "validate_overlapping_vlans", resourceName, "validate_overlapping_vlans"),
				),
			},
			{
				Config:      CreateAccFabricWideSettingsPolicyDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccFabricWideSettingsPolicyDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: RestoreFabricWideSetting(fabricWideSettings),
			},
		},
	})
}

func CreateAccFabricWideSettingsPolicyConfigDataSource() string {
	fmt.Println("=== STEP  testing fabric_wide_settings_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {

	}

	data "aci_fabric_wide_settings" "test" {

		depends_on = [ aci_fabric_wide_settings.test ]
	}
	`)
	return resource
}

func CreateAccFabricWideSettingsPolicyDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing fabric_wide_settings_policy Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {

	}

	data "aci_fabric_wide_settings" "test" {

		%s = "%s"
		depends_on = [ aci_fabric_wide_settings.test ]
	}
	`, key, value)
	return resource
}

func CreateAccFabricWideSettingsPolicyDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing fabric_wide_settings_policy Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_fabric_wide_settings" "test" {

		%s = "%s"
	}

	data "aci_fabric_wide_settings" "test" {

		depends_on = [ aci_fabric_wide_settings.test ]
	}
	`, key, value)
	return resource
}
