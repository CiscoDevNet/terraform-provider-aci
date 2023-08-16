package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPowerSupplyRedundancyPolicy_Basic(t *testing.T) {
	var power_supply_redundancy_policy models.PsuInstPol
	psu_inst_pol_name := acctest.RandString(5)
	description := "power_supply_redundancy_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPowerSupplyRedundancyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPowerSupplyRedundancyPolicyConfig_basic(psu_inst_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPowerSupplyRedundancyPolicyExists("aci_power_supply_redundancy_policy.foo_power_supply_redundancy_policy", &power_supply_redundancy_policy),
					testAccCheckAciPowerSupplyRedundancyPolicyAttributes(psu_inst_pol_name, description, &power_supply_redundancy_policy),
				),
			},
		},
	})
}

func TestAccAciPowerSupplyRedundancyPolicy_Update(t *testing.T) {
	var power_supply_redundancy_policy models.PsuInstPol
	psu_inst_pol_name := acctest.RandString(5)
	description := "power_supply_redundancy_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPowerSupplyRedundancyPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPowerSupplyRedundancyPolicyConfig_basic(psu_inst_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPowerSupplyRedundancyPolicyExists("aci_power_supply_redundancy_policy.foo_power_supply_redundancy_policy", &power_supply_redundancy_policy),
					testAccCheckAciPowerSupplyRedundancyPolicyAttributes(psu_inst_pol_name, description, &power_supply_redundancy_policy),
				),
			},
			{
				Config: testAccCheckAciPowerSupplyRedundancyPolicyConfig_basic(psu_inst_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPowerSupplyRedundancyPolicyExists("aci_power_supply_redundancy_policy.foo_power_supply_redundancy_policy", &power_supply_redundancy_policy),
					testAccCheckAciPowerSupplyRedundancyPolicyAttributes(psu_inst_pol_name, description, &power_supply_redundancy_policy),
				),
			},
		},
	})
}

func testAccCheckAciPowerSupplyRedundancyPolicyConfig_basic(psu_inst_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_power_supply_redundancy_policy" "foo_power_supply_redundancy_policy" {
		name 		= "%s"
		description = "power_supply_redundancy_policy created while acceptance testing"

	}

	`, psu_inst_pol_name)
}

func testAccCheckAciPowerSupplyRedundancyPolicyExists(name string, power_supply_redundancy_policy *models.PsuInstPol) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Power Supply Redundancy Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Power Supply Redundancy Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		power_supply_redundancy_policyFound := models.PsuInstPolFromContainer(cont)
		if power_supply_redundancy_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Power Supply Redundancy Policy %s not found", rs.Primary.ID)
		}
		*power_supply_redundancy_policy = *power_supply_redundancy_policyFound
		return nil
	}
}

func testAccCheckAciPowerSupplyRedundancyPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_power_supply_redundancy_policy" {
			cont, err := client.Get(rs.Primary.ID)
			power_supply_redundancy_policy := models.PsuInstPolFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Power Supply Redundancy Policy %s Still exists", power_supply_redundancy_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPowerSupplyRedundancyPolicyAttributes(psu_inst_pol_name, description string, power_supply_redundancy_policy *models.PsuInstPol) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if psu_inst_pol_name != GetMOName(power_supply_redundancy_policy.DistinguishedName) {
			return fmt.Errorf("Bad psuinst_pol %s", GetMOName(power_supply_redundancy_policy.DistinguishedName))
		}

		if description != power_supply_redundancy_policy.Description {
			return fmt.Errorf("Bad power_supply_redundancy_policy Description %s", power_supply_redundancy_policy.Description)
		}
		return nil
	}
}
