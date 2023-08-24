package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBFDMultihopNodePolicy_Basic(t *testing.T) {
	var bfdmultihop_node_policy models.BFDMultihopNodePolicy
	fv_tenant_name := acctest.RandString(5)
	bfd_mh_node_pol_name := acctest.RandString(5)
	description := "bfdmultihop_node_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDMultihopNodePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDMultihopNodePolicyConfig_basic(fv_tenant_name, bfd_mh_node_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDMultihopNodePolicyExists("aci_bfdmultihop_node_policy.foo_bfdmultihop_node_policy", &bfdmultihop_node_policy),
					testAccCheckAciBFDMultihopNodePolicyAttributes(fv_tenant_name, bfd_mh_node_pol_name, description, &bfdmultihop_node_policy),
				),
			},
		},
	})
}

func TestAccAciBFDMultihopNodePolicy_Update(t *testing.T) {
	var bfdmultihop_node_policy models.BFDMultihopNodePolicy
	fv_tenant_name := acctest.RandString(5)
	bfd_mh_node_pol_name := acctest.RandString(5)
	description := "bfdmultihop_node_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDMultihopNodePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDMultihopNodePolicyConfig_basic(fv_tenant_name, bfd_mh_node_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDMultihopNodePolicyExists("aci_bfdmultihop_node_policy.foo_bfdmultihop_node_policy", &bfdmultihop_node_policy),
					testAccCheckAciBFDMultihopNodePolicyAttributes(fv_tenant_name, bfd_mh_node_pol_name, description, &bfdmultihop_node_policy),
				),
			},
			{
				Config: testAccCheckAciBFDMultihopNodePolicyConfig_basic(fv_tenant_name, bfd_mh_node_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDMultihopNodePolicyExists("aci_bfdmultihop_node_policy.foo_bfdmultihop_node_policy", &bfdmultihop_node_policy),
					testAccCheckAciBFDMultihopNodePolicyAttributes(fv_tenant_name, bfd_mh_node_pol_name, description, &bfdmultihop_node_policy),
				),
			},
		},
	})
}

func testAccCheckAciBFDMultihopNodePolicyConfig_basic(fv_tenant_name, bfd_mh_node_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_bfdmultihop_node_policy" "foo_bfdmultihop_node_policy" {
		name 		= "%s"
		description = "bfdmultihop_node_policy created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	`, fv_tenant_name, bfd_mh_node_pol_name)
}

func testAccCheckAciBFDMultihopNodePolicyExists(name string, bfdmultihop_node_policy *models.BFDMultihopNodePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BFD Multihop Node Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BFD Multihop Node Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bfdmultihop_node_policyFound := models.BFDMultihopNodePolicyFromContainer(cont)
		if bfdmultihop_node_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BFD Multihop Node Policy %s not found", rs.Primary.ID)
		}
		*bfdmultihop_node_policy = *bfdmultihop_node_policyFound
		return nil
	}
}

func testAccCheckAciBFDMultihopNodePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bfdmultihop_node_policy" {
			cont, err := client.Get(rs.Primary.ID)
			bfdmultihop_node_policy := models.BFDMultihopNodePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BFD Multihop Node Policy %s Still exists", bfdmultihop_node_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBFDMultihopNodePolicyAttributes(fv_tenant_name, bfd_mh_node_pol_name, description string, bfdmultihop_node_policy *models.BFDMultihopNodePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if bfd_mh_node_pol_name != GetMOName(bfdmultihop_node_policy.DistinguishedName) {
			return fmt.Errorf("Bad bfdmh_node_pol %s", GetMOName(bfdmultihop_node_policy.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(bfdmultihop_node_policy.DistinguishedName, bfdmultihop_node_policy.Rn)) {
			return fmt.Errorf(" Bad fvtenant %s", GetMOName(GetParentDn(bfdmultihop_node_policy.DistinguishedName, bfdmultihop_node_policy.Rn)))
		}
		if description != bfdmultihop_node_policy.Description {
			return fmt.Errorf("Bad bfdmultihop_node_policy Description %s", bfdmultihop_node_policy.Description)
		}
		return nil
	}
}
