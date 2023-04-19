package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciVPCDomainPolicy_Basic(t *testing.T) {
	var vpc_domain_policy models.VPCDomainPolicy
	description := "vpc_domain_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVPCDomainPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists("aci_vpc_domain_policy.test", &vpc_domain_policy),
					testAccCheckAciVPCDomainPolicyAttributes(description, &vpc_domain_policy),
				),
			},
		},
	})
}

func TestAccAciVPCDomainPolicy_update(t *testing.T) {
	var vpc_domain_policy models.VPCDomainPolicy
	description := "vpc_domain_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVPCDomainPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVPCDomainPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists("aci_vpc_domain_policy.test", &vpc_domain_policy),
					testAccCheckAciVPCDomainPolicyAttributes(description, &vpc_domain_policy),
				),
			},
			{
				Config: testAccCheckAciVPCDomainPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVPCDomainPolicyExists("aci_vpc_domain_policy.test", &vpc_domain_policy),
					testAccCheckAciVPCDomainPolicyAttributes(description, &vpc_domain_policy),
				),
			},
		},
	})
}

func testAccCheckAciVPCDomainPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_vpc_domain_policy" "test" {
		name 		= "test"
		description = "%s"
		dead_intvl = "200"
		annotation = "test_annotation"
		name_alias = "test_alias"
	}
	`, description)
}

func testAccCheckAciVPCDomainPolicyExists(name string, vpc_domain_policy *models.VPCDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VPC Domain Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC Domain Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vpc_domain_policyFound := models.VPCDomainPolicyFromContainer(cont)
		if vpc_domain_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VPC Domain Policy %s not found", rs.Primary.ID)
		}
		*vpc_domain_policy = *vpc_domain_policyFound
		return nil
	}
}

func testAccCheckAciVPCDomainPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vpc_domain_policy" {
			cont, err := client.Get(rs.Primary.ID)
			vpc_domain_policy := models.VPCDomainPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VPC Domain Policy %s Still exists", vpc_domain_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVPCDomainPolicyAttributes(description string, vpc_domain_policy *models.VPCDomainPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(vpc_domain_policy.DistinguishedName) {
			return fmt.Errorf("Bad vpc_inst_pol %s", GetMOName(vpc_domain_policy.DistinguishedName))
		}

		if description != vpc_domain_policy.Description {
			return fmt.Errorf("Bad vpc_domain_policy Description %s", vpc_domain_policy.Description)
		}

		if "200" != vpc_domain_policy.DeadIntvl {
			return fmt.Errorf("Bad vpc_domain_policy DeadIntvl %s", vpc_domain_policy.DeadIntvl)
		}

		if "test_annotation" != vpc_domain_policy.Annotation {
			return fmt.Errorf("Bad vpc_domain_policy Annotation %s", vpc_domain_policy.Annotation)
		}

		if "test_alias" != vpc_domain_policy.NameAlias {
			return fmt.Errorf("Bad vpc_domain_policy NameAlias")
		}
		return nil
	}
}
