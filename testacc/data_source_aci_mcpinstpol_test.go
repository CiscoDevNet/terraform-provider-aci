package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciMCPInstancePolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_mcp_instance_policy.test"
	dataSourceName := "data.aci_mcp_instance_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	key := acctest.RandString(5)
	mcpInstancePolicy, err := aci.GetRemoteMiscablingProtocolInstancePolicy(sharedAciClient(), "uni/infra/mcpInstP-default")
	if err != nil {
		t.Errorf("reading initial config of MCP Instance Policy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMCPInstancePolicyConfigDataSource(key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "admin_st", resourceName, "admin_st"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ctrl", resourceName, "ctrl"),
					resource.TestCheckResourceAttrPair(dataSourceName, "init_delay_time", resourceName, "init_delay_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_detect_mult", resourceName, "loop_detect_mult"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_protect_act", resourceName, "loop_protect_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "loop_protect_act", resourceName, "loop_protect_act"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tx_freq", resourceName, "tx_freq"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tx_freq_msec", resourceName, "tx_freq_msec"),
				),
			},
			{
				Config:      CreateAccMCPInstancePolicyDataSourceUpdate(key, randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMCPInstancePolicyDataSourceUpdatedResource(key, "annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: CreateAccMCPInstancePolicyInitialConfig(key, mcpInstancePolicy),
			},
		},
	})
}

func CreateAccMCPInstancePolicyConfigDataSource(rName string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName)
	return resource
}

func CreateAccMCPInstancePolicyDataSourceUpdate(rName, key, value string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
		%s = "%s"
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}

func CreateAccMCPInstancePolicyDataSourceUpdatedResource(rName, key, value string) string {
	fmt.Println("=== STEP  testing mcp_instance_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_mcp_instance_policy" "test" {
	
		key  = "%s"
		%s = "%s"
	}

	data "aci_mcp_instance_policy" "test" {
	
		depends_on = [ aci_mcp_instance_policy.test ]
	}
	`, rName, key, value)
	return resource
}
