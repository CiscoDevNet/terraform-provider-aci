package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPIMInterfacePolicy_Basic(t *testing.T) {
	var piminterface_policy models.PIMInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	pim_if_pol_name := acctest.RandString(5)
	description := "piminterface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPIMInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPIMInterfacePolicyConfig_basic(fv_tenant_name, pim_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPIMInterfacePolicyExists("aci_piminterface_policy.foo_piminterface_policy", &piminterface_policy),
					testAccCheckAciPIMInterfacePolicyAttributes(fv_tenant_name, pim_if_pol_name, description, &piminterface_policy),
				),
			},
		},
	})
}

func TestAccAciPIMInterfacePolicy_Update(t *testing.T) {
	var piminterface_policy models.PIMInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	pim_if_pol_name := acctest.RandString(5)
	description := "piminterface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPIMInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPIMInterfacePolicyConfig_basic(fv_tenant_name, pim_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPIMInterfacePolicyExists("aci_piminterface_policy.foo_piminterface_policy", &piminterface_policy),
					testAccCheckAciPIMInterfacePolicyAttributes(fv_tenant_name, pim_if_pol_name, description, &piminterface_policy),
				),
			},
			{
				Config: testAccCheckAciPIMInterfacePolicyConfig_basic(fv_tenant_name, pim_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPIMInterfacePolicyExists("aci_piminterface_policy.foo_piminterface_policy", &piminterface_policy),
					testAccCheckAciPIMInterfacePolicyAttributes(fv_tenant_name, pim_if_pol_name, description, &piminterface_policy),
				),
			},
		},
	})
}

func testAccCheckAciPIMInterfacePolicyConfig_basic(fv_tenant_name, pim_if_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_piminterface_policy" "foo_piminterface_policy" {
		name 		= "%s"
		description = "piminterface_policy created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	`, fv_tenant_name, pim_if_pol_name)
}

func testAccCheckAciPIMInterfacePolicyExists(name string, piminterface_policy *models.PIMInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("PIM Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No PIM Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		piminterface_policyFound := models.PIMInterfacePolicyFromContainer(cont)
		if piminterface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("PIM Interface Policy %s not found", rs.Primary.ID)
		}
		*piminterface_policy = *piminterface_policyFound
		return nil
	}
}

func testAccCheckAciPIMInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_piminterface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			piminterface_policy := models.PIMInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("PIM Interface Policy %s Still exists", piminterface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciPIMInterfacePolicyAttributes(fv_tenant_name, pim_if_pol_name, description string, piminterface_policy *models.PIMInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pim_if_pol_name != GetMOName(piminterface_policy.DistinguishedName) {
			return fmt.Errorf("Bad pimif_pol %s", GetMOName(piminterface_policy.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(piminterface_policy.DistinguishedName, piminterface_policy.Rn)) {
			return fmt.Errorf(" Bad fvtenant %s", GetMOName(GetParentDn(piminterface_policy.DistinguishedName, piminterface_policy.Rn)))
		}
		if description != piminterface_policy.Description {
			return fmt.Errorf("Bad piminterface_policy Description %s", piminterface_policy.Description)
		}
		return nil
	}
}
