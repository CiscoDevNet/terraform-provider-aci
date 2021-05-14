package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciEndPointRetentionPolicy_Basic(t *testing.T) {
	var end_point_retention_policy models.EndPointRetentionPolicy
	description := "end_point_retention_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndPointRetentionPolicyConfig_basic(description, "protocol"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists("aci_end_point_retention_policy.fooend_point_retention_policy", &end_point_retention_policy),
					testAccCheckAciEndPointRetentionPolicyAttributes(description, "protocol", &end_point_retention_policy),
				),
			},
		},
	})
}

func TestAccAciEndPointRetentionPolicy_update(t *testing.T) {
	var end_point_retention_policy models.EndPointRetentionPolicy
	description := "end_point_retention_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndPointRetentionPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndPointRetentionPolicyConfig_basic(description, "protocol"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists("aci_end_point_retention_policy.fooend_point_retention_policy", &end_point_retention_policy),
					testAccCheckAciEndPointRetentionPolicyAttributes(description, "protocol", &end_point_retention_policy),
				),
			},
			{
				Config: testAccCheckAciEndPointRetentionPolicyConfig_basic(description, "rarp-flood"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndPointRetentionPolicyExists("aci_end_point_retention_policy.fooend_point_retention_policy", &end_point_retention_policy),
					testAccCheckAciEndPointRetentionPolicyAttributes(description, "rarp-flood", &end_point_retention_policy),
				),
			},
		},
	})
}

func testAccCheckAciEndPointRetentionPolicyConfig_basic(description, bounce_trig string) string {
	return fmt.Sprintf(`
	resource "aci_tenant" "tenant_for_ret_pol" {
		name        = "tenant_for_ret_pol"
		description = "This tenant is created by terraform ACI provider"
	}
	resource "aci_end_point_retention_policy" "fooend_point_retention_policy" {
		tenant_dn   		= "${aci_tenant.tenant_for_ret_pol.id}"
		description 		= "%s"
		name                = "demo_ret_pol"
		annotation          = "tag_ret_pol"
		bounce_age_intvl    = "630"
		bounce_trig         = "%s"
		hold_intvl          = "6"
		local_ep_age_intvl  = "900"
		move_freq           = "256"
		name_alias          = "alias_demo"
		remote_ep_age_intvl = "300"
	}  
	`, description, bounce_trig)
}

func testAccCheckAciEndPointRetentionPolicyExists(name string, end_point_retention_policy *models.EndPointRetentionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("End Point Retention Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No End Point Retention Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		end_point_retention_policyFound := models.EndPointRetentionPolicyFromContainer(cont)
		if end_point_retention_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("End Point Retention Policy %s not found", rs.Primary.ID)
		}
		*end_point_retention_policy = *end_point_retention_policyFound
		return nil
	}
}

func testAccCheckAciEndPointRetentionPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_end_point_retention_policy" {
			cont, err := client.Get(rs.Primary.ID)
			end_point_retention_policy := models.EndPointRetentionPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("End Point Retention Policy %s Still exists", end_point_retention_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciEndPointRetentionPolicyAttributes(description, bounce_trig string, end_point_retention_policy *models.EndPointRetentionPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != end_point_retention_policy.Description {
			return fmt.Errorf("Bad end_point_retention_policy Description %s", end_point_retention_policy.Description)
		}

		if "demo_ret_pol" != end_point_retention_policy.Name {
			return fmt.Errorf("Bad end_point_retention_policy name %s", end_point_retention_policy.Name)
		}

		if "tag_ret_pol" != end_point_retention_policy.Annotation {
			return fmt.Errorf("Bad end_point_retention_policy annotation %s", end_point_retention_policy.Annotation)
		}

		if "630" != end_point_retention_policy.BounceAgeIntvl {
			return fmt.Errorf("Bad end_point_retention_policy bounce_age_intvl %s", end_point_retention_policy.BounceAgeIntvl)
		}

		if bounce_trig != end_point_retention_policy.BounceTrig {
			return fmt.Errorf("Bad end_point_retention_policy bounce_trig %s", end_point_retention_policy.BounceTrig)
		}

		if "6" != end_point_retention_policy.HoldIntvl {
			return fmt.Errorf("Bad end_point_retention_policy hold_intvl %s", end_point_retention_policy.HoldIntvl)
		}

		if "900" != end_point_retention_policy.LocalEpAgeIntvl {
			return fmt.Errorf("Bad end_point_retention_policy local_ep_age_intvl %s", end_point_retention_policy.LocalEpAgeIntvl)
		}

		if "256" != end_point_retention_policy.MoveFreq {
			return fmt.Errorf("Bad end_point_retention_policy move_freq %s", end_point_retention_policy.MoveFreq)
		}

		if "alias_demo" != end_point_retention_policy.NameAlias {
			return fmt.Errorf("Bad end_point_retention_policy name_alias %s", end_point_retention_policy.NameAlias)
		}

		if "300" != end_point_retention_policy.RemoteEpAgeIntvl {
			return fmt.Errorf("Bad end_point_retention_policy remote_ep_age_intvl %s", end_point_retention_policy.RemoteEpAgeIntvl)
		}

		return nil
	}
}
