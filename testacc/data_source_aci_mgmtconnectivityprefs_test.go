package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciMgmtconnectivitypreferenceDataSource_Basic(t *testing.T) {
	resourceName := "aci_mgmt_preference.test"
	dataSourceName := "data.aci_mgmt_preference.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAciMgmtconnectivitypreferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: CreateAccMgmtPreferenceConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "interface_pref", resourceName, "interface_pref"),
				),
			},
			{
				Config:      CreateAccMgmtPreferenceDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccMgmtPreferenceDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
		},
	})
}

func CreateAccMgmtPreferenceConfigDataSource() string {
	fmt.Println("=== STEP  testing mgmt_preference Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
	}

	data "aci_mgmt_preference" "test" {
	
		depends_on = [ aci_mgmt_preference.test ]
	}
	`)
	return resource
}

func CreateAccMgmtPreferenceDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing mgmt_preference Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
	}

	data "aci_mgmt_preference" "test" {
	
		%s = "%s"
		depends_on = [ aci_mgmt_preference.test ]
	}
	`, key, value)
	return resource
}

func CreateAccMgmtPreferenceDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing mgmt_preference Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_mgmt_preference" "test" {
	
		%s = "%s"
	}

	data "aci_mgmt_preference" "test" {
	
		depends_on = [ aci_mgmt_preference.test ]
	}
	`, key, value)
	return resource
}
