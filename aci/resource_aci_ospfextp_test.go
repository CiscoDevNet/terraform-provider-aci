package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciL3outOspfExternalPolicy_Basic(t *testing.T) {
	var l3out_ospf_external_policy models.L3outOspfExternalPolicy
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outOspfExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists("aci_l3out_ospf_external_policy.fool3out_ospf_external_policy", &l3out_ospf_external_policy),
					testAccCheckAciL3outOspfExternalPolicyAttributes(description, &l3out_ospf_external_policy),
				),
			},
		},
	})
}

func TestAccAciL3outOspfExternalPolicy_update(t *testing.T) {
	var l3out_ospf_external_policy models.L3outOspfExternalPolicy
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3outOspfExternalPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3outOspfExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists("aci_l3out_ospf_external_policy.fool3out_ospf_external_policy", &l3out_ospf_external_policy),
					testAccCheckAciL3outOspfExternalPolicyAttributes(description, &l3out_ospf_external_policy),
				),
			},
			{
				Config: testAccCheckAciL3outOspfExternalPolicyConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3outOspfExternalPolicyExists("aci_l3out_ospf_external_policy.fool3out_ospf_external_policy", &l3out_ospf_external_policy),
					testAccCheckAciL3outOspfExternalPolicyAttributes(description, &l3out_ospf_external_policy),
				),
			},
		},
	})
}

func testAccCheckAciL3outOspfExternalPolicyConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3out_ospf_external_policy" "fool3out_ospf_external_policy" {
		l3_outside_dn     = aci_l3_outside.example.id
		description       = "%s"
		annotation        = "example"
		area_cost         = "1"
		area_ctrl         = "redistribute"
		area_id           = "0.0.0.1"
		area_type         = "nssa"
		multipod_internal = "no"
		name_alias        = "example"
	}
	`, description)
}

func testAccCheckAciL3outOspfExternalPolicyExists(name string, l3out_ospf_external_policy *models.L3outOspfExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3outOspfExternalPolicy %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3outOspfExternalPolicy dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3out_ospf_external_policyFound := models.L3outOspfExternalPolicyFromContainer(cont)
		if l3out_ospf_external_policyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3outOspfExternalPolicy %s not found", rs.Primary.ID)
		}
		*l3out_ospf_external_policy = *l3out_ospf_external_policyFound
		return nil
	}
}

func testAccCheckAciL3outOspfExternalPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_external_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3out_ospf_external_policy := models.L3outOspfExternalPolicyFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3outOspfExternalPolicy %s Still exists", l3out_ospf_external_policy.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3outOspfExternalPolicyAttributes(description string, l3out_ospf_external_policy *models.L3outOspfExternalPolicy) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3out_ospf_external_policy.Description {
			return fmt.Errorf("Bad l3out_ospf_external_policy Description %s", l3out_ospf_external_policy.Description)
		}

		if "example" != l3out_ospf_external_policy.Annotation {
			return fmt.Errorf("Bad l3out_ospf_external_policy annotation %s", l3out_ospf_external_policy.Annotation)
		}

		if "1" != l3out_ospf_external_policy.AreaCost {
			return fmt.Errorf("Bad l3out_ospf_external_policy area_cost %s", l3out_ospf_external_policy.AreaCost)
		}

		if "redistribute" != l3out_ospf_external_policy.AreaCtrl {
			return fmt.Errorf("Bad l3out_ospf_external_policy area_ctrl %s", l3out_ospf_external_policy.AreaCtrl)
		}

		if "0.0.0.1" != l3out_ospf_external_policy.AreaId {
			return fmt.Errorf("Bad l3out_ospf_external_policy area_id %s", l3out_ospf_external_policy.AreaId)
		}

		if "nssa" != l3out_ospf_external_policy.AreaType {
			return fmt.Errorf("Bad l3out_ospf_external_policy area_type %s", l3out_ospf_external_policy.AreaType)
		}

		if "no" != l3out_ospf_external_policy.MultipodInternal {
			return fmt.Errorf("Bad l3out_ospf_external_policy multipod_internal %s", l3out_ospf_external_policy.MultipodInternal)
		}

		if "example" != l3out_ospf_external_policy.NameAlias {
			return fmt.Errorf("Bad l3out_ospf_external_policy name_alias %s", l3out_ospf_external_policy.NameAlias)
		}

		return nil
	}
}
