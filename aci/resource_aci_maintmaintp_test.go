package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciMaintenancePolicy_Basic(t *testing.T) {
	var maintenance_policy models.MaintenancePolicy
	description := "maintenance_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMaintenancePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists("aci_maintenance_policy.foomaintenance_policy", &maintenance_policy),
					testAccCheckAciMaintenancePolicyAttributes(description, &maintenance_policy),
				),
			},
			{
				ResourceName:      "aci_maintenance_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciMaintenancePolicy_update(t *testing.T) {
	var maintenance_policy models.MaintenancePolicy
	description := "maintenance_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMaintenancePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMaintenancePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists("aci_maintenance_policy.foomaintenance_policy", &maintenance_policy),
					testAccCheckAciMaintenancePolicyAttributes(description, &maintenance_policy),
				),
			},
			{
				Config: testAccCheckAciMaintenancePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMaintenancePolicyExists("aci_maintenance_policy.foomaintenance_policy", &maintenance_policy),
					testAccCheckAciMaintenancePolicyAttributes(description, &maintenance_policy),
				),
			},
		},
	})
}

func testAccCheckAciMaintenancePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_maintenance_policy" "foomaintenance_policy" {
		description = "%s"
		
		name  = "example"
		  admin_st  = "untriggered"
		  annotation  = "example"
		  graceful  = "yes"
		  ignore_compat  = "no"
		  internal_label  = "example"
		  name_alias  = "example"
		  notif_cond  = "notifyNever"
		  run_mode  = "pauseNever"
		  version  = "example"
		  version_check_override  = "untriggered"
		}
	`, description)
}

func testAccCheckAciMaintenancePolicyExists(name string, maintenance_policy *models.MaintenancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Maintenance Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Maintenance Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		maintenance_policyFound := models.MaintenancePolicyFromContainer(cont)
		if maintenance_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Maintenance Policy %s not found", rs.Primary.ID)
		}
		*maintenance_policy = *maintenance_policyFound
		return nil
	}
}

func testAccCheckAciMaintenancePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_maintenance_policy" {
			cont, err := client.Get(rs.Primary.ID)
			maintenance_policy := models.MaintenancePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Maintenance Policy %s Still exists", maintenance_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciMaintenancePolicyAttributes(description string, maintenance_policy *models.MaintenancePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != maintenance_policy.Description {
			return fmt.Errorf("Bad maintenance_policy Description %s", maintenance_policy.Description)
		}

		if "example" != maintenance_policy.Name {
			return fmt.Errorf("Bad maintenance_policy name %s", maintenance_policy.Name)
		}

		if "untriggered" != maintenance_policy.AdminSt {
			return fmt.Errorf("Bad maintenance_policy admin_st %s", maintenance_policy.AdminSt)
		}

		if "example" != maintenance_policy.Annotation {
			return fmt.Errorf("Bad maintenance_policy annotation %s", maintenance_policy.Annotation)
		}

		if "yes" != maintenance_policy.Graceful {
			return fmt.Errorf("Bad maintenance_policy graceful %s", maintenance_policy.Graceful)
		}

		if "no" != maintenance_policy.IgnoreCompat {
			return fmt.Errorf("Bad maintenance_policy ignore_compat %s", maintenance_policy.IgnoreCompat)
		}

		if "example" != maintenance_policy.InternalLabel {
			return fmt.Errorf("Bad maintenance_policy internal_label %s", maintenance_policy.InternalLabel)
		}

		if "example" != maintenance_policy.NameAlias {
			return fmt.Errorf("Bad maintenance_policy name_alias %s", maintenance_policy.NameAlias)
		}

		if "notifyNever" != maintenance_policy.NotifCond {
			return fmt.Errorf("Bad maintenance_policy notif_cond %s", maintenance_policy.NotifCond)
		}

		if "pauseNever" != maintenance_policy.RunMode {
			return fmt.Errorf("Bad maintenance_policy run_mode %s", maintenance_policy.RunMode)
		}

		if "example" != maintenance_policy.Version {
			return fmt.Errorf("Bad maintenance_policy version %s", maintenance_policy.Version)
		}

		if "untriggered" != maintenance_policy.VersionCheckOverride {
			return fmt.Errorf("Bad maintenance_policy version_check_override %s", maintenance_policy.VersionCheckOverride)
		}

		return nil
	}
}
