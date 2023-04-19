package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRSAProvider_Basic(t *testing.T) {
	var rsa_provider models.RSAProvider
	description := "rsa_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRSAProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists("aci_rsa_provider.foorsa_provider", &rsa_provider),
					testAccCheckAciRSAProviderAttributes(description, &rsa_provider),
				),
			},
		},
	})
}

func TestAccAciRSAProvider_Update(t *testing.T) {
	var rsa_provider models.RSAProvider
	description := "rsa_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRSAProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRSAProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists("aci_rsa_provider.foorsa_provider", &rsa_provider),
					testAccCheckAciRSAProviderAttributes(description, &rsa_provider),
				),
			},
			{
				Config: testAccCheckAciRSAProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRSAProviderExists("aci_rsa_provider.foorsa_provider", &rsa_provider),
					testAccCheckAciRSAProviderAttributes(description, &rsa_provider),
				),
			},
		},
	})
}

func testAccCheckAciRSAProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_rsa_provider" "foorsa_provider" {
		name 				   = "test"
		name_alias             = "rsa_provider_alias"
		annotation             = "orchestrator:terraform"
		description 		   = "%s"
		auth_port              = "1812"
		auth_protocol          = "pap"
		key                    = ""
		monitor_server         = "disabled"
		monitoring_password    = ""
		monitoring_user        = "default"
		retries                = "1"
		timeout                = "5"
	}

	`, description)
}

func testAccCheckAciRSAProviderExists(name string, rsa_provider *models.RSAProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("RSA Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RSA Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		rsa_providerFound := models.RSAProviderFromContainer(cont)
		if rsa_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("RSA Provider %s not found", rs.Primary.ID)
		}
		*rsa_provider = *rsa_providerFound
		return nil
	}
}

func testAccCheckAciRSAProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_rsa_provider" {
			cont, err := client.Get(rs.Primary.ID)
			rsa_provider := models.RSAProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("RSA Provider %s Still exists", rsa_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRSAProviderAttributes(description string, rsa_provider *models.RSAProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != rsa_provider.Name {
			return fmt.Errorf("Bad aaa_rsa_provider %s", rsa_provider.Name)
		}

		if description != rsa_provider.Description {
			return fmt.Errorf("Bad rsa_provider Description %s", rsa_provider.Description)
		}

		if "rsa_provider_alias" != rsa_provider.NameAlias {
			return fmt.Errorf("Bad rsa_provider Name Alias %s", rsa_provider.NameAlias)
		}

		if "1" != rsa_provider.Retries {
			return fmt.Errorf("Bad rsa_provider Retries %s", rsa_provider.Retries)
		}

		if "5" != rsa_provider.Timeout {
			return fmt.Errorf("Bad rsa_provider Description %s", rsa_provider.Timeout)
		}
		return nil
	}
}
