package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciWebTokenDataDataSource_Basic(t *testing.T) {
	resourceName := "aci_global_security.test"
	dataSourceName := "data.aci_global_security.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	aaaUserEp, err := aci.GetRemoteUserManagement(sharedAciClient(), "uni/userext")
	if err != nil {
		t.Errorf("reading initial config of aaaUserEp")
	}
	aaaPwdProfile, err := aci.GetRemotePasswordChangeExpirationPolicy(sharedAciClient(), "uni/userext/pwdprofile")
	if err != nil {
		t.Errorf("reading initial config of aaaPwdProfile")
	}
	aaaBlockLoginProfile, err := aci.GetRemoteBlockUserLoginsPolicy(sharedAciClient(), "uni/userext/blockloginp")
	if err != nil {
		t.Errorf("reading initial config of aaaBlockLoginProfile")
	}
	pkiWebTokenData, err := aci.GetRemoteWebTokenData(sharedAciClient(), "uni/userext/pkiext/webtokendata")
	if err != nil {
		t.Errorf("reading initial config of pkiWebTokenData")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccWebTokenDataConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "pwd_strength_check", resourceName, "pwd_strength_check"),
					resource.TestCheckResourceAttrPair(dataSourceName, "change_count", resourceName, "change_count"),
					resource.TestCheckResourceAttrPair(dataSourceName, "change_during_interval", resourceName, "change_during_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "change_interval", resourceName, "change_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "expiration_warn_time", resourceName, "expiration_warn_time"),
					resource.TestCheckResourceAttrPair(dataSourceName, "history_count", resourceName, "history_count"),
					resource.TestCheckResourceAttrPair(dataSourceName, "no_change_interval", resourceName, "no_change_interval"),
					resource.TestCheckResourceAttrPair(dataSourceName, "block_duration", resourceName, "block_duration"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enable_login_block", resourceName, "enable_login_block"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_failed_attempts", resourceName, "max_failed_attempts"),
					resource.TestCheckResourceAttrPair(dataSourceName, "max_failed_attempts_window", resourceName, "max_failed_attempts_window"),
					resource.TestCheckResourceAttrPair(dataSourceName, "maximum_validity_period", resourceName, "maximum_validity_period"),
					resource.TestCheckResourceAttrPair(dataSourceName, "session_record_flags", resourceName, "session_record_flags"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ui_idle_timeout_seconds", resourceName, "ui_idle_timeout_seconds"),
					resource.TestCheckResourceAttrPair(dataSourceName, "webtoken_timeout_seconds", resourceName, "webtoken_timeout_seconds"),
				),
			},
			{
				Config:      CreateAccWebTokenDataDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccWebTokenDataDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreGlobalUser(aaaUserEp, aaaPwdProfile, aaaBlockLoginProfile, pkiWebTokenData),
			},
		},
	})
}

func CreateAccWebTokenDataConfigDataSource() string {
	fmt.Println("=== STEP  testing global_security Data Source")
	resource := fmt.Sprintf(`

	resource "aci_global_security" "test" {

	}

	data "aci_global_security" "test" {

		depends_on = [ aci_global_security.test ]
	}
	`)
	return resource
}

func CreateAccWebTokenDataDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing global_security Data Source with random attribute")
	resource := fmt.Sprintf(`

	resource "aci_global_security" "test" {

	}

	data "aci_global_security" "test" {

		%s = "%s"
		depends_on = [ aci_global_security.test ]
	}
	`, key, value)
	return resource
}

func CreateAccWebTokenDataDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing global_security Data Source with updated resource")
	resource := fmt.Sprintf(`

	resource "aci_global_security" "test" {

		%s = "%s"
	}

	data "aci_global_security" "test" {

		depends_on = [ aci_global_security.test ]
	}
	`, key, value)
	return resource
}
