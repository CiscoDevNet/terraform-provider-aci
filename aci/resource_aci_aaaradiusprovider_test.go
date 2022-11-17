package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciRADIUSProvider_Basic(t *testing.T) {
	var radius_provider models.RADIUSProvider
	description := "radius_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRADIUSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRADIUSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderExists("aci_radius_provider.fooradius_provider", &radius_provider),
					testAccCheckAciRADIUSProviderAttributes(description, &radius_provider),
				),
			},
		},
	})
}

func TestAccAciRADIUSProvider_Update(t *testing.T) {
	var radius_provider models.RADIUSProvider
	description := "radius_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciRADIUSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciRADIUSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderExists("aci_radius_provider.fooradius_provider", &radius_provider),
					testAccCheckAciRADIUSProviderAttributes(description, &radius_provider),
				),
			},
			{
				Config: testAccCheckAciRADIUSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciRADIUSProviderExists("aci_radius_provider.fooradius_provider", &radius_provider),
					testAccCheckAciRADIUSProviderAttributes(description, &radius_provider),
				),
			},
		},
	})
}

func testAccCheckAciRADIUSProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_radius_provider" "fooradius_provider" {
		name = "test"
		type = "radius"
		description = "%s"
		annotation = "test_annotation_value"
		name_alias = "test_name_alias_value"
		auth_port = "1812"
		auth_protocol = "chap"
		monitor_server = "disabled"
		monitoring_user = "test_monitoring_user_value"
		retries = "1"
		timeout = "35"
	}

	`, description)
}

func testAccCheckAciRADIUSProviderExists(name string, radius_provider *models.RADIUSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("RADIUS Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No RADIUS Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		radius_providerFound := models.RADIUSProviderFromContainer(cont)
		if radius_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("RADIUS Provider %s not found", rs.Primary.ID)
		}
		*radius_provider = *radius_providerFound
		return nil
	}
}

func testAccCheckAciRADIUSProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_radius_provider" {
			cont, err := client.Get(rs.Primary.ID)
			radius_provider := models.RADIUSProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("RADIUS Provider %s Still exists", radius_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciRADIUSProviderAttributes(description string, radius_provider *models.RADIUSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != GetMOName(radius_provider.DistinguishedName) {
			return fmt.Errorf("Bad aaa_radius_provider %s", GetMOName(radius_provider.DistinguishedName))
		}

		if description != radius_provider.Description {
			return fmt.Errorf("Bad radius_provider Description %s", radius_provider.Description)
		}

		if "test_annotation_value" != radius_provider.Annotation {
			return fmt.Errorf("Bad radius_provider Annotation %s", radius_provider.Annotation)
		}

		if "test_name_alias_value" != radius_provider.NameAlias {
			return fmt.Errorf("Bad radius_provider NameAlias %s", radius_provider.NameAlias)
		}

		if "1812" != radius_provider.AuthPort {
			return fmt.Errorf("Bad radius_provider AuthPort %s", radius_provider.AuthPort)
		}

		if "chap" != radius_provider.AuthProtocol {
			return fmt.Errorf("Bad radius_provider AuthProtocol %s", radius_provider.AuthProtocol)
		}

		if "disabled" != radius_provider.MonitorServer {
			return fmt.Errorf("Bad radius_provider MonitorServer %s", radius_provider.MonitorServer)
		}

		if "test_monitoring_user_value" != radius_provider.MonitoringUser {
			return fmt.Errorf("Bad radius_provider MonitoringUser %s", radius_provider.MonitoringUser)
		}

		if "1" != radius_provider.Retries {
			return fmt.Errorf("Bad radius_provider Retries %s", radius_provider.Retries)
		}

		if "35" != radius_provider.Timeout {
			return fmt.Errorf("Bad radius_provider Timeout %s", radius_provider.Timeout)
		}
		return nil
	}
}
