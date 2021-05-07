package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outRouteTagPolicy_Basic(t *testing.T) {
	var l3out_route_tag_policy models.L3outRouteTagPolicy
	description := "route_tag_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outRouteTagPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists("aci_l3out_route_tag_policy.fool3out_route_tag_policy", &l3out_route_tag_policy),
					testAccCheckAciL3outRouteTagPolicyAttributes(description, &l3out_route_tag_policy),
				),
			},
		},
	})
}

func TestAccAciL3outRouteTagPolicy_update(t *testing.T) {
	var l3out_route_tag_policy models.L3outRouteTagPolicy
	description := "route_tag_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outRouteTagPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outRouteTagPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists("aci_l3out_route_tag_policy.fool3out_route_tag_policy", &l3out_route_tag_policy),
					testAccCheckAciL3outRouteTagPolicyAttributes(description, &l3out_route_tag_policy),
				),
			},
			{
				Config: testAccCheckAciL3outRouteTagPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outRouteTagPolicyExists("aci_l3out_route_tag_policy.fool3out_route_tag_policy", &l3out_route_tag_policy),
					testAccCheckAciL3outRouteTagPolicyAttributes(description, &l3out_route_tag_policy),
				),
			},
		},
	})
}

func testAccCheckAciL3outRouteTagPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_route_tag_policy" "fool3out_route_tag_policy" {
		tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		name_alias  = "example"
  		tag  = "1"
	}
	`, description)
}

func testAccCheckAciL3outRouteTagPolicyExists(name string, l3out_route_tag_policy *models.L3outRouteTagPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3out Route Tag Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3out Route Tag Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_route_tag_policyFound := models.L3outRouteTagPolicyFromContainer(cont)
		if l3out_route_tag_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3out Route Tag Policy %s not found", rs.Primary.ID)
		}
		*l3out_route_tag_policy = *l3out_route_tag_policyFound
		return nil
	}
}

func testAccCheckAciL3outRouteTagPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_route_tag_policy" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_route_tag_policy := models.L3outRouteTagPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3out Route Tag Policy %s Still exists", l3out_route_tag_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outRouteTagPolicyAttributes(description string, l3out_route_tag_policy *models.L3outRouteTagPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_route_tag_policy.Description {
			return fmt.Errorf("Bad l3out_route_tag_policy Description %s", l3out_route_tag_policy.Description)
		}

		if "example" != l3out_route_tag_policy.Name {
			return fmt.Errorf("Bad l3out_route_tag_policy name %s", l3out_route_tag_policy.Name)
		}

		if "example" != l3out_route_tag_policy.Annotation {
			return fmt.Errorf("Bad l3out_route_tag_policy annotation %s", l3out_route_tag_policy.Annotation)
		}

		if "example" != l3out_route_tag_policy.NameAlias {
			return fmt.Errorf("Bad l3out_route_tag_policy name_alias %s", l3out_route_tag_policy.NameAlias)
		}

		if "1" != l3out_route_tag_policy.Tag {
			return fmt.Errorf("Bad l3out_route_tag_policy tag %s", l3out_route_tag_policy.Tag)
		}

		return nil
	}
}
