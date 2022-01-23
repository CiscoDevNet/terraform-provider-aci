package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSNMPCommunity_Basic(t *testing.T) {
	var snmp_community models.SNMPCommunity
	snmp_pol_name := acctest.RandString(5)
	snmp_community_p_name := acctest.RandString(5)
	description := "snmp_community created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(snmp_pol_name, snmp_community_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_snmp_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(snmp_pol_name, snmp_community_p_name, description, &snmp_community),
				),
			},
		},
	})
}

func TestAccAciSNMPCommunity_Update(t *testing.T) {
	var snmp_community models.SNMPCommunity
	snmp_pol_name := acctest.RandString(5)
	snmp_community_p_name := acctest.RandString(5)
	description := "snmp_community created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(snmp_pol_name, snmp_community_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_snmp_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(snmp_pol_name, snmp_community_p_name, description, &snmp_community),
				),
			},
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(snmp_pol_name, snmp_community_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_snmp_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(snmp_pol_name, snmp_community_p_name, description, &snmp_community),
				),
			},
		},
	})
}

func testAccCheckAciSNMPCommunityConfig_basic(snmp_pol_name, snmp_community_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_snmp_policy" "foosnmp_policy" {
		name 		= "%s"
		description = "snmp_policy created while acceptance testing"

	}

	resource "aci_snmp_community" "foosnmp_community" {
		name 		= "%s"
		description = "snmp_community created while acceptance testing"
		parent_dn = aci_snmp_policy.foosnmp_policy.id
	}

	`, snmp_pol_name, snmp_community_p_name)
}

func testAccCheckAciSNMPCommunityExists(name string, snmp_community *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SNMP Community %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SNMP Community dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		snmp_communityFound := models.SNMPCommunityFromContainer(cont)
		if snmp_communityFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SNMP Community %s not found", rs.Primary.ID)
		}
		*snmp_community = *snmp_communityFound
		return nil
	}
}

func testAccCheckAciSNMPCommunityDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_snmp_community" {
			cont, err := client.Get(rs.Primary.ID)
			snmp_community := models.SNMPCommunityFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SNMP Community %s Still exists", snmp_community.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSNMPCommunityAttributes(snmp_pol_name, snmp_community_p_name, description string, snmp_community *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if snmp_community_p_name != GetMOName(snmp_community.DistinguishedName) {
			return fmt.Errorf("Bad snmp_community_p %s", GetMOName(snmp_community.DistinguishedName))
		}

		if snmp_pol_name != GetMOName(GetParentDn(snmp_community.DistinguishedName)) {
			return fmt.Errorf(" Bad snmp_pol %s", GetMOName(GetParentDn(snmp_community.DistinguishedName)))
		}
		if description != snmp_community.Description {
			return fmt.Errorf("Bad snmp_community Description %s", snmp_community.Description)
		}
		if "Test_Annotation" != snmp_community.Annotation {
			return fmt.Errorf("Bad snmp_community Annotation %s", snmp_community.Annotation)
		}

		if "Test_name_alias" != snmp_community.NameAlias {
			return fmt.Errorf("Bad snmp_community Name Alias %s", snmp_community.NameAlias)
		}
		return nil
	}
}
