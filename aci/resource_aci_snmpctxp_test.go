package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSNMPContextProfile_Basic(t *testing.T) {
	var snmp_context_profile models.SNMPContextProfile
	name_alias := "snmp_context_alias"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPContextProfileConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPContextProfileExists("aci_vrf_snmp_context.foosnmp_context_profile", &snmp_context_profile),
					testAccCheckAciSNMPContextProfileAttributes(name_alias, &snmp_context_profile),
				),
			},
		},
	})
}

func TestAccAciSNMPContextProfile_Update(t *testing.T) {
	var snmp_context_profile models.SNMPContextProfile
	name_alias := "snmp_context_alias"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSNMPContextProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSNMPContextProfileConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPContextProfileExists("aci_vrf_snmp_context.foosnmp_context_profile", &snmp_context_profile),
					testAccCheckAciSNMPContextProfileAttributes(name_alias, &snmp_context_profile),
				),
			},
			{
				Config: testAccCheckAciSNMPContextProfileConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSNMPContextProfileExists("aci_vrf_snmp_context.foosnmp_context_profile", &snmp_context_profile),
					testAccCheckAciSNMPContextProfileAttributes(name_alias, &snmp_context_profile),
				),
			},
		},
	})
}

func testAccCheckAciSNMPContextProfileConfig_basic(name_alias string) string {
	return fmt.Sprintf(`
	resource "aci_vrf_snmp_context" "foosnmp_context_profile" {
		name = "test"
		vrf_snmp_context_dn = "uni/tn-AS_tenant/ctx-harsh_test/snmpctx"
		annotation = "test_annotation"
		name_alias = "%s"
	}
	`, name_alias)
}

func testAccCheckAciSNMPContextProfileExists(name string, snmp_context_profile *models.SNMPContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SNMP Context  Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SNMP Context  Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		snmp_context_profileFound := models.SNMPContextProfileFromContainer(cont)
		if snmp_context_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SNMP Context  Profile %s not found", rs.Primary.ID)
		}
		*snmp_context_profile = *snmp_context_profileFound
		return nil
	}
}

func testAccCheckAciSNMPContextProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vrf_snmp_context" {
			cont, err := client.Get(rs.Primary.ID)
			snmp_context_profile := models.SNMPContextProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SNMP Context  Profile %s Still exists", snmp_context_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSNMPContextProfileAttributes(name_alias string, snmp_context_profile *models.SNMPContextProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != snmp_context_profile.Name {
			return fmt.Errorf("Bad snmp_ctx_p %s", snmp_context_profile.Name)
		}

		if "test_annotation" != snmp_context_profile.Annotation {
			return fmt.Errorf("Bad vrf_snmp_context Annotation %s", snmp_context_profile.Annotation)
		}

		if name_alias != snmp_context_profile.NameAlias {
			return fmt.Errorf("Bad vrf_snmp_context NameAlias %s", snmp_context_profile.NameAlias)
		}
		return nil
	}
}
