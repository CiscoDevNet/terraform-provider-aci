package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciCoopPolicyDataSource_Basic(t *testing.T) {
	resourceName := "aci_coop_policy.test"
	dataSourceName := "data.aci_coop_policy.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	coopPolicy, err := aci.GetRemoteCOOPGroupPolicy(sharedAciClient(), "uni/fabric/pol-default")
	if err != nil {
		t.Errorf("reading initial config of coopPolicy")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccCoopPolicyConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
				),
			},
			{
				Config:      CreateAccCoopPolicyDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccCoopPolicyDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: CreateAccCoopPolicyInitialConfig(coopPolicy),
			},
		},
	})
}

func CreateAccCoopPolicyConfigDataSource() string {
	fmt.Println("=== STEP  testing coop_policy Data Source")
	resource := fmt.Sprintf(`
	resource "aci_coop_policy" "test" {
	}

	data "aci_coop_policy" "test" {
		depends_on = [ aci_coop_policy.test ]
	}
	`)
	return resource
}

func CreateAccCoopPolicyDSWithInvalidName() string {
	fmt.Println("=== STEP  testing coop_policy Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	}

	data "aci_coop_policy" "test" {
		depends_on = [ aci_coop_policy.test ]
	}
	`)
	return resource
}

func CreateAccCoopPolicyDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
	}

	data "aci_coop_policy" "test" {
		%s = "%s"
		depends_on = [ aci_coop_policy.test ]
	}
	`, key, value)
	return resource
}

func CreateAccCoopPolicyDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing coop_policy Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_coop_policy" "test" {
		%s = "%s"
	}

	data "aci_coop_policy" "test" {
		depends_on = [ aci_coop_policy.test ]
	}
	`, key, value)
	return resource
}
