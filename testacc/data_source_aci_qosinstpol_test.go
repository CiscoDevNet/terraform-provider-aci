package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciQOSInstancePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_qos_instance_policy.test"
	dataSourceName := "data.aci_qos_instance_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	QOSInstancePolicy, err := aci.GetRemoteQOSInstancePolicy(sharedAciClient(), "uni/infra/qosinst-default")
	if err != nil {
		t.Errorf("reading initial config of QOSInstancePolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccQOSInstancePolicyConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_age_timer", resourceName, "etrap_age_timer"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_bw_thresh", resourceName, "etrap_bw_thresh"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_byte_ct", resourceName, "etrap_byte_ct"),
					resource.TestCheckResourceAttrPair(dataSourceName, "etrap_st", resourceName, "etrap_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_flush_interval", resourceName, "fabric_flush_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "fabric_flush_st", resourceName, "fabric_flush_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl.#", resourceName, "ctrl.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uburst_spine_queues", resourceName, "uburst_spine_queues"),
					resource.TestCheckResourceAttrPair(dataSourceName, "uburst_tor_queues", resourceName, "uburst_tor_queues"),
				),
			},
			{
				Config:      CreateAccQOSInstancePolicyDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccQOSInstancePolicyDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreQOSInstancePolicy(QOSInstancePolicy),
			},
		},
	})
}

func CreateAccQOSInstancePolicyConfigDataSource() string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {}

	data "aci_qos_instance_policy" "test" {
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`)
	return resource
}

func CreateAccQOSInstancePolicyDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {}

	data "aci_qos_instance_policy" "test" {
		%s = "%s"
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, key, value)
	return resource
}

func CreateAccQOSInstancePolicyDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing qos_instance_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_qos_instance_policy" "test" {
		%s = "%s"
	}

	data "aci_qos_instance_policy" "test" {
		depends_on = [ aci_qos_instance_policy.test ]
	}
	`, key, value)
	return resource
}
