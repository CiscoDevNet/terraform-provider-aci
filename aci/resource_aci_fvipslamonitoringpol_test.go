package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciIPSLAMonitoringPolicy_Basic(t *testing.T) {
	var ipsla_monitoring_policy models.IPSLAMonitoringPolicy
	fv_tenant_name := acctest.RandString(5)
	fv_ipsla_monitoring_pol_name := acctest.RandString(5)
	description := "ipsla_monitoring_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciIPSLAMonitoringPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciIPSLAMonitoringPolicyConfig_basic(fv_tenant_name, fv_ipsla_monitoring_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPSLAMonitoringPolicyExists("aci_ip_sla_monitoring_policy.foo_ipsla_monitoring_policy", &ipsla_monitoring_policy),
					testAccCheckAciIPSLAMonitoringPolicyAttributes(fv_tenant_name, fv_ipsla_monitoring_pol_name, description, &ipsla_monitoring_policy),
				),
			},
		},
	})
}

func TestAccAciIPSLAMonitoringPolicy_Update(t *testing.T) {
	var ipsla_monitoring_policy models.IPSLAMonitoringPolicy
	fv_tenant_name := acctest.RandString(5)
	fv_ipsla_monitoring_pol_name := acctest.RandString(5)
	description := "ipsla_monitoring_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciIPSLAMonitoringPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciIPSLAMonitoringPolicyConfig_basic(fv_tenant_name, fv_ipsla_monitoring_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPSLAMonitoringPolicyExists("aci_ip_sla_monitoring_policy.foo_ipsla_monitoring_policy", &ipsla_monitoring_policy),
					testAccCheckAciIPSLAMonitoringPolicyAttributes(fv_tenant_name, fv_ipsla_monitoring_pol_name, description, &ipsla_monitoring_policy),
				),
			},
			{
				Config: testAccCheckAciIPSLAMonitoringPolicyConfig_basic(fv_tenant_name, fv_ipsla_monitoring_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIPSLAMonitoringPolicyExists("aci_ip_sla_monitoring_policy.foo_ipsla_monitoring_policy", &ipsla_monitoring_policy),
					testAccCheckAciIPSLAMonitoringPolicyAttributes(fv_tenant_name, fv_ipsla_monitoring_pol_name, description, &ipsla_monitoring_policy),
				),
			},
		},
	})
}

func testAccCheckAciIPSLAMonitoringPolicyConfig_basic(fv_tenant_name, fv_ipsla_monitoring_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_ip_sla_monitoring_policy" "foo_ip_sla_monitoring_policy" {
		name 		= "%s"
		description = "ipsla_monitoring_policy created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	`, fv_tenant_name, fv_ipsla_monitoring_pol_name)
}

func testAccCheckAciIPSLAMonitoringPolicyExists(name string, ipsla_monitoring_policy *models.IPSLAMonitoringPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("IP SLA Monitoring Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP SLA Monitoring Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ipsla_monitoring_policyFound := models.IPSLAMonitoringPolicyFromContainer(cont)
		if ipsla_monitoring_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("IP SLA Monitoring Policy %s not found", rs.Primary.ID)
		}
		*ipsla_monitoring_policy = *ipsla_monitoring_policyFound
		return nil
	}
}

func testAccCheckAciIPSLAMonitoringPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ip_sla_monitoring_policy" {
			cont, err := client.Get(rs.Primary.ID)
			ipsla_monitoring_policy := models.IPSLAMonitoringPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("IP SLA Monitoring Policy %s Still exists", ipsla_monitoring_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciIPSLAMonitoringPolicyAttributes(fv_tenant_name, fv_ipsla_monitoring_pol_name, description string, ipsla_monitoring_policy *models.IPSLAMonitoringPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fv_ipsla_monitoring_pol_name != GetMOName(ipsla_monitoring_policy.DistinguishedName) {
			return fmt.Errorf("Bad fvipsla_monitoring_pol %s", GetMOName(ipsla_monitoring_policy.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(ipsla_monitoring_policy.DistinguishedName, ipsla_monitoring_policy.Rn)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(ipsla_monitoring_policy.DistinguishedName, ipsla_monitoring_policy.Rn)))
		}
		if description != ipsla_monitoring_policy.Description {
			return fmt.Errorf("Bad ip_sla_monitoring_policy Description %s", ipsla_monitoring_policy.Description)
		}
		return nil
	}
}
