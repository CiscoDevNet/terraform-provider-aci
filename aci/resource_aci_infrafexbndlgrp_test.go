package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciFexBundleGroup_Basic(t *testing.T) {
	var fex_bundle_group models.FexBundleGroup
	description := "fex_bundle_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFexBundleGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists("aci_fex_bundle_group.example", &fex_bundle_group),
					testAccCheckAciFexBundleGroupAttributes(description, &fex_bundle_group),
				),
			},
		},
	})
}

func TestAccAciFexBundleGroup_update(t *testing.T) {
	var fex_bundle_group models.FexBundleGroup
	description := "fex_bundle_group created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFexBundleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFexBundleGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists("aci_fex_bundle_group.example", &fex_bundle_group),
					testAccCheckAciFexBundleGroupAttributes(description, &fex_bundle_group),
				),
			},
			{
				Config: testAccCheckAciFexBundleGroupConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFexBundleGroupExists("aci_fex_bundle_group.example", &fex_bundle_group),
					testAccCheckAciFexBundleGroupAttributes(description, &fex_bundle_group),
				),
			},
		},
	})
}

func testAccCheckAciFexBundleGroupConfig_basic(description string) string {
	return fmt.Sprintf(`
	resource "aci_fex_bundle_group" "example" {
		fex_profile_dn  = "${aci_fex_profile.example.id}"
		description = "%s"
		name  = "fex_grp_check"
		annotation  = "fex_grp_check"
	  	name_alias  = "fex_grp_check"
	}
	`, description)
}

func testAccCheckAciFexBundleGroupExists(name string, fex_bundle_group *models.FexBundleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Fex Bundle Group %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Fex Bundle Group dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fex_bundle_groupFound := models.FexBundleGroupFromContainer(cont)
		if fex_bundle_groupFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Fex Bundle Group %s not found", rs.Primary.ID)
		}
		*fex_bundle_group = *fex_bundle_groupFound
		return nil
	}
}

func testAccCheckAciFexBundleGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_fex_bundle_group" {
			cont, err := client.Get(rs.Primary.ID)
			fex_bundle_group := models.FexBundleGroupFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Fex Bundle Group %s Still exists", fex_bundle_group.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFexBundleGroupAttributes(description string, fex_bundle_group *models.FexBundleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != fex_bundle_group.Description {
			return fmt.Errorf("Bad fex_bundle_group Description %s", fex_bundle_group.Description)
		}

		if "fex_grp_check" != fex_bundle_group.Name {
			return fmt.Errorf("Bad fex_bundle_group name %s", fex_bundle_group.Name)
		}

		if "fex_grp_check" != fex_bundle_group.Annotation {
			return fmt.Errorf("Bad fex_bundle_group annotation %s", fex_bundle_group.Annotation)
		}

		if "fex_grp_check" != fex_bundle_group.NameAlias {
			return fmt.Errorf("Bad fex_bundle_group name_alias %s", fex_bundle_group.NameAlias)
		}

		return nil
	}
}
