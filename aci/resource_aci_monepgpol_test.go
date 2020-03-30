package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciMonitoringPolicy_Basic(t *testing.T) {
	var monitoring_policy models.MonitoringPolicy
	description := "monitoring_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMonitoringPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMonitoringPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMonitoringPolicyExists("aci_monitoring_policy.foomonitoring_policy", &monitoring_policy),
					testAccCheckAciMonitoringPolicyAttributes(description, &monitoring_policy),
				),
			},
			{
				ResourceName:      "aci_monitoring_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciMonitoringPolicy_update(t *testing.T) {
	var monitoring_policy models.MonitoringPolicy
	description := "monitoring_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciMonitoringPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciMonitoringPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMonitoringPolicyExists("aci_monitoring_policy.foomonitoring_policy", &monitoring_policy),
					testAccCheckAciMonitoringPolicyAttributes(description, &monitoring_policy),
				),
			},
			{
				Config: testAccCheckAciMonitoringPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciMonitoringPolicyExists("aci_monitoring_policy.foomonitoring_policy", &monitoring_policy),
					testAccCheckAciMonitoringPolicyAttributes(description, &monitoring_policy),
				),
			},
		},
	})
}

func testAccCheckAciMonitoringPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_monitoring_policy" "foomonitoring_policy" {
		  tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciMonitoringPolicyExists(name string, monitoring_policy *models.MonitoringPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Monitoring Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Monitoring Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		monitoring_policyFound := models.MonitoringPolicyFromContainer(cont)
		if monitoring_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Monitoring Policy %s not found", rs.Primary.ID)
		}
		*monitoring_policy = *monitoring_policyFound
		return nil
	}
}

func testAccCheckAciMonitoringPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_monitoring_policy" {
			cont, err := client.Get(rs.Primary.ID)
			monitoring_policy := models.MonitoringPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Monitoring Policy %s Still exists", monitoring_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciMonitoringPolicyAttributes(description string, monitoring_policy *models.MonitoringPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != monitoring_policy.Description {
			return fmt.Errorf("Bad monitoring_policy Description %s", monitoring_policy.Description)
		}

		if "example" != monitoring_policy.Name {
			return fmt.Errorf("Bad monitoring_policy name %s", monitoring_policy.Name)
		}

		if "example" != monitoring_policy.Annotation {
			return fmt.Errorf("Bad monitoring_policy annotation %s", monitoring_policy.Annotation)
		}

		if "example" != monitoring_policy.NameAlias {
			return fmt.Errorf("Bad monitoring_policy name_alias %s", monitoring_policy.NameAlias)
		}

		return nil
	}
}
