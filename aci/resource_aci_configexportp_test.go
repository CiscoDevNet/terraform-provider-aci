package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciConfigurationExportPolicy_Basic(t *testing.T) {
	var configuration_export_policy models.ConfigurationExportPolicy
	description := "configuration_export_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConfigurationExportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists("aci_configuration_export_policy.fooconfiguration_export_policy", &configuration_export_policy),
					testAccCheckAciConfigurationExportPolicyAttributes(description, &configuration_export_policy),
				),
			},
		},
	})
}

func TestAccAciConfigurationExportPolicy_update(t *testing.T) {
	var configuration_export_policy models.ConfigurationExportPolicy
	description := "configuration_export_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciConfigurationExportPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciConfigurationExportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists("aci_configuration_export_policy.fooconfiguration_export_policy", &configuration_export_policy),
					testAccCheckAciConfigurationExportPolicyAttributes(description, &configuration_export_policy),
				),
			},
			{
				Config: testAccCheckAciConfigurationExportPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciConfigurationExportPolicyExists("aci_configuration_export_policy.fooconfiguration_export_policy", &configuration_export_policy),
					testAccCheckAciConfigurationExportPolicyAttributes(description, &configuration_export_policy),
				),
			},
		},
	})
}

func testAccCheckAciConfigurationExportPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_configuration_export_policy" "fooconfiguration_export_policy" {
		name                  = "example"
		description           = "%s"
		admin_st              = "untriggered"
		annotation            = "example"
		format                = "json"
		include_secure_fields = "yes"
		max_snapshot_count    = "10"
		name_alias            = "example"
		snapshot              = "yes"
		target_dn             = "uni/tn-crest_test_kishan_tenant"
		}
	`, description)
}

func testAccCheckAciConfigurationExportPolicyExists(name string, configuration_export_policy *models.ConfigurationExportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Configuration Export Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Configuration Export Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		configuration_export_policyFound := models.ConfigurationExportPolicyFromContainer(cont)
		if configuration_export_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Configuration Export Policy %s not found", rs.Primary.ID)
		}
		*configuration_export_policy = *configuration_export_policyFound
		return nil
	}
}

func testAccCheckAciConfigurationExportPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_configuration_export_policy" {
			cont, err := client.Get(rs.Primary.ID)
			configuration_export_policy := models.ConfigurationExportPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Configuration Export Policy %s Still exists", configuration_export_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciConfigurationExportPolicyAttributes(description string, configuration_export_policy *models.ConfigurationExportPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != configuration_export_policy.Description {
			return fmt.Errorf("Bad configuration_export_policy Description %s", configuration_export_policy.Description)
		}

		if "example" != configuration_export_policy.Name {
			return fmt.Errorf("Bad configuration_export_policy name %s", configuration_export_policy.Name)
		}

		if "untriggered" != configuration_export_policy.AdminSt {
			return fmt.Errorf("Bad configuration_export_policy admin_st %s", configuration_export_policy.AdminSt)
		}

		if "example" != configuration_export_policy.Annotation {
			return fmt.Errorf("Bad configuration_export_policy annotation %s", configuration_export_policy.Annotation)
		}

		if "json" != configuration_export_policy.Format {
			return fmt.Errorf("Bad configuration_export_policy format %s", configuration_export_policy.Format)
		}

		if "no" != configuration_export_policy.IncludeSecureFields {
			return fmt.Errorf("Bad configuration_export_policy include_secure_fields %s", configuration_export_policy.IncludeSecureFields)
		}

		if "example" != configuration_export_policy.MaxSnapshotCount {
			return fmt.Errorf("Bad configuration_export_policy max_snapshot_count %s", configuration_export_policy.MaxSnapshotCount)
		}

		if "example" != configuration_export_policy.NameAlias {
			return fmt.Errorf("Bad configuration_export_policy name_alias %s", configuration_export_policy.NameAlias)
		}

		if "yes" != configuration_export_policy.Snapshot {
			return fmt.Errorf("Bad configuration_export_policy snapshot %s", configuration_export_policy.Snapshot)
		}

		if "example" != configuration_export_policy.TargetDn {
			return fmt.Errorf("Bad configuration_export_policy target_dn %s", configuration_export_policy.TargetDn)
		}

		return nil
	}
}
