package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciEndpointSecurityGroup_Basic(t *testing.T) {
	var endpoint_security_group models.EndpointSecurityGroup
	description := "endpoint_security_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndpointSecurityGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists("aci_endpoint_security_group.fooendpoint_security_group", &endpoint_security_group),
					testAccCheckAciEndpointSecurityGroupAttributes(description, &endpoint_security_group),
				),
			},
			{
				ResourceName:      "aci_endpoint_security_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciEndpointSecurityGroup_update(t *testing.T) {
	var endpoint_security_group models.EndpointSecurityGroup
	description := "endpoint_security_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciEndpointSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciEndpointSecurityGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists("aci_endpoint_security_group.fooendpoint_security_group", &endpoint_security_group),
					testAccCheckAciEndpointSecurityGroupAttributes(description, &endpoint_security_group),
				),
			},
			{
				Config: testAccCheckAciEndpointSecurityGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciEndpointSecurityGroupExists("aci_endpoint_security_group.fooendpoint_security_group", &endpoint_security_group),
					testAccCheckAciEndpointSecurityGroupAttributes(description, &endpoint_security_group),
				),
			},
		},
	})
}

func testAccCheckAciEndpointSecurityGroupConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_endpoint_security_group" "fooendpoint_security_group" {
		  application_profile_dn  = "${aci_application_profile.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  exception_tag  = "example"
		  flood_on_encap  = "example"
		  match_t  = "example"
		  name_alias  = "example"
		  pc_enf_pref  = "example"
		  pref_gr_memb  = "example"
		  prio  = "example"
		  userdom  = "example"
		}
	`, description)
}

func testAccCheckAciEndpointSecurityGroupExists(name string, endpoint_security_group *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Endpoint Security Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Endpoint Security Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		endpoint_security_groupFound := models.EndpointSecurityGroupFromContainer(cont)
		if endpoint_security_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Endpoint Security Group %s not found", rs.Primary.ID)
		}
		*endpoint_security_group = *endpoint_security_groupFound
		return nil
	}
}

func testAccCheckAciEndpointSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_endpoint_security_group" {
			cont, err := client.Get(rs.Primary.ID)
			endpoint_security_group := models.EndpointSecurityGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Endpoint Security Group %s Still exists", endpoint_security_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciEndpointSecurityGroupAttributes(description string, endpoint_security_group *models.EndpointSecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != endpoint_security_group.Description {
			return fmt.Errorf("Bad endpoint_security_group Description %s", endpoint_security_group.Description)
		}

		if "example" != endpoint_security_group.Name {
			return fmt.Errorf("Bad endpoint_security_group name %s", endpoint_security_group.Name)
		}

		if "example" != endpoint_security_group.Annotation {
			return fmt.Errorf("Bad endpoint_security_group annotation %s", endpoint_security_group.Annotation)
		}

		if "example" != endpoint_security_group.ExceptionTag {
			return fmt.Errorf("Bad endpoint_security_group exception_tag %s", endpoint_security_group.ExceptionTag)
		}

		if "example" != endpoint_security_group.FloodOnEncap {
			return fmt.Errorf("Bad endpoint_security_group flood_on_encap %s", endpoint_security_group.FloodOnEncap)
		}

		if "example" != endpoint_security_group.MatchT {
			return fmt.Errorf("Bad endpoint_security_group match_t %s", endpoint_security_group.MatchT)
		}

		if "example" != endpoint_security_group.NameAlias {
			return fmt.Errorf("Bad endpoint_security_group name_alias %s", endpoint_security_group.NameAlias)
		}

		if "example" != endpoint_security_group.PcEnfPref {
			return fmt.Errorf("Bad endpoint_security_group pc_enf_pref %s", endpoint_security_group.PcEnfPref)
		}

		if "example" != endpoint_security_group.PrefGrMemb {
			return fmt.Errorf("Bad endpoint_security_group pref_gr_memb %s", endpoint_security_group.PrefGrMemb)
		}

		if "example" != endpoint_security_group.Prio {
			return fmt.Errorf("Bad endpoint_security_group prio %s", endpoint_security_group.Prio)
		}

		if "example" != endpoint_security_group.Userdom {
			return fmt.Errorf("Bad endpoint_security_group userdom %s", endpoint_security_group.Userdom)
		}

		return nil
	}
}
