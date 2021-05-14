package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFirmwarePolicy_Basic(t *testing.T) {
	var firmware_policy models.FirmwarePolicy
	description := "firmware_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwarePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists("aci_firmware_policy.foofirmware_policy", &firmware_policy),
					testAccCheckAciFirmwarePolicyAttributes(description, &firmware_policy),
				),
			},
			{
				ResourceName:      "aci_firmware_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFirmwarePolicy_update(t *testing.T) {
	var firmware_policy models.FirmwarePolicy
	description := "firmware_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFirmwarePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFirmwarePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists("aci_firmware_policy.foofirmware_policy", &firmware_policy),
					testAccCheckAciFirmwarePolicyAttributes(description, &firmware_policy),
				),
			},
			{
				Config: testAccCheckAciFirmwarePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFirmwarePolicyExists("aci_firmware_policy.foofirmware_policy", &firmware_policy),
					testAccCheckAciFirmwarePolicyAttributes(description, &firmware_policy),
				),
			},
		},
	})
}

func testAccCheckAciFirmwarePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_firmware_policy" "foofirmware_policy" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  effective_on_reboot  = "no"
		  ignore_compat  = "no"
		  internal_label  = "example"
		  name_alias  = "example"
		  version  = "example"
		  version_check_override  = "trigger"
		}
	`, description)
}

func testAccCheckAciFirmwarePolicyExists(name string, firmware_policy *models.FirmwarePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Firmware Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Firmware Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		firmware_policyFound := models.FirmwarePolicyFromContainer(cont)
		if firmware_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Firmware Policy %s not found", rs.Primary.ID)
		}
		*firmware_policy = *firmware_policyFound
		return nil
	}
}

func testAccCheckAciFirmwarePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_firmware_policy" {
			cont, err := client.Get(rs.Primary.ID)
			firmware_policy := models.FirmwarePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Firmware Policy %s Still exists", firmware_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFirmwarePolicyAttributes(description string, firmware_policy *models.FirmwarePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != firmware_policy.Description {
			return fmt.Errorf("Bad firmware_policy Description %s", firmware_policy.Description)
		}

		if "example" != firmware_policy.Name {
			return fmt.Errorf("Bad firmware_policy name %s", firmware_policy.Name)
		}

		if "example" != firmware_policy.Annotation {
			return fmt.Errorf("Bad firmware_policy annotation %s", firmware_policy.Annotation)
		}

		if "no" != firmware_policy.EffectiveOnReboot {
			return fmt.Errorf("Bad firmware_policy effective_on_reboot %s", firmware_policy.EffectiveOnReboot)
		}

		if "no" != firmware_policy.IgnoreCompat {
			return fmt.Errorf("Bad firmware_policy ignore_compat %s", firmware_policy.IgnoreCompat)
		}

		if "example" != firmware_policy.InternalLabel {
			return fmt.Errorf("Bad firmware_policy internal_label %s", firmware_policy.InternalLabel)
		}

		if "example" != firmware_policy.NameAlias {
			return fmt.Errorf("Bad firmware_policy name_alias %s", firmware_policy.NameAlias)
		}

		if "example" != firmware_policy.Version {
			return fmt.Errorf("Bad firmware_policy version %s", firmware_policy.Version)
		}

		if "trigger" != firmware_policy.VersionCheckOverride {
			return fmt.Errorf("Bad firmware_policy version_check_override %s", firmware_policy.VersionCheckOverride)
		}

		return nil
	}
}
