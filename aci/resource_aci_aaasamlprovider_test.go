package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciSAMLProvider_Basic(t *testing.T) {
	var saml_provider models.SAMLProvider
	description := "saml_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSAMLProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists("aci_saml_provider.foosaml_provider", &saml_provider),
					testAccCheckAciSAMLProviderAttributes(description, &saml_provider),
				),
			},
		},
	})
}

func TestAccAciSAMLProvider_Update(t *testing.T) {
	var saml_provider models.SAMLProvider
	description := "saml_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciSAMLProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciSAMLProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists("aci_saml_provider.foosaml_provider", &saml_provider),
					testAccCheckAciSAMLProviderAttributes(description, &saml_provider),
				),
			},
			{
				Config: testAccCheckAciSAMLProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciSAMLProviderExists("aci_saml_provider.foosaml_provider", &saml_provider),
					testAccCheckAciSAMLProviderAttributes(description, &saml_provider),
				),
			},
		},
	})
}

func testAccCheckAciSAMLProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_saml_provider" "foosaml_provider" {
		name 		= "test"
		description = "%s"
		name_alias = "saml_provider_alias"
		annotation = "orchestrator:terraform"
		entity_id = "entity_id_example"
		gui_banner_message = "gui_banner_message_example"
		https_proxy = ""
		id_p = "adfs"
		key = ""
		metadata_url = ""
		monitor_server = "disabled"
		monitoring_password = ""
		monitoring_user = "default"
		retries = "1"
		sig_alg = "SIG_RSA_SHA256"
		timeout = "5"
		tp = ""
		want_assertions_encrypted = "yes"
		want_assertions_signed = "yes"
		want_requests_signed = "yes"
		want_response_signed = "yes"
	}

	`, description)
}

func testAccCheckAciSAMLProviderExists(name string, saml_provider *models.SAMLProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("SAML Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No SAML Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		saml_providerFound := models.SAMLProviderFromContainer(cont)
		if saml_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("SAML Provider %s not found", rs.Primary.ID)
		}
		*saml_provider = *saml_providerFound
		return nil
	}
}

func testAccCheckAciSAMLProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_saml_provider" {
			cont, err := client.Get(rs.Primary.ID)
			saml_provider := models.SAMLProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("SAML Provider %s Still exists", saml_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciSAMLProviderAttributes(description string, saml_provider *models.SAMLProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != saml_provider.Name {
			return fmt.Errorf("Bad aaa_saml_provider %s", saml_provider.Name)
		}

		if description != saml_provider.Description {
			return fmt.Errorf("Bad saml_provider Description %s", saml_provider.Description)
		}

		if "saml_provider_alias" != saml_provider.NameAlias {
			return fmt.Errorf("Bad saml_provider Name Alias %s", saml_provider.NameAlias)
		}

		if "adfs" != saml_provider.IdP {
			return fmt.Errorf("Bad saml_provider IdP %s", saml_provider.IdP)
		}

		if "1" != saml_provider.Retries {
			return fmt.Errorf("Bad saml_provider Retries %s", saml_provider.Retries)
		}

		if "5" != saml_provider.Timeout {
			return fmt.Errorf("Bad saml_provider Timeout %s", saml_provider.Timeout)
		}
		return nil
	}
}
