package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outBgpExternalPolicy_Basic(t *testing.T) {
	var l3out_bgp_external_policy models.L3outBgpExternalPolicy
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBgpExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBgpExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBgpExternalPolicyExists("aci_external_profile.fooexternal_profile", &l3out_bgp_external_policy),
					testAccCheckAciL3outBgpExternalPolicyAttributes(description, &l3out_bgp_external_policy),
				),
			},
		},
	})
}

func TestAccAciL3outBgpExternalPolicy_update(t *testing.T) {
	var l3out_bgp_external_policy models.L3outBgpExternalPolicy
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outBgpExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outBgpExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBgpExternalPolicyExists("aci_external_profile.fooexternal_profile", &l3out_bgp_external_policy),
					testAccCheckAciL3outBgpExternalPolicyAttributes(description, &l3out_bgp_external_policy),
				),
			},
			{
				Config: testAccCheckAciL3outBgpExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outBgpExternalPolicyExists("aci_external_profile.fooexternal_profile", &l3out_bgp_external_policy),
					testAccCheckAciL3outBgpExternalPolicyAttributes(description, &l3out_bgp_external_policy),
				),
			},
		},
	})
}

func testAccCheckAciL3outBgpExternalPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_bgp_external_policy" "fool3out_bgp_external_policy" {
		l3_outside_dn  = aci_l3_outside.example.id
		description = "%s"
  		annotation  = "example"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciL3outBgpExternalPolicyExists(name string, l3out_bgp_external_policy *models.L3outBgpExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3outBgpExternalPolicy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3outBgpExternalPolicy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_bgp_external_policyFound := models.L3outBgpExternalPolicyFromContainer(cont)
		if l3out_bgp_external_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3outBgpExternalPolicy %s not found", rs.Primary.ID)
		}
		*l3out_bgp_external_policy = *l3out_bgp_external_policyFound
		return nil
	}
}

func testAccCheckAciL3outBgpExternalPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_external_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_bgp_external_policy := models.L3outBgpExternalPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3outBgpExternalPolicy %s Still exists", l3out_bgp_external_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outBgpExternalPolicyAttributes(description string, l3out_bgp_external_policy *models.L3outBgpExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_bgp_external_policy.Description {
			return fmt.Errorf("Bad l3out_bgp_external_policy Description %s", l3out_bgp_external_policy.Description)
		}

		if "example" != l3out_bgp_external_policy.Annotation {
			return fmt.Errorf("Bad l3out_bgp_external_policy annotation %s", l3out_bgp_external_policy.Annotation)
		}

		if "example" != l3out_bgp_external_policy.NameAlias {
			return fmt.Errorf("Bad l3out_bgp_external_policy name_alias %s", l3out_bgp_external_policy.NameAlias)
		}

		return nil
	}
}
