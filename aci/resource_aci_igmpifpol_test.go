package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciIGMPInterfacePolicy_Basic(t *testing.T) {
	var igmpinterface_policy models.IGMPInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	igmp_if_pol_name := acctest.RandString(5)
	description := "igmpinterface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciIGMPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciIGMPInterfacePolicyConfig_basic(fv_tenant_name, igmp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIGMPInterfacePolicyExists("aci_igmpinterface_policy.foo_igmpinterface_policy", &igmpinterface_policy),
					testAccCheckAciIGMPInterfacePolicyAttributes(fv_tenant_name, igmp_if_pol_name, description, &igmpinterface_policy),
				),
			},
		},
	})
}

func TestAccAciIGMPInterfacePolicy_Update(t *testing.T) {
	var igmpinterface_policy models.IGMPInterfacePolicy
	fv_tenant_name := acctest.RandString(5)
	igmp_if_pol_name := acctest.RandString(5)
	description := "igmpinterface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciIGMPInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciIGMPInterfacePolicyConfig_basic(fv_tenant_name, igmp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIGMPInterfacePolicyExists("aci_igmpinterface_policy.foo_igmpinterface_policy", &igmpinterface_policy),
					testAccCheckAciIGMPInterfacePolicyAttributes(fv_tenant_name, igmp_if_pol_name, description, &igmpinterface_policy),
				),
			},
			{
				Config: testAccCheckAciIGMPInterfacePolicyConfig_basic(fv_tenant_name, igmp_if_pol_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciIGMPInterfacePolicyExists("aci_igmpinterface_policy.foo_igmpinterface_policy", &igmpinterface_policy),
					testAccCheckAciIGMPInterfacePolicyAttributes(fv_tenant_name, igmp_if_pol_name, description, &igmpinterface_policy),
				),
			},
		},
	})
}

func testAccCheckAciIGMPInterfacePolicyConfig_basic(fv_tenant_name, igmp_if_pol_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_igmpinterface_policy" "foo_igmpinterface_policy" {
		name 		= "%s"
		description = "igmpinterface_policy created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	`, fv_tenant_name, igmp_if_pol_name)
}

func testAccCheckAciIGMPInterfacePolicyExists(name string, igmpinterface_policy *models.IGMPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("IGMP Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IGMP Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		igmpinterface_policyFound := models.IGMPInterfacePolicyFromContainer(cont)
		if igmpinterface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("IGMP Interface Policy %s not found", rs.Primary.ID)
		}
		*igmpinterface_policy = *igmpinterface_policyFound
		return nil
	}
}

func testAccCheckAciIGMPInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_igmpinterface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			igmpinterface_policy := models.IGMPInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("IGMP Interface Policy %s Still exists", igmpinterface_policy.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciIGMPInterfacePolicyAttributes(fv_tenant_name, igmp_if_pol_name, description string, igmpinterface_policy *models.IGMPInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if igmp_if_pol_name != GetMOName(igmpinterface_policy.DistinguishedName) {
			return fmt.Errorf("Bad igmpif_pol %s", GetMOName(igmpinterface_policy.DistinguishedName))
		}

		if fv_tenant_name != GetMOName(GetParentDn(igmpinterface_policy.DistinguishedName, igmpinterface_policy.Rn)) {
			return fmt.Errorf(" Bad fvtenant %s", GetMOName(GetParentDn(igmpinterface_policy.DistinguishedName, igmpinterface_policy.Rn)))
		}
		if description != igmpinterface_policy.Description {
			return fmt.Errorf("Bad igmpinterface_policy Description %s", igmpinterface_policy.Description)
		}
		return nil
	}
}
