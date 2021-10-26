package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciOSPFInterfacePolicy_Basic(t *testing.T) {
	var ospf_interface_policy models.OSPFInterfacePolicy
	description := "ospf_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFInterfacePolicyConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists("aci_ospf_interface_policy.fooospf_interface_policy", &ospf_interface_policy),
					testAccCheckAciOSPFInterfacePolicyAttributes(description, "unspecified", &ospf_interface_policy),
				),
			},
		},
	})
}

func TestAccAciOSPFInterfacePolicy_update(t *testing.T) {
	var ospf_interface_policy models.OSPFInterfacePolicy
	description := "ospf_interface_policy created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciOSPFInterfacePolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciOSPFInterfacePolicyConfig_basic(description, "unspecified"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists("aci_ospf_interface_policy.fooospf_interface_policy", &ospf_interface_policy),
					testAccCheckAciOSPFInterfacePolicyAttributes(description, "unspecified", &ospf_interface_policy),
				),
			},
			{
				Config: testAccCheckAciOSPFInterfacePolicyConfig_basic(description, "passive"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciOSPFInterfacePolicyExists("aci_ospf_interface_policy.fooospf_interface_policy", &ospf_interface_policy),
					testAccCheckAciOSPFInterfacePolicyAttributes(description, "passive", &ospf_interface_policy),
				),
			},
		},
	})
}

func testAccCheckAciOSPFInterfacePolicyConfig_basic(description, ctrl string) string {
	return fmt.Sprintf(`

	resource "aci_ospf_interface_policy" "fooospf_interface_policy" {
		tenant_dn    = aci_tenant.example.id
		description  = "%s"
		name         = "demo_ospfpol"
		annotation   = "tag_ospf"
		cost         = "unspecified"
		ctrl         = "%s"
		dead_intvl   = "40"
		hello_intvl  = "10"
		name_alias   = "alias_ospf"
		nw_t         = "unspecified"
		pfx_suppress = "inherit"
		prio         = "1"
		rexmit_intvl = "5"
		xmit_delay   = "1"
	}
	`, description, ctrl)
}

func testAccCheckAciOSPFInterfacePolicyExists(name string, ospf_interface_policy *models.OSPFInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("OSPF Interface Policy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No OSPF Interface Policy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ospf_interface_policyFound := models.OSPFInterfacePolicyFromContainer(cont)
		if ospf_interface_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("OSPF Interface Policy %s not found", rs.Primary.ID)
		}
		*ospf_interface_policy = *ospf_interface_policyFound
		return nil
	}
}

func testAccCheckAciOSPFInterfacePolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_ospf_interface_policy" {
			cont, err := client.Get(rs.Primary.ID)
			ospf_interface_policy := models.OSPFInterfacePolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("OSPF Interface Policy %s Still exists", ospf_interface_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciOSPFInterfacePolicyAttributes(description, ctrl string, ospf_interface_policy *models.OSPFInterfacePolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != ospf_interface_policy.Description {
			return fmt.Errorf("Bad ospf_interface_policy Description %s", ospf_interface_policy.Description)
		}

		if "demo_ospfpol" != ospf_interface_policy.Name {
			return fmt.Errorf("Bad ospf_interface_policy name %s", ospf_interface_policy.Name)
		}

		if "tag_ospf" != ospf_interface_policy.Annotation {
			return fmt.Errorf("Bad ospf_interface_policy annotation %s", ospf_interface_policy.Annotation)
		}

		if "unspecified" != ospf_interface_policy.Cost {
			return fmt.Errorf("Bad ospf_interface_policy cost %s", ospf_interface_policy.Cost)
		}

		policyCtrl := ospf_interface_policy.Ctrl
		if policyCtrl == "" {
			policyCtrl = "unspecified"
		}
		if ctrl != policyCtrl {
			return fmt.Errorf("Bad ospf_interface_policy ctrl %s", ospf_interface_policy.Ctrl)
		}

		if "40" != ospf_interface_policy.DeadIntvl {
			return fmt.Errorf("Bad ospf_interface_policy dead_intvl %s", ospf_interface_policy.DeadIntvl)
		}

		if "10" != ospf_interface_policy.HelloIntvl {
			return fmt.Errorf("Bad ospf_interface_policy hello_intvl %s", ospf_interface_policy.HelloIntvl)
		}

		if "alias_ospf" != ospf_interface_policy.NameAlias {
			return fmt.Errorf("Bad ospf_interface_policy name_alias %s", ospf_interface_policy.NameAlias)
		}

		if "unspecified" != ospf_interface_policy.NwT {
			return fmt.Errorf("Bad ospf_interface_policy nw_t %s", ospf_interface_policy.NwT)
		}

		if "inherit" != ospf_interface_policy.PfxSuppress {
			return fmt.Errorf("Bad ospf_interface_policy pfx_suppress %s", ospf_interface_policy.PfxSuppress)
		}

		if "1" != ospf_interface_policy.Prio {
			return fmt.Errorf("Bad ospf_interface_policy prio %s", ospf_interface_policy.Prio)
		}

		if "5" != ospf_interface_policy.RexmitIntvl {
			return fmt.Errorf("Bad ospf_interface_policy rexmit_intvl %s", ospf_interface_policy.RexmitIntvl)
		}

		if "1" != ospf_interface_policy.XmitDelay {
			return fmt.Errorf("Bad ospf_interface_policy xmit_delay %s", ospf_interface_policy.XmitDelay)
		}

		return nil
	}
}
