package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBfdMultihopInterfacePolicy_Basic(t *testing.T) {
	var aci_bfd_multihop_interface_policy models.AciBfdMultihopInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	bfd_mh_if_pol_name := acctest.RandString(5)
	description := "aci_bfd_multihop_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBfdMultihopInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBfdMultihopInterfacePolicyConfig_basic(fv_tenant_name, bfd_mh_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfacePolicyExists("aci_bfd_multihop_interface_policy.foo_bfd_interface_policy", &aci_bfd_multihop_interface_policy),
					testAccCheckAciBfdMultihopInterfacePolicyAttributes(fv_tenant_name, bfd_mh_if_pol_name, description, &aci_bfd_multihop_interface_policy),
				),
			},
		},
	})
}

func TestAccAciBfdMultihopInterfacePolicy_Update(t *testing.T) {
	var aci_bfd_multihop_interface_policy models.AciBfdMultihopInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	bfd_mh_if_pol_name := acctest.RandString(5)
	description := "aci_bfd_multihop_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBfdMultihopInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBfdMultihopInterfacePolicyConfig_basic(fv_tenant_name, bfd_mh_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfacePolicyExists("aci_bfd_multihop_interface_policy.foo_bfd_interface_policy", &aci_bfd_multihop_interface_policy),
					testAccCheckAciBfdMultihopInterfacePolicyAttributes(fv_tenant_name, bfd_mh_if_pol_name, description, &aci_bfd_multihop_interface_policy),
				),
			},
			{
				Config: testAccCheckAciBfdMultihopInterfacePolicyConfig_basic(fv_tenant_name, bfd_mh_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBfdMultihopInterfacePolicyExists("aci_bfd_multihop_interface_policy.foo_bfd_interface_policy", &aci_bfd_multihop_interface_policy),
					testAccCheckAciBfdMultihopInterfacePolicyAttributes(fv_tenant_name, bfd_mh_if_pol_name, description, &aci_bfd_multihop_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciBfdMultihopInterfacePolicyConfig_basic(fv_tenant_name, bfd_mh_if_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "footenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_bfd_multihop_interface_policy" "foo_bfd_interface_policy" {
		name 		= "%s"
		description = "aci_bfd_multihop_interface_policy created while acceptance testing"
		tenant_dn = aci_tenant.footenant.id
	}

	`, fv_tenant_name, bfd_mh_if_pol_name)
}

func testAccCheckAciBfdMultihopInterfacePolicyExists(name string, aci_bfd_multihop_interface_policy *models.AciBfdMultihopInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Aci BFD Multihop Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Aci BFD Multihop Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bfd_interface_policyFound := models.BFDInterfacePolicyFromContainer(cont)
		if bfd_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Aci BFD Multihop Interface Policy %s not found", rs.Primary.ID)
		}
		*aci_bfd_multihop_interface_policy = *bfd_interface_policyFound
		return nil
	}
}

func testAccCheckAciBfdMultihopInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_bfd_multihop_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			aci_bfd_multihop_interface_policy := models.BFDInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Aci BFD Multihop Interface Policy %s Still exists", aci_bfd_multihop_interface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciBfdMultihopInterfacePolicyAttributes(fv_tenant_name, bfd_mh_if_pol_name, description string, aci_bfd_multihop_interface_policy *models.AciBfdMultihopInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if bfd_mh_if_pol_name != GetMOName(aci_bfd_multihop_interface_policy.DistinguishedName) {
			return fmt.Errorf("Bad bfd_mh_if_pol %s", GetMOName(aci_bfd_multihop_interface_policy.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(aci_bfd_multihop_interface_policy.DistinguishedName)) {
			return fmt.Errorf(" Bad fv_tenant %s", GetMOName(GetParentDn(aci_bfd_multihop_interface_policy.DistinguishedName)))
		}
		if description != aci_bfd_multihop_interface_policy.Description {
			return fmt.Errorf("Bad aci_bfd_multihop_interface_policy Description %s", aci_bfd_multihop_interface_policy.Description)
		}
		return nil
	}
}
