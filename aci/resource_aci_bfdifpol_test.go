package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciBFDInterfacePolicy_Basic(t *testing.T) {
	var bfd_interface_policy models.BFDInterfacePolicy
	description := "bfd_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfacePolicyExists("aci_bfd_interface_policy.test", &bfd_interface_policy),
					testAccCheckAciBFDInterfacePolicyAttributes(description, &bfd_interface_policy),
				),
			},
		},
	})
}

func TestAccAciBFDInterfacePolicy_update(t *testing.T) {
	var bfd_interface_policy models.BFDInterfacePolicy
	description := "bfd_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciBFDInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciBFDInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfacePolicyExists("aci_bfd_interface_policy.test", &bfd_interface_policy),
					testAccCheckAciBFDInterfacePolicyAttributes(description, &bfd_interface_policy),
				),
			},
			{
				Config: testAccCheckAciBFDInterfacePolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciBFDInterfacePolicyExists("aci_bfd_interface_policy.test", &bfd_interface_policy),
					testAccCheckAciBFDInterfacePolicyAttributes(description, &bfd_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciBFDInterfacePolicyConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_bfd_interface_policy" "test" {
		tenant_dn  = aci_tenant.example.id
		description = "%s"
		name  = "example"
  		admin_st = "disabled"
  		annotation  = "example"
  		ctrl = "opt-subif"
  		detect_mult  = "3"
  		echo_admin_st = "disabled"
  		echo_rx_intvl  = "50"
  		min_rx_intvl  = "50"
  		min_tx_intvl  = "50"
  		name_alias  = "example"
	}
	`, description)
}

func testAccCheckAciBFDInterfacePolicyExists(name string, bfd_interface_policy *models.BFDInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("BFD Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No BFD Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		bfd_interface_policyFound := models.BFDInterfacePolicyFromContainer(cont)
		if bfd_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("BFD Interface Policy %s not found", rs.Primary.ID)
		}
		*bfd_interface_policy = *bfd_interface_policyFound
		return nil
	}
}

func testAccCheckAciBFDInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_bfd_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			bfd_interface_policy := models.BFDInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("BFD Interface Policy %s Still exists", bfd_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciBFDInterfacePolicyAttributes(description string, bfd_interface_policy *models.BFDInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != bfd_interface_policy.Description {
			return fmt.Errorf("Bad bfd_interface_policy Description %s", bfd_interface_policy.Description)
		}

		if "example" != bfd_interface_policy.Name {
			return fmt.Errorf("Bad bfd_interface_policy name %s", bfd_interface_policy.Name)
		}

		if "disabled" != bfd_interface_policy.AdminSt {
			return fmt.Errorf("Bad bfd_interface_policy admin_st %s", bfd_interface_policy.AdminSt)
		}

		if "example" != bfd_interface_policy.Annotation {
			return fmt.Errorf("Bad bfd_interface_policy annotation %s", bfd_interface_policy.Annotation)
		}

		if "opt-subif" != bfd_interface_policy.Ctrl {
			return fmt.Errorf("Bad bfd_interface_policy ctrl %s", bfd_interface_policy.Ctrl)
		}

		if "3" != bfd_interface_policy.DetectMult {
			return fmt.Errorf("Bad bfd_interface_policy detect_mult %s", bfd_interface_policy.DetectMult)
		}

		if "disabled" != bfd_interface_policy.EchoAdminSt {
			return fmt.Errorf("Bad bfd_interface_policy echo_admin_st %s", bfd_interface_policy.EchoAdminSt)
		}

		if "50" != bfd_interface_policy.EchoRxIntvl {
			return fmt.Errorf("Bad bfd_interface_policy echo_rx_intvl %s", bfd_interface_policy.EchoRxIntvl)
		}

		if "50" != bfd_interface_policy.MinRxIntvl {
			return fmt.Errorf("Bad bfd_interface_policy min_rx_intvl %s", bfd_interface_policy.MinRxIntvl)
		}

		if "50" != bfd_interface_policy.MinTxIntvl {
			return fmt.Errorf("Bad bfd_interface_policy min_tx_intvl %s", bfd_interface_policy.MinTxIntvl)
		}

		if "example" != bfd_interface_policy.NameAlias {
			return fmt.Errorf("Bad bfd_interface_policy name_alias %s", bfd_interface_policy.NameAlias)
		}

		return nil
	}
}
