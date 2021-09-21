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
	description := "snmp_community created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_vrf_snmp_context_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(description, &snmp_community),
				),
			},
		},
	})
}

func TestAccAciSNMPCommunity_Update(t *testing.T) {
	var snmp_community models.SNMPCommunity
	description := "snmp_community created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPCommunityDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_vrf_snmp_context_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(description, &snmp_community),
				),
			},
			{
				Config: testAccCheckAciSNMPCommunityConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPCommunityExists("aci_vrf_snmp_context_community.foosnmp_community", &snmp_community),
					testAccCheckAciSNMPCommunityAttributes(description, &snmp_community),
				),
			},
		},
	})
}

func testAccCheckAciSNMPCommunityConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_vrf_snmp_context_community" "foosnmp_community" {
		name 		= "test"
		description = "%s"
		vrf_dn = aci_vrf.test.id
		annotation = "Test_Annotation"
		name_alias = "Test_name_alias"
	}

	`, description)
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
		if rs.Type == "aci_vrf_snmp_context_community" {
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

func testAccCheckAciSNMPCommunityAttributes(description string, snmp_community *models.SNMPCommunity) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != snmp_community.Name {
			return fmt.Errorf("Bad snmp_community_p %s", snmp_community.Name)
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
