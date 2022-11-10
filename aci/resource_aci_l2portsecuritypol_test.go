package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPortSecurityPolicy_Basic(t *testing.T) {
	var port_security_policy models.PortSecurityPolicy
	description := "port_security_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPortSecurityPolicyConfig_basic(description, "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists("aci_port_security_policy.fooport_security_policy", &port_security_policy),
					testAccCheckAciPortSecurityPolicyAttributes(description, "60", &port_security_policy),
				),
			},
		},
	})
}

func TestAccAciPortSecurityPolicy_update(t *testing.T) {
	var port_security_policy models.PortSecurityPolicy
	description := "port_security_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPortSecurityPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPortSecurityPolicyConfig_basic(description, "60"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists("aci_port_security_policy.fooport_security_policy", &port_security_policy),
					testAccCheckAciPortSecurityPolicyAttributes(description, "60", &port_security_policy),
				),
			},
			{
				Config: testAccCheckAciPortSecurityPolicyConfig_basic(description, "600"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPortSecurityPolicyExists("aci_port_security_policy.fooport_security_policy", &port_security_policy),
					testAccCheckAciPortSecurityPolicyAttributes(description, "600", &port_security_policy),
				),
			},
		},
	})
}

func testAccCheckAciPortSecurityPolicyConfig_basic(description, timeout string) string {
	return fmt.Sprintf(`

	resource "aci_port_security_policy" "fooport_security_policy" {
		description = "%s"
		name        = "demo_port_pol"
		annotation  = "tag_port_pol"
		maximum     = "12"
		name_alias  = "alias_port_pol"
		timeout     = "%s"
		violation   = "protect"
	}
	`, description, timeout)
}

func testAccCheckAciPortSecurityPolicyExists(name string, port_security_policy *models.PortSecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Port Security Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Port Security Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		port_security_policyFound := models.PortSecurityPolicyFromContainer(cont)
		if port_security_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Port Security Policy %s not found", rs.Primary.ID)
		}
		*port_security_policy = *port_security_policyFound
		return nil
	}
}

func testAccCheckAciPortSecurityPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_port_security_policy" {
			cont, err := client.Get(rs.Primary.ID)
			port_security_policy := models.PortSecurityPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Port Security Policy %s Still exists", port_security_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciPortSecurityPolicyAttributes(description, timeout string, port_security_policy *models.PortSecurityPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != port_security_policy.Description {
			return fmt.Errorf("Bad port_security_policy Description %s", port_security_policy.Description)
		}

		if "demo_port_pol" != port_security_policy.Name {
			return fmt.Errorf("Bad port_security_policy name %s", port_security_policy.Name)
		}

		if "tag_port_pol" != port_security_policy.Annotation {
			return fmt.Errorf("Bad port_security_policy annotation %s", port_security_policy.Annotation)
		}

		if "12" != port_security_policy.Maximum {
			return fmt.Errorf("Bad port_security_policy maximum %s", port_security_policy.Maximum)
		}

		if "alias_port_pol" != port_security_policy.NameAlias {
			return fmt.Errorf("Bad port_security_policy name_alias %s", port_security_policy.NameAlias)
		}

		if timeout != port_security_policy.Timeout {
			return fmt.Errorf("Bad port_security_policy timeout %s", port_security_policy.Timeout)
		}

		if "protect" != port_security_policy.Violation {
			return fmt.Errorf("Bad port_security_policy violation %s", port_security_policy.Violation)
		}

		return nil
	}
}
