package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciExternalNetworkInstanceProfile_Basic(t *testing.T) {
	var external_network_instance_profile models.ExternalNetworkInstanceProfile
	description := "external_network_instance_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciExternalNetworkInstanceProfileConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists("aci_external_network_instance_profile.fooexternal_network_instance_profile", &external_network_instance_profile),
					testAccCheckAciExternalNetworkInstanceProfileAttributes(description, "AtleastOne", &external_network_instance_profile),
				),
			},
		},
	})
}

func TestAccAciExternalNetworkInstanceProfile_update(t *testing.T) {
	var external_network_instance_profile models.ExternalNetworkInstanceProfile
	description := "external_network_instance_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalNetworkInstanceProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciExternalNetworkInstanceProfileConfig_basic(description, "AtleastOne"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists("aci_external_network_instance_profile.fooexternal_network_instance_profile", &external_network_instance_profile),
					testAccCheckAciExternalNetworkInstanceProfileAttributes(description, "AtleastOne", &external_network_instance_profile),
				),
			},
			{
				Config: testAccCheckAciExternalNetworkInstanceProfileConfig_basic(description, "All"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalNetworkInstanceProfileExists("aci_external_network_instance_profile.fooexternal_network_instance_profile", &external_network_instance_profile),
					testAccCheckAciExternalNetworkInstanceProfileAttributes(description, "All", &external_network_instance_profile),
				),
			},
		},
	})
}

func testAccCheckAciExternalNetworkInstanceProfileConfig_basic(description, match_t string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "example" {
		name       = "test_acc_tenant"
	}

	resource "aci_l3_outside" "example" {
		tenant_dn      = aci_tenant.example.id
		name           = "demo_l3out"
		target_dscp = "CS0"
	}

	resource "aci_external_network_instance_profile" "fooexternal_network_instance_profile" {
		l3_outside_dn  = "${aci_l3_outside.example.id}"
		description    = "%s"
		name           = "demo_inst_prof"
		annotation     = "tag_network_profile"
		exception_tag  = "2"
		flood_on_encap = "disabled"
		match_t        = "%s"
		name_alias     = "alias_profile"
		pref_gr_memb   = "exclude"
		prio           = "level1"
		target_dscp    = "unspecified"
	}
	  
	`, description, match_t)
}

func testAccCheckAciExternalNetworkInstanceProfileExists(name string, external_network_instance_profile *models.ExternalNetworkInstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("External Network Instance Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No External Network Instance Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		external_network_instance_profileFound := models.ExternalNetworkInstanceProfileFromContainer(cont)
		if external_network_instance_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("External Network Instance Profile %s not found", rs.Primary.ID)
		}
		*external_network_instance_profile = *external_network_instance_profileFound
		return nil
	}
}

func testAccCheckAciExternalNetworkInstanceProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_external_network_instance_profile" {
			cont, err := client.Get(rs.Primary.ID)
			external_network_instance_profile := models.ExternalNetworkInstanceProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("External Network Instance Profile %s Still exists", external_network_instance_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciExternalNetworkInstanceProfileAttributes(description, match_t string, external_network_instance_profile *models.ExternalNetworkInstanceProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != external_network_instance_profile.Description {
			return fmt.Errorf("Bad external_network_instance_profile Description %s", external_network_instance_profile.Description)
		}

		if "demo_inst_prof" != external_network_instance_profile.Name {
			return fmt.Errorf("Bad external_network_instance_profile name %s", external_network_instance_profile.Name)
		}

		if "tag_network_profile" != external_network_instance_profile.Annotation {
			return fmt.Errorf("Bad external_network_instance_profile annotation %s", external_network_instance_profile.Annotation)
		}

		if "2" != external_network_instance_profile.ExceptionTag {
			return fmt.Errorf("Bad external_network_instance_profile exception_tag %s", external_network_instance_profile.ExceptionTag)
		}

		if "disabled" != external_network_instance_profile.FloodOnEncap {
			return fmt.Errorf("Bad external_network_instance_profile flood_on_encap %s", external_network_instance_profile.FloodOnEncap)
		}

		if match_t != external_network_instance_profile.MatchT {
			return fmt.Errorf("Bad external_network_instance_profile match_t %s", external_network_instance_profile.MatchT)
		}

		if "alias_profile" != external_network_instance_profile.NameAlias {
			return fmt.Errorf("Bad external_network_instance_profile name_alias %s", external_network_instance_profile.NameAlias)
		}

		if "exclude" != external_network_instance_profile.PrefGrMemb {
			return fmt.Errorf("Bad external_network_instance_profile pref_gr_memb %s", external_network_instance_profile.PrefGrMemb)
		}

		if "level1" != external_network_instance_profile.Prio {
			return fmt.Errorf("Bad external_network_instance_profile prio %s", external_network_instance_profile.Prio)
		}

		if "unspecified" != external_network_instance_profile.TargetDscp {
			return fmt.Errorf("Bad external_network_instance_profile target_dscp %s", external_network_instance_profile.TargetDscp)
		}

		return nil
	}
}
