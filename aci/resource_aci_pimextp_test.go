package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciExternalProfile_Basic(t *testing.T) {
	var external_profile models.ExternalProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	pim_ext_p_name := acctest.RandString(5)
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciExternalProfileConfig_basic(fv_tenant_name, l3ext_out_name, pim_ext_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalProfileExists("aci_external_profile.foo_external_profile", &external_profile),
					testAccCheckAciExternalProfileAttributes(fv_tenant_name, l3ext_out_name, pim_ext_p_name, description, &external_profile),
				),
			},
		},
	})
}

func TestAccAciExternalProfile_Update(t *testing.T) {
	var external_profile models.ExternalProfile
	fv_tenant_name := acctest.RandString(5)
	l3ext_out_name := acctest.RandString(5)
	pim_ext_p_name := acctest.RandString(5)
	description := "external_profile created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciExternalProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciExternalProfileConfig_basic(fv_tenant_name, l3ext_out_name, pim_ext_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalProfileExists("aci_external_profile.foo_external_profile", &external_profile),
					testAccCheckAciExternalProfileAttributes(fv_tenant_name, l3ext_out_name, pim_ext_p_name, description, &external_profile),
				),
			},
			{
				Config: testAccCheckAciExternalProfileConfig_basic(fv_tenant_name, l3ext_out_name, pim_ext_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciExternalProfileExists("aci_external_profile.foo_external_profile", &external_profile),
					testAccCheckAciExternalProfileAttributes(fv_tenant_name, l3ext_out_name, pim_ext_p_name, description, &external_profile),
				),
			},
		},
	})
}

func testAccCheckAciExternalProfileConfig_basic(fv_tenant_name, l3ext_out_name, pim_ext_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_tenant" "foo_tenant" {
		name 		= "%s"
		description = "tenant created while acceptance testing"

	}

	resource "aci_l3outside" "foo_l3outside" {
		name 		= "%s"
		description = "l3outside created while acceptance testing"
		tenant_dn = aci_tenant.foo_tenant.id
	}

	resource "aci_external_profile" "foo_external_profile" {
		name 		= "%s"
		description = "external_profile created while acceptance testing"
		l3outside_dn = aci_l3outside.foo_l3outside.id
	}

	`, fv_tenant_name, l3ext_out_name, pim_ext_p_name)
}

func testAccCheckAciExternalProfileExists(name string, external_profile *models.ExternalProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("External Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No External Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		external_profileFound := models.ExternalProfileFromContainer(cont)
		if external_profileFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("External Profile %s not found", rs.Primary.ID)
		}
		*external_profile = *external_profileFound
		return nil
	}
}

func testAccCheckAciExternalProfileDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_external_profile" {
			cont, err := client.Get(rs.Primary.ID)
			external_profile := models.ExternalProfileFromContainer(cont)
			if err == nil {
				return fmt.Errorf("External Profile %s Still exists", external_profile.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciExternalProfileAttributes(fv_tenant_name, l3ext_out_name, pim_ext_p_name, description string, external_profile *models.ExternalProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pim_ext_p_name != GetMOName(external_profile.DistinguishedName) {
			return fmt.Errorf("Bad pimext_p %s", GetMOName(external_profile.DistinguishedName))
		}

		if l3ext_out_name != GetMOName(GetParentDn(external_profile.DistinguishedName, external_profile.Rn)) {
			return fmt.Errorf(" Bad l3extout %s", GetMOName(GetParentDn(external_profile.DistinguishedName, external_profile.Rn)))
		}
		if description != external_profile.Description {
			return fmt.Errorf("Bad external_profile Description %s", external_profile.Description)
		}
		return nil
	}
}
