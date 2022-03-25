package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciInfraRsDomP_Basic(t *testing.T) {
	var InfraRsDomP models.InfraRsDomP
	infra_att_entity_p_name := acctest.RandString(5)
	infra_rs_dom_p_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInfraRsDomPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInfraRsDomPConfig_basic(infra_att_entity_p_name, infra_rs_dom_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInfraRsDomPExists("aci_aaep_to_domain.fooInfraRsDomP", &InfraRsDomP),
					testAccCheckAciInfraRsDomPAttributes(infra_att_entity_p_name, infra_rs_dom_p_name, &InfraRsDomP),
				),
			},
		},
	})
}

func TestAccAciInfraRsDomP_Update(t *testing.T) {
	var InfraRsDomP models.InfraRsDomP
	infra_att_entity_p_name := acctest.RandString(5)
	infra_rs_dom_p_name := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciInfraRsDomPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciInfraRsDomPConfig_basic(infra_att_entity_p_name, infra_rs_dom_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInfraRsDomPExists("aci_aaep_to_domain.fooInfraRsDomP", &InfraRsDomP),
					testAccCheckAciInfraRsDomPAttributes(infra_att_entity_p_name, infra_rs_dom_p_name, &InfraRsDomP),
				),
			},
			{
				Config: testAccCheckAciInfraRsDomPConfig_basic(infra_att_entity_p_name, infra_rs_dom_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciInfraRsDomPExists("aci_InfraRsDomP.fooInfraRsDomP", &InfraRsDomP),
					testAccCheckAciInfraRsDomPAttributes(infra_att_entity_p_name, infra_rs_dom_p_name, &InfraRsDomP),
				),
			},
		},
	})
}

func testAccCheckAciInfraRsDomPConfig_basic(infra_att_entity_p_name, infra_rs_dom_p_name string) string {
	return fmt.Sprintf(`
	resource "aci_attachable_access_entity_profile" "fooattachable_access_entity_profile" {
		name        = "%s"
		description = "attachable_access_entity_profile created while acceptance testing"
	  }
	  
	  resource "aci_l3_domain_profile" "fool3_domain_profile" {
		name       = "%s"
		annotation = "example"
		name_alias = "example"
	  }
	  
	  resource "aci_aaep_to_domain" "fooInfraRsDomP" {
		t_dn                                = aci_l3_domain_profile.fool3_domain_profile.id
		attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
	  } 
	`, infra_att_entity_p_name, infra_rs_dom_p_name, infra_rs_dom_p_name)
}

func testAccCheckAciInfraRsDomPExists(name string, InfraRsDomP *models.InfraRsDomP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("InfraRsDomP %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No InfraRsDomP dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		InfraRsDomPFound := models.InfraRsDomPFromContainer(cont)
		if InfraRsDomPFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("InfraRsDomP %s not found", rs.Primary.ID)
		}
		*InfraRsDomP = *InfraRsDomPFound
		return nil
	}
}

func testAccCheckAciInfraRsDomPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_aaep_to_domain" {
			cont, err := client.Get(rs.Primary.ID)
			InfraRsDomP := models.InfraRsDomPFromContainer(cont)
			if err == nil {
				return fmt.Errorf("InfraRsDomP %s Still exists", InfraRsDomP.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciInfraRsDomPAttributes(infra_att_entity_p_name, infra_rs_dom_p_name string, InfraRsDomP *models.InfraRsDomP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if infra_rs_dom_p_name != GetMOName(InfraRsDomP.DistinguishedName) {
			return fmt.Errorf("Bad infra_rs_dom_p %s", GetMOName(InfraRsDomP.DistinguishedName))
		}

		if infra_att_entity_p_name != GetMOName(GetParentDn(InfraRsDomP.DistinguishedName, InfraRsDomP.Rn)) {
			return fmt.Errorf(" Bad infra_att_entity_p %s", GetMOName(GetParentDn(InfraRsDomP.DistinguishedName, InfraRsDomP.Rn)))
		}
		return nil
	}
}
