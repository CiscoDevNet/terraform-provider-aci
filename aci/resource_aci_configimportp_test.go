package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciConfigurationImportPolicy_Basic(t *testing.T) {
	var configuration_import_policy models.ConfigurationImportPolicy
	description := "configuration_import_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConfigurationImportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists("aci_configuration_import_policy.fooconfiguration_import_policy", &configuration_import_policy),
					testAccCheckAciConfigurationImportPolicyAttributes(description, &configuration_import_policy),
				),
			},
			{
				ResourceName:      "aci_configuration_import_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciConfigurationImportPolicy_update(t *testing.T) {
	var configuration_import_policy models.ConfigurationImportPolicy
	description := "configuration_import_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConfigurationImportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConfigurationImportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists("aci_configuration_import_policy.fooconfiguration_import_policy", &configuration_import_policy),
					testAccCheckAciConfigurationImportPolicyAttributes(description, &configuration_import_policy),
				),
			},
			{
				Config: testAccCheckAciConfigurationImportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationImportPolicyExists("aci_configuration_import_policy.fooconfiguration_import_policy", &configuration_import_policy),
					testAccCheckAciConfigurationImportPolicyAttributes(description, &configuration_import_policy),
				),
			},
		},
	})
}

func testAccCheckAciConfigurationImportPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_configuration_import_policy" "fooconfiguration_import_policy" {
		description = "%s"
		
		name  = "example"
		  admin_st  = "triggered"
		  annotation  = "example"
		  fail_on_decrypt_errors  = "no"
		  file_name  = "example"
		  import_mode  = "atomic"
		  import_type  = "merge"
		  name_alias  = "example"
		  snapshot  = "no"
		}
	`, description)
}

func testAccCheckAciConfigurationImportPolicyExists(name string, configuration_import_policy *models.ConfigurationImportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Configuration Import Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Configuration Import Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		configuration_import_policyFound := models.ConfigurationImportPolicyFromContainer(cont)
		if configuration_import_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Configuration Import Policy %s not found", rs.Primary.ID)
		}
		*configuration_import_policy = *configuration_import_policyFound
		return nil
	}
}

func testAccCheckAciConfigurationImportPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_configuration_import_policy" {
			cont, err := client.Get(rs.Primary.ID)
			configuration_import_policy := models.ConfigurationImportPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Configuration Import Policy %s Still exists", configuration_import_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciConfigurationImportPolicyAttributes(description string, configuration_import_policy *models.ConfigurationImportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != configuration_import_policy.Description {
			return fmt.Errorf("Bad configuration_import_policy Description %s", configuration_import_policy.Description)
		}

		if "example" != configuration_import_policy.Name {
			return fmt.Errorf("Bad configuration_import_policy name %s", configuration_import_policy.Name)
		}

		if "triggered" != configuration_import_policy.AdminSt {
			return fmt.Errorf("Bad configuration_import_policy admin_st %s", configuration_import_policy.AdminSt)
		}

		if "example" != configuration_import_policy.Annotation {
			return fmt.Errorf("Bad configuration_import_policy annotation %s", configuration_import_policy.Annotation)
		}

		if "no" != configuration_import_policy.FailOnDecryptErrors {
			return fmt.Errorf("Bad configuration_import_policy fail_on_decrypt_errors %s", configuration_import_policy.FailOnDecryptErrors)
		}

		if "example" != configuration_import_policy.FileName {
			return fmt.Errorf("Bad configuration_import_policy file_name %s", configuration_import_policy.FileName)
		}

		if "atomic" != configuration_import_policy.ImportMode {
			return fmt.Errorf("Bad configuration_import_policy import_mode %s", configuration_import_policy.ImportMode)
		}

		if "merge" != configuration_import_policy.ImportType {
			return fmt.Errorf("Bad configuration_import_policy import_type %s", configuration_import_policy.ImportType)
		}

		if "example" != configuration_import_policy.NameAlias {
			return fmt.Errorf("Bad configuration_import_policy name_alias %s", configuration_import_policy.NameAlias)
		}

		if "no" != configuration_import_policy.Snapshot {
			return fmt.Errorf("Bad configuration_import_policy snapshot %s", configuration_import_policy.Snapshot)
		}

		return nil
	}
}
