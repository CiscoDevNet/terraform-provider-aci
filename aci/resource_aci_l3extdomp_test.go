package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL3DomainProfile_Basic(t *testing.T) {
	var l3_domain_profile models.L3DomainProfile
	description := "l3_domain_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3DomainProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists("aci_l3_domain_profile.fool3_domain_profile", &l3_domain_profile),
					testAccCheckAciL3DomainProfileAttributes(description, &l3_domain_profile),
				),
			},
			{
				ResourceName:      "aci_l3_domain_profile",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciL3DomainProfile_update(t *testing.T) {
	var l3_domain_profile models.L3DomainProfile
	description := "l3_domain_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL3DomainProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL3DomainProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists("aci_l3_domain_profile.fool3_domain_profile", &l3_domain_profile),
					testAccCheckAciL3DomainProfileAttributes(description, &l3_domain_profile),
				),
			},
			{
				Config: testAccCheckAciL3DomainProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL3DomainProfileExists("aci_l3_domain_profile.fool3_domain_profile", &l3_domain_profile),
					testAccCheckAciL3DomainProfileAttributes(description, &l3_domain_profile),
				),
			},
		},
	})
}

func testAccCheckAciL3DomainProfileConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_l3_domain_profile" "fool3_domain_profile" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciL3DomainProfileExists(name string, l3_domain_profile *models.L3DomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L3 Domain Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L3 Domain Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l3_domain_profileFound := models.L3DomainProfileFromContainer(cont)
		if l3_domain_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L3 Domain Profile %s not found", rs.Primary.ID)
		}
		*l3_domain_profile = *l3_domain_profileFound
		return nil
	}
}

func testAccCheckAciL3DomainProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l3_domain_profile" {
			cont, err := client.Get(rs.Primary.ID)
			l3_domain_profile := models.L3DomainProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L3 Domain Profile %s Still exists", l3_domain_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL3DomainProfileAttributes(description string, l3_domain_profile *models.L3DomainProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != l3_domain_profile.Description {
			return fmt.Errorf("Bad l3_domain_profile Description %s", l3_domain_profile.Description)
		}

		if "example" != l3_domain_profile.Name {
			return fmt.Errorf("Bad l3_domain_profile name %s", l3_domain_profile.Name)
		}

		if "example" != l3_domain_profile.Annotation {
			return fmt.Errorf("Bad l3_domain_profile annotation %s", l3_domain_profile.Annotation)
		}

		if "example" != l3_domain_profile.NameAlias {
			return fmt.Errorf("Bad l3_domain_profile name_alias %s", l3_domain_profile.NameAlias)
		}

		return nil
	}
}
