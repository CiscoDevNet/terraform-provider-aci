package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciLDAPProvider_Basic(t *testing.T) {
	var ldap_provider models.LDAPProvider
	description := "ldap_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists("aci_ldap_provider.fooldap_provider", &ldap_provider),
					testAccCheckAciLDAPProviderAttributes(description, &ldap_provider),
				),
			},
		},
	})
}

func TestAccAciLDAPProvider_Update(t *testing.T) {
	var ldap_provider models.LDAPProvider
	description := "ldap_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciLDAPProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciLDAPProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists("aci_ldap_provider.fooldap_provider", &ldap_provider),
					testAccCheckAciLDAPProviderAttributes(description, &ldap_provider),
				),
			},
			{
				Config: testAccCheckAciLDAPProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciLDAPProviderExists("aci_ldap_provider.fooldap_provider", &ldap_provider),
					testAccCheckAciLDAPProviderAttributes(description, &ldap_provider),
				),
			},
		},
	})
}

func testAccCheckAciLDAPProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_ldap_provider" "fooldap_provider" {
		name = "test"
		type = "duo"
		description = "%s"
		annotation = "test_annotation"
		name_alias = "test_name_alias"
		ssl_validation_level = "strict"
		attribute = "test_attribute_value"
		basedn = "test_basedn_value"
		enable_ssl = "yes"
		filter = "test_filter_value"
		key = "test_key_value"
		monitor_server = "enabled"
		monitoring_password = "test_monitoring_password_value"
		monitoring_user = "test_monitoring_user_value"
		port = "389"
		retries = "1"
		rootdn = "test_rootdn_value"
		timeout = "30"
	}

	`, description)
}

func testAccCheckAciLDAPProviderExists(name string, ldap_provider *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("LDAP Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No LDAP Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		ldap_providerFound := models.LDAPProviderFromContainer(cont)
		if ldap_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("LDAP Provider %s not found", rs.Primary.ID)
		}
		*ldap_provider = *ldap_providerFound
		return nil
	}
}

func testAccCheckAciLDAPProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_ldap_provider" {
			cont, err := client.Get(rs.Primary.ID)
			ldap_provider := models.LDAPProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("LDAP Provider %s Still exists", ldap_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciLDAPProviderAttributes(description string, ldap_provider *models.LDAPProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(ldap_provider.DistinguishedName) {
			return fmt.Errorf("Bad aaa_ldap_provider %s", GetMOName(ldap_provider.DistinguishedName))
		}

		if description != ldap_provider.Description {
			return fmt.Errorf("Bad ldap_provider Description %s", ldap_provider.Description)
		}

		if "test_annotation" != ldap_provider.Annotation {
			return fmt.Errorf("Bad ldap_provider Annotation %s", ldap_provider.Annotation)
		}

		if "test_name_alias" != ldap_provider.NameAlias {
			return fmt.Errorf("Bad ldap_provider NameAlias %s", ldap_provider.NameAlias)
		}

		if "strict" != ldap_provider.SSLValidationLevel {
			return fmt.Errorf("Bad ldap_provider SSLValidationLevel %s", ldap_provider.SSLValidationLevel)
		}

		if "test_attribute_value" != ldap_provider.Attribute {
			return fmt.Errorf("Bad ldap_provider Attribute %s", ldap_provider.Attribute)
		}

		if "test_basedn_value" != ldap_provider.Basedn {
			return fmt.Errorf("Bad ldap_provider Basedn %s", ldap_provider.Basedn)
		}

		if "yes" != ldap_provider.EnableSSL {
			return fmt.Errorf("Bad ldap_provider EnableSSL %s", ldap_provider.EnableSSL)
		}

		if "test_filter_value" != ldap_provider.Filter {
			return fmt.Errorf("Bad ldap_provider Filter %s", ldap_provider.Filter)
		}

		if "enabled" != ldap_provider.MonitorServer {
			return fmt.Errorf("Bad ldap_provider MonitorServer %s", ldap_provider.MonitorServer)
		}

		if "test_monitoring_user_value" != ldap_provider.MonitoringUser {
			return fmt.Errorf("Bad ldap_provider MonitoringUser %s", ldap_provider.MonitoringUser)
		}

		if "389" != ldap_provider.Port {
			return fmt.Errorf("Bad ldap_provider Port %s", ldap_provider.Port)
		}

		if "1" != ldap_provider.Retries {
			return fmt.Errorf("Bad ldap_provider Retries %s", ldap_provider.Retries)
		}

		if "test_rootdn_value" != ldap_provider.Rootdn {
			return fmt.Errorf("Bad ldap_provider Rootdn %s", ldap_provider.Rootdn)
		}

		if "30" != ldap_provider.Timeout {
			return fmt.Errorf("Bad ldap_provider Timeout %s", ldap_provider.Timeout)
		}
		return nil
	}
}
