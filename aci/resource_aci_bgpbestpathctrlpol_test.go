package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciBgpBestPathPolicy_Basic(t *testing.T) {
	var bgp_best_path_policy models.BgpBestPathPolicy
	description := "bgp_best_path_controlpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpBestPathPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists("aci_bgp_best_path_policy.foobgp_best_path_policy", &bgp_best_path_policy),
					testAccCheckAciBgpBestPathPolicyAttributes(description, &bgp_best_path_policy),
				),
			},
		},
	})
}

func TestAccAciBgpBestPathPolicy_update(t *testing.T) {
	var bgp_best_path_policy models.BgpBestPathPolicy
	description := "bgp_best_path_controlpolicy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBgpBestPathPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBgpBestPathPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists("aci_bgp_best_path_policy.foobgp_best_path_policy", &bgp_best_path_policy),
					testAccCheckAciBgpBestPathPolicyAttributes(description, &bgp_best_path_policy),
				),
			},
			{
				Config: testAccCheckAciBgpBestPathPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBgpBestPathPolicyExists("aci_bgp_best_path_policy.foobgp_best_path_policy", &bgp_best_path_policy),
					testAccCheckAciBgpBestPathPolicyAttributes(description, &bgp_best_path_policy),
				),
			},
		},
	})
}

func testAccCheckAciBgpBestPathPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_bgp_best_path_policy" "foobgp_best_path_policy" {
		tenant_dn  = "${aci_tenant.example.id}"
		description = "%s"
		name  = "example"
  		annotation  = "example"
  		ctrl = "asPathMultipathRelax"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciBgpBestPathPolicyExists(name string, bgp_best_path_policy *models.BgpBestPathPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BgpBestPathPolicy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BgpBestPathPolicy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bgp_best_path_policyFound := models.BgpBestPathPolicyFromContainer(cont)
		if bgp_best_path_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BgpBestPathPolicy %s not found", rs.Primary.ID)
		}
		*bgp_best_path_policy = *bgp_best_path_policyFound
		return nil
	}
}

func testAccCheckAciBgpBestPathPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bgp_best_path_controlpolicy" {
			cont, err := client.Get(rs.Primary.ID)
			bgp_best_path_policy := models.BgpBestPathPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BgpBestPathPolicy %s Still exists", bgp_best_path_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBgpBestPathPolicyAttributes(description string, bgp_best_path_policy *models.BgpBestPathPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bgp_best_path_policy.Description {
			return fmt.Errorf("Bad bgp_best_path_policy Description %s", bgp_best_path_policy.Description)
		}

		if "example" != bgp_best_path_policy.Name {
			return fmt.Errorf("Bad bgp_best_path_policy name %s", bgp_best_path_policy.Name)
		}

		if "example" != bgp_best_path_policy.Annotation {
			return fmt.Errorf("Bad bgp_best_path_policy annotation %s", bgp_best_path_policy.Annotation)
		}

		if "asPathMultipathRelax" != bgp_best_path_policy.Ctrl {
			return fmt.Errorf("Bad bgp_best_path_policy ctrl %s", bgp_best_path_policy.Ctrl)
		}

		if "example" != bgp_best_path_policy.NameAlias {
			return fmt.Errorf("Bad bgp_best_path_policy name_alias %s", bgp_best_path_policy.NameAlias)
		}

		return nil
	}
}
