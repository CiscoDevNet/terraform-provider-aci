package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFEXProfile_Basic(t *testing.T) {
	var fex_profile models.FEXProfile
	description := "fex_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFEXProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFEXProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists("aci_fex_profile.example", &fex_profile),
					testAccCheckAciFEXProfileAttributes(description, &fex_profile),
				),
			},
		},
	})
}

func TestAccAciFEXProfile_update(t *testing.T) {
	var fex_profile models.FEXProfile
	description := "fex_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFEXProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFEXProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists("aci_fex_profile.example", &fex_profile),
					testAccCheckAciFEXProfileAttributes(description, &fex_profile),
				),
			},
			{
				Config: testAccCheckAciFEXProfileConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFEXProfileExists("aci_fex_profile.example", &fex_profile),
					testAccCheckAciFEXProfileAttributes(description, &fex_profile),
				),
			},
		},
	})
}

func testAccCheckAciFEXProfileConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_fex_profile" "example" {
		description = "%s"
		name  = "check"
		annotation  = "check"
		name_alias  = "check"
	}
	`, description)
}

func testAccCheckAciFEXProfileExists(name string, fex_profile *models.FEXProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("FEX Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No FEX Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fex_profileFound := models.FEXProfileFromContainer(cont)
		if fex_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("FEX Profile %s not found", rs.Primary.ID)
		}
		*fex_profile = *fex_profileFound
		return nil
	}
}

func testAccCheckAciFEXProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_fex_profile" {
			cont, err := client.Get(rs.Primary.ID)
			fex_profile := models.FEXProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("FEX Profile %s Still exists", fex_profile.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFEXProfileAttributes(description string, fex_profile *models.FEXProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != fex_profile.Description {
			return fmt.Errorf("Bad fex_profile Description %s", fex_profile.Description)
		}

		if "check" != fex_profile.Name {
			return fmt.Errorf("Bad fex_profile name %s", fex_profile.Name)
		}

		if "check" != fex_profile.Annotation {
			return fmt.Errorf("Bad fex_profile annotation %s", fex_profile.Annotation)
		}

		if "check" != fex_profile.NameAlias {
			return fmt.Errorf("Bad fex_profile name_alias %s", fex_profile.NameAlias)
		}

		return nil
	}
}
