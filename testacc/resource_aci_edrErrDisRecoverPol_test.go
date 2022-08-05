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

func TestAccAciErrorDisableRecovery_Basic(t *testing.T) {
	var error_disable_recovery_default models.ErrorDisabledRecoveryPolicy
	var error_disable_recovery_updated models.ErrorDisabledRecoveryPolicy
	resourceName := "aci_error_disable_recovery.test"
	edrErrDisRecoverPol, err := aci.GetRemoteErrorDisabledRecoveryPolicy(sharedAciClient(), "uni/infra/edrErrDisRecoverPol-default")
	if err != nil {
		t.Errorf("reading initial config of edrErrDisRecoverPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryExists(resourceName, &error_disable_recovery_default),
					resource.TestCheckResourceAttrSet(resourceName, "err_dis_recov_intvl"),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryExists(resourceName, &error_disable_recovery_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_error_disable_recovery"),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "30"),
					testAccCheckAciErrorDisabledRecoveryIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: restoreErrorDisabledRecoveryPolicy(edrErrDisRecoverPol),
			},
		},
	})
}

func TestAccAciErrorDisableRecovery_Update(t *testing.T) {
	var error_disable_recovery_default models.ErrorDisabledRecoveryPolicy
	var error_disable_recovery_updated models.ErrorDisabledRecoveryPolicy
	resourceName := "aci_error_disable_recovery.test"
	edrErrDisRecoverPol, err := aci.GetRemoteErrorDisabledRecoveryPolicy(sharedAciClient(), "uni/infra/edrErrDisRecoverPol-default")
	if err != nil {
		t.Errorf("reading initial config of edrErrDisRecoverPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryExists(resourceName, &error_disable_recovery_default),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryUpdatedAttr("err_dis_recov_intvl", "65535"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryExists(resourceName, &error_disable_recovery_updated),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "65535"),
					testAccCheckAciErrorDisabledRecoveryIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
			{
				Config: CreateAccErrorDisabledRecoveryUpdatedAttr("err_dis_recov_intvl", "32767"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciErrorDisabledRecoveryExists(resourceName, &error_disable_recovery_updated),
					resource.TestCheckResourceAttr(resourceName, "err_dis_recov_intvl", "32767"),
					testAccCheckAciErrorDisabledRecoveryIdEqual(&error_disable_recovery_default, &error_disable_recovery_updated),
				),
			},
			{
				Config: restoreErrorDisabledRecoveryPolicy(edrErrDisRecoverPol),
			},
		},
	})
}

func TestAccAciErrorDisableRecovery_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	edrErrDisRecoverPol, err := aci.GetRemoteErrorDisabledRecoveryPolicy(sharedAciClient(), "uni/infra/edrErrDisRecoverPol-default")
	if err != nil {
		t.Errorf("reading initial config of edrErrDisRecoverPol")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: CreateAccErrorDisabledRecoveryConfig(),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("err_dis_recov_intvl", "29"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("err_dis_recov_intvl", "65536"),
				ExpectError: regexp.MustCompile(`out of range`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr("err_dis_recov_intvl", randomValue),
				ExpectError: regexp.MustCompile(`unknown property value`),
			},
			{
				Config:      CreateAccErrorDisabledRecoveryUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restoreErrorDisabledRecoveryPolicy(edrErrDisRecoverPol),
			},
		},
	})
}

func restoreErrorDisabledRecoveryPolicy(edrErrDisRecoverPol *models.ErrorDisabledRecoveryPolicy) string {
	resource := fmt.Sprintf(`
	resource "aci_error_disable_recovery" "test" {
		annotation = "%s"
		description = "%s"
		name_alias = "%s"
		err_dis_recov_intvl = "%s"
	}
	`, edrErrDisRecoverPol.Annotation, edrErrDisRecoverPol.Description, edrErrDisRecoverPol.NameAlias, edrErrDisRecoverPol.ErrDisRecovIntvl)
	return resource
}

func testAccCheckAciErrorDisabledRecoveryExists(name string, error_disable_recovery *models.ErrorDisabledRecoveryPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Error Disable Recovery %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Error Disable Recovery dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		error_disable_recoveryFound := models.ErrorDisabledRecoveryPolicyFromContainer(cont)
		if error_disable_recoveryFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Error Disable Recovery %s not found", rs.Primary.ID)
		}
		*error_disable_recovery = *error_disable_recoveryFound
		return nil
	}
}

func testAccCheckAciErrorDisabledRecoveryIdEqual(m1, m2 *models.ErrorDisabledRecoveryPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("error_disable_recovery DNs are not equal")
		}
		return nil
	}
}

func CreateAccErrorDisabledRecoveryConfig() string {
	fmt.Println("=== STEP  Testing error_disable_recovery creation")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {}
	`)
	return resource
}

func CreateAccErrorDisabledRecoveryConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Testing error_disable_recovery creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_error_disable_recovery"
        err_dis_recov_intvl = "30"
	}
	`)

	return resource
}

func CreateAccErrorDisabledRecoveryUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  Testing error_disable_recovery attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_error_disable_recovery" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
