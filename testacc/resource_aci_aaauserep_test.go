package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/terraform-providers/terraform-provider-aci/aci"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciWebTokenData_Basic(t *testing.T) {
	var global_security_default models.UserManagement
	var global_security_updated models.UserManagement
	resourceName := "aci_global_security.test"
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
				Config: CreateAccWebTokenDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_default),
					resource.TestCheckResourceAttrSet(resourceName, "pwd_strength_check"),
					resource.TestCheckResourceAttrSet(resourceName, "change_count"),
					resource.TestCheckResourceAttrSet(resourceName, "change_during_interval"),
					resource.TestCheckResourceAttrSet(resourceName, "change_interval"),
					resource.TestCheckResourceAttrSet(resourceName, "expiration_warn_time"),
					resource.TestCheckResourceAttrSet(resourceName, "history_count"),
					resource.TestCheckResourceAttrSet(resourceName, "no_change_interval"),
					resource.TestCheckResourceAttrSet(resourceName, "block_duration"),
					resource.TestCheckResourceAttrSet(resourceName, "enable_login_block"),
					resource.TestCheckResourceAttrSet(resourceName, "max_failed_attempts"),
					resource.TestCheckResourceAttrSet(resourceName, "max_failed_attempts_window"),
					resource.TestCheckResourceAttrSet(resourceName, "maximum_validity_period"),
					resource.TestCheckResourceAttrSet(resourceName, "ui_idle_timeout_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "webtoken_timeout_seconds"),
				),
			},
			{
				Config: CreateAccWebTokenDataConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_global_security"),
					resource.TestCheckResourceAttr(resourceName, "pwd_strength_check", "no"),
					resource.TestCheckResourceAttr(resourceName, "change_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "change_during_interval", "disable"),
					resource.TestCheckResourceAttr(resourceName, "change_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "expiration_warn_time", "0"),
					resource.TestCheckResourceAttr(resourceName, "history_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "no_change_interval", "0"),
					resource.TestCheckResourceAttr(resourceName, "block_duration", "1"),
					resource.TestCheckResourceAttr(resourceName, "enable_login_block", "disable"),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts", "1"),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts_window", "1"),
					resource.TestCheckResourceAttr(resourceName, "maximum_validity_period", "4"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.0", "login"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.1", "logout"),
					resource.TestCheckResourceAttr(resourceName, "ui_idle_timeout_seconds", "60"),
					resource.TestCheckResourceAttr(resourceName, "webtoken_timeout_seconds", "300"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restoreGlobalUser(aaaUserEp, aaaPwdProfile, aaaBlockLoginProfile, pkiWebTokenData),
			},
		},
	})
}

func TestAccAciWebTokenData_Update(t *testing.T) {
	var global_security_default models.UserManagement
	var global_security_updated models.UserManagement
	resourceName := "aci_global_security.test"
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
				Config: CreateAccWebTokenDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_default),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("pwd_strength_check", "yes"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "pwd_strength_check", "yes"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("change_count", "10"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "change_count", "10"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("change_count", "5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "change_count", "5"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("change_during_interval", "enable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "change_during_interval", "enable"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("change_interval", "745"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "change_interval", "745"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("change_interval", "372"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "change_interval", "372"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("expiration_warn_time", "15"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "expiration_warn_time", "15"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("expiration_warn_time", "30"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "expiration_warn_time", "30"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("history_count", "15"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "history_count", "15"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("history_count", "7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "history_count", "7"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("no_change_interval", "745"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "no_change_interval", "745"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("no_change_interval", "372"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "no_change_interval", "372"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("block_duration", "1440"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "block_duration", "1440"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("block_duration", "720"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "block_duration", "720"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("enable_login_block", "enable"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "enable_login_block", "enable"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("max_failed_attempts", "15"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts", "15"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("max_failed_attempts", "7"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts", "7"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("max_failed_attempts_window", "720"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts_window", "720"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("max_failed_attempts_window", "360"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "max_failed_attempts_window", "360"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("maximum_validity_period", "24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "maximum_validity_period", "24"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("maximum_validity_period", "14"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "maximum_validity_period", "14"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("ui_idle_timeout_seconds", "32000"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_idle_timeout_seconds", "32000"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("ui_idle_timeout_seconds", "65525"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "ui_idle_timeout_seconds", "65525"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("webtoken_timeout_seconds", "4545"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "webtoken_timeout_seconds", "4545"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttr("webtoken_timeout_seconds", "9600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "webtoken_timeout_seconds", "9600"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttrList("session_record_flags", StringListtoString([]string{"refresh", "logout", "login"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.0", "refresh"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.1", "logout"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.2", "login"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: CreateAccWebTokenDataUpdatedAttrList("session_record_flags", StringListtoString([]string{"login", "logout", "refresh"})),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciWebTokenDataExists(resourceName, &global_security_updated),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.0", "login"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.1", "logout"),
					resource.TestCheckResourceAttr(resourceName, "session_record_flags.2", "refresh"),
					testAccCheckAciWebTokenDataIdEqual(&global_security_default, &global_security_updated),
				),
			},
			{
				Config: restoreGlobalUser(aaaUserEp, aaaPwdProfile, aaaBlockLoginProfile, pkiWebTokenData),
			},
		},
	})
}

func TestAccAciWebTokenData_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
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
				Config: CreateAccWebTokenDataConfig(),
			},

			{
				Config:      CreateAccWebTokenDataUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccWebTokenDataUpdatedAttr("pwd_strength_check", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_count", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_count", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_count", "11"),
				ExpectError: regexp.MustCompile(`is out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_during_interval", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_interval", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("change_interval", "746"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("expiration_warn_time", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("expiration_warn_time", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("expiration_warn_time", "31"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("history_count", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("history_count", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("history_count", "16"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("no_change_interval", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("no_change_interval", "-1"),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("no_change_interval", "746"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("block_duration", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("block_duration", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("block_duration", "1441"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("enable_login_block", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts", "16"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts_window", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts_window", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("max_failed_attempts_window", "721"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("maximum_validity_period", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("maximum_validity_period", "0"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("maximum_validity_period", "721"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("ui_idle_timeout_seconds", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("ui_idle_timeout_seconds", "59"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("ui_idle_timeout_seconds", "65526"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("webtoken_timeout_seconds", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("webtoken_timeout_seconds", "299"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr("webtoken_timeout_seconds", "9601"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttrList("session_record_flags", StringListtoString([]string{"login", "login"})),
				ExpectError: regexp.MustCompile(`duplication is not supported in list`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttrList("session_record_flags", StringListtoString([]string{randomValue})),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},
			{
				Config:      CreateAccWebTokenDataUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restoreGlobalUser(aaaUserEp, aaaPwdProfile, aaaBlockLoginProfile, pkiWebTokenData),
			},
		},
	})
}

func restoreGlobalUser(aaaUserEp *models.UserManagement, aaaPwdProfile *models.PasswordChangeExpirationPolicy, aaaBlockLoginProfile *models.BlockUserLoginsPolicy, pkiWebTokenData *models.WebTokenData) string {
	resource := fmt.Sprintf(`
	resource "aci_global_security" "test" {
		annotation = "%s"
		description = "%s"
		name_alias = "%s"
		pwd_strength_check = "%s"
		change_count = "%s"
		change_during_interval = "%s"
		change_interval = "%s"
		expiration_warn_time = "%s"
		history_count = "%s"
		no_change_interval = "%s"
		block_duration = "%s"
		enable_login_block = "%s"
		max_failed_attempts = "%s"
		max_failed_attempts_window = "%s"
		maximum_validity_period = "%s"
		ui_idle_timeout_seconds = "%s"
		webtoken_timeout_seconds = "%s"
		session_record_flags = %s
	}
	`, aaaUserEp.Annotation, aaaUserEp.Description, aaaUserEp.NameAlias, aaaUserEp.PwdStrengthCheck, aaaPwdProfile.ChangeCount, aaaPwdProfile.ChangeDuringInterval, aaaPwdProfile.ChangeInterval, aaaPwdProfile.ExpirationWarnTime, aaaPwdProfile.HistoryCount, aaaPwdProfile.NoChangeInterval, aaaBlockLoginProfile.BlockDuration, aaaBlockLoginProfile.EnableLoginBlock, aaaBlockLoginProfile.MaxFailedAttempts, aaaBlockLoginProfile.MaxFailedAttemptsWindow, pkiWebTokenData.MaximumValidityPeriod, pkiWebTokenData.UiIdleTimeoutSeconds, pkiWebTokenData.WebtokenTimeoutSeconds, StringListtoString(convertToStringArray(pkiWebTokenData.SessionRecordFlags)))
	return resource
}

func testAccCheckAciWebTokenDataExists(name string, global_security *models.UserManagement) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Web Token Data %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Web Token Data dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		global_securityFound := models.UserManagementFromContainer(cont)
		if global_securityFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Web Token Data %s not found", rs.Primary.ID)
		}
		*global_security = *global_securityFound
		return nil
	}
}

func testAccCheckAciWebTokenDataIdEqual(m1, m2 *models.UserManagement) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("global_security DNs are not equal")
		}
		return nil
	}
}

func CreateAccWebTokenDataConfig() string {
	fmt.Println("=== STEP  testing global_security creation")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {}
	`)
	return resource
}

func CreateAccWebTokenDataConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing global_security creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_global_security"
		pwd_strength_check = "no"
		change_count = "0"
		change_during_interval = "disable"
		change_interval = "0"
		expiration_warn_time = "0"
		history_count = "0"
		no_change_interval = "0"
		block_duration = "1"
		enable_login_block = "disable"
		max_failed_attempts = "1"
		max_failed_attempts_window = "1"
		maximum_validity_period = "4"
		session_record_flags = ["login" , "logout"]
		ui_idle_timeout_seconds = "60"
		webtoken_timeout_seconds = "300"
	}
	`)

	return resource
}

func CreateAccWebTokenDataUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing global_security attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}

func CreateAccWebTokenDataUpdatedAttrList(attribute, value string) string {
	fmt.Printf("=== STEP  testing global_security attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_global_security" "test" {
	
		%s = %s
	}
	`, attribute, value)
	return resource
}
