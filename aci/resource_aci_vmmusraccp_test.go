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

func TestAccAciVMMCredential_Basic(t *testing.T) {
	var vmm_credential models.VMMCredential
	vmm_prov_p_name := acctest.RandString(5)
	vmm_dom_p_name := acctest.RandString(5)
	vmm_usr_acc_p_name := acctest.RandString(5)
	description := "vmm_credential created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciVMMCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciVMMCredentialConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_usr_acc_p_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciVMMCredentialExists("aci_vmm_credential.foovmm_credential", &vmm_credential),
					testAccCheckAciVMMCredentialAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_usr_acc_p_name, description, &vmm_credential),
				),
			},
		},
	})
}

func testAccCheckAciVMMCredentialConfig_basic(vmm_prov_p_name, vmm_dom_p_name, vmm_usr_acc_p_name string) string {
	return fmt.Sprintf(`

	resource "aci_provider_profile" "fooprovider_profile" {
		name 		= "%s"
		description = "provider_profile created while acceptance testing"

	}

	resource "aci_vmm_domain" "foovmm_domain" {
		name 		= "%s"
		description = "vmm_domain created while acceptance testing"
		provider_profile_dn = aci_provider_profile.fooprovider_profile.id
	}

	resource "aci_vmm_credential" "foovmm_credential" {
		name 		= "%s"
		description = "vmm_credential created while acceptance testing"
		vmm_domain_dn = aci_vmm_domain.foovmm_domain.id
	}

	`, vmm_prov_p_name, vmm_dom_p_name, vmm_usr_acc_p_name)
}

func testAccCheckAciVMMCredentialExists(name string, vmm_credential *models.VMMCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("VMM Credential %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VMM Credential dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		vmm_credentialFound := models.VMMCredentialFromContainer(cont)
		if vmm_credentialFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("VMM Credential %s not found", rs.Primary.ID)
		}
		*vmm_credential = *vmm_credentialFound
		return nil
	}
}

func testAccCheckAciVMMCredentialDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_vmm_credential" {
			cont, err := client.Get(rs.Primary.ID)
			vmm_credential := models.VMMCredentialFromContainer(cont)
			if err == nil {
				return fmt.Errorf("VMM Credential %s Still exists", vmm_credential.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciVMMCredentialAttributes(vmm_prov_p_name, vmm_dom_p_name, vmm_usr_acc_p_name, description string, vmm_credential *models.VMMCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if vmm_usr_acc_p_name != GetMOName(vmm_credential.DistinguishedName) {
			return fmt.Errorf("Bad vmm_usr_acc_p %s", GetMOName(vmm_credential.DistinguishedName))
		}

		if vmm_dom_p_name != GetMOName(vmm_credential.DistinguishedName) {
			return fmt.Errorf(" Bad vmm_dom_p %s", GetMOName(vmm_credential.DistinguishedName))
		}
		if description != vmm_credential.Description {
			return fmt.Errorf("Bad vmm_credential Description %s", vmm_credential.Description)
		}
		return nil
	}
}
