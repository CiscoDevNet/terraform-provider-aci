package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciErrorDisableRecoveryDataSource_Basic(t *testing.T) {
	resourceName := "aci_error_disable_recovery.test"
	dataSourceName := "data.aci_error_disable_recovery.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	edrErrDisRecoverPol, err := aci.GetRemoteErrorDisabledRecoveryPolicy(sharedAciClient(), "uni/infra/edrErrDisRecoverPol-default")
	if err != nil {
		t.Errorf("reading initial config of edrErrDisRecoverPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "err_dis_recov_intvl", resourceName, "err_dis_recov_intvl"),
				),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccErrorDisabledRecoveryDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreErrorDisabledRecoveryPolicy(edrErrDisRecoverPol),
			},
		},
	})
}

func CreateAccErrorDisabledRecoveryConfigDataSource() string {
	fmt.Println("=== STEP  Testing error_disable_recovery Data Source")
	resource := fmt.Sprintf(`

	resource "aci_error_disable_recovery" "test" {

	}

	data "aci_error_disable_recovery" "test" {

		depends_on = [ aci_error_disable_recovery.test ]
	}
	`)
	return resource
}

func CreateAccErrorDisabledRecoveryDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  Testing error_disable_recovery Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_error_disable_recovery" "test" {

	}

	data "aci_error_disable_recovery" "test" {

		%s = "%s"
		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, key, value)
	return resource
}

func CreateAccErrorDisabledRecoveryDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  Testing error_disable_recovery Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_error_disable_recovery" "test" {

		%s = "%s"
	}

	data "aci_error_disable_recovery" "test" {

		depends_on = [ aci_error_disable_recovery.test ]
	}
	`, key, value)
	return resource
}
