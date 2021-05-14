package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciHSRPGroupPolicy_Basic(t *testing.T) {
	var hsrp_group_policy models.HSRPGroupPolicy
	description := "hsrp_group_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPGroupPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists("aci_hsrp_group_policy.foohsrp_group_policy", &hsrp_group_policy),
					testAccCheckAciHSRPGroupPolicyAttributes(description, &hsrp_group_policy),
				),
			},
		},
	})
}

func TestAccAciHSRPGroupPolicy_update(t *testing.T) {
	var hsrp_group_policy models.HSRPGroupPolicy
	description := "hsrp_group_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciHSRPGroupPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciHSRPGroupPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists("aci_hsrp_group_policy.foohsrp_group_policy", &hsrp_group_policy),
					testAccCheckAciHSRPGroupPolicyAttributes(description, &hsrp_group_policy),
				),
			},
			{
				Config: testAccCheckAciHSRPGroupPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciHSRPGroupPolicyExists("aci_hsrp_group_policy.foohsrp_group_policy", &hsrp_group_policy),
					testAccCheckAciHSRPGroupPolicyAttributes(description, &hsrp_group_policy),
				),
			},
		},
	})
}

func testAccCheckAciHSRPGroupPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_hsrp_group_policy" "foohsrp_group_policy" {
		tenant_dn  = "${aci_tenant.example.id}"
		name  = "example"
		annotation  = "example"
		description = "%s"
		ctrl = "preempt"
		hello_intvl  = "3000"
		hold_intvl  = "10000"
		name_alias  = "example"
		preempt_delay_min  = "60"
		preempt_delay_reload  = "60"
		preempt_delay_sync  = "60"
		prio  = "100"
		timeout  = "60"
		hsrp_group_policy_type = "md5"
	}
	`, description)
}

func testAccCheckAciHSRPGroupPolicyExists(name string, hsrp_group_policy *models.HSRPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("HSRP Group Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No HSRP Group Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		hsrp_group_policyFound := models.HSRPGroupPolicyFromContainer(cont)
		if hsrp_group_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("HSRP Group Policy %s not found", rs.Primary.ID)
		}
		*hsrp_group_policy = *hsrp_group_policyFound
		return nil
	}
}

func testAccCheckAciHSRPGroupPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_hsrp_group_policy" {
			cont, err := client.Get(rs.Primary.ID)
			hsrp_group_policy := models.HSRPGroupPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("HSRP Group Policy %s Still exists", hsrp_group_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciHSRPGroupPolicyAttributes(description string, hsrp_group_policy *models.HSRPGroupPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != hsrp_group_policy.Description {
			return fmt.Errorf("Bad hsrp_group_policy Description %s", hsrp_group_policy.Description)
		}

		if "example" != hsrp_group_policy.Name {
			return fmt.Errorf("Bad hsrp_group_policy name %s", hsrp_group_policy.Name)
		}

		if "example" != hsrp_group_policy.Annotation {
			return fmt.Errorf("Bad hsrp_group_policy annotation %s", hsrp_group_policy.Annotation)
		}

		if "preempt" != hsrp_group_policy.Ctrl {
			return fmt.Errorf("Bad hsrp_group_policy ctrl %s", hsrp_group_policy.Ctrl)
		}

		if "3000" != hsrp_group_policy.HelloIntvl {
			return fmt.Errorf("Bad hsrp_group_policy hello_intvl %s", hsrp_group_policy.HelloIntvl)
		}

		if "10000" != hsrp_group_policy.HoldIntvl {
			return fmt.Errorf("Bad hsrp_group_policy hold_intvl %s", hsrp_group_policy.HoldIntvl)
		}

		if "example" != hsrp_group_policy.NameAlias {
			return fmt.Errorf("Bad hsrp_group_policy name_alias %s", hsrp_group_policy.NameAlias)
		}

		if "60" != hsrp_group_policy.PreemptDelayMin {
			return fmt.Errorf("Bad hsrp_group_policy preempt_delay_min %s", hsrp_group_policy.PreemptDelayMin)
		}

		if "60" != hsrp_group_policy.PreemptDelayReload {
			return fmt.Errorf("Bad hsrp_group_policy preempt_delay_reload %s", hsrp_group_policy.PreemptDelayReload)
		}

		if "60" != hsrp_group_policy.PreemptDelaySync {
			return fmt.Errorf("Bad hsrp_group_policy preempt_delay_sync %s", hsrp_group_policy.PreemptDelaySync)
		}

		if "100" != hsrp_group_policy.Prio {
			return fmt.Errorf("Bad hsrp_group_policy prio %s", hsrp_group_policy.Prio)
		}

		if "60" != hsrp_group_policy.Timeout {
			return fmt.Errorf("Bad hsrp_group_policy timeout %s", hsrp_group_policy.Timeout)
		}

		if "md5" != hsrp_group_policy.HSRPGroupPolicy_type {
			return fmt.Errorf("Bad hsrp_group_policy hsrp_group_policy_type %s", hsrp_group_policy.HSRPGroupPolicy_type)
		}

		return nil
	}
}
