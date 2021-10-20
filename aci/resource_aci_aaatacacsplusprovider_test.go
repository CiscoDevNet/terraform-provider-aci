package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciTACACSProvider_Basic(t *testing.T) {
	var TACACS_provider models.TACACSProvider
	description := "TACACS_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists("aci_tacacs_provider.fooTACACS_provider", &TACACS_provider),
					testAccCheckAciTACACSProviderAttributes(description, &TACACS_provider),
				),
			},
		},
	})
}

func TestAccAciTACACSProvider_Update(t *testing.T) {
	var TACACS_provider models.TACACSProvider
	description := "TACACS_provider created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciTACACSProviderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciTACACSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists("aci_tacacs_provider.fooTACACS_provider", &TACACS_provider),
					testAccCheckAciTACACSProviderAttributes(description, &TACACS_provider),
				),
			},
			{
				Config: testAccCheckAciTACACSProviderConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciTACACSProviderExists("aci_tacacs_provider.fooTACACS_provider", &TACACS_provider),
					testAccCheckAciTACACSProviderAttributes(description, &TACACS_provider),
				),
			},
		},
	})
}

func testAccCheckAciTACACSProviderConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_tacacs_provider" "fooTACACS_provider" {
		name 		= "test"
		description = "%s"
		annotation  = "example"
		name_alias  = "tacacs_provider_alias"

	}

	`, description)
}

func testAccCheckAciTACACSProviderExists(name string, TACACS_provider *models.TACACSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("TACACS Provider %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No TACACS Provider dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		TACACS_providerFound := models.TACACSProviderFromContainer(cont)
		if TACACS_providerFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("TACACS Provider %s not found", rs.Primary.ID)
		}
		*TACACS_provider = *TACACS_providerFound
		return nil
	}
}

func testAccCheckAciTACACSProviderDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "aci_tacacs_provider" {
			cont, err := client.Get(rs.Primary.ID)
			TACACS_provider := models.TACACSProviderFromContainer(cont)
			if err == nil {
				return fmt.Errorf("TACACS Provider %s Still exists", TACACS_provider.DistinguishedName)
			}
		} else {
			continue
		}
	}
	return nil
}

func testAccCheckAciTACACSProviderAttributes(description string, TACACS_provider *models.TACACSProvider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "test" != TACACS_provider.Name {
			return fmt.Errorf("Bad aaa_tacacs_plus_provider %s", TACACS_provider.Name)
		}

		if description != TACACS_provider.Description {
			return fmt.Errorf("Bad TACACS_provider Description %s", TACACS_provider.Description)
		}

		if "example" != TACACS_provider.Annotation {
			return fmt.Errorf("Bad TACACS_provider Annotation %s", TACACS_provider.Annotation)
		}

		if "tacacs_provider_alias" != TACACS_provider.NameAlias {
			return fmt.Errorf("Bad TACACS_provider Name Alias %s", TACACS_provider.NameAlias)
		}
		return nil
	}
}
