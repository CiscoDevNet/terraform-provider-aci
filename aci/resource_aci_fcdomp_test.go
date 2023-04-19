package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFCDomain_Basic(t *testing.T) {
	var fc_domain models.FCDomain
	description := "fc_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFCDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists("aci_fc_domain.foofc_domain", &fc_domain),
					testAccCheckAciFCDomainAttributes(description, &fc_domain),
				),
			},
			{
				ResourceName:      "aci_fc_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFCDomain_update(t *testing.T) {
	var fc_domain models.FCDomain
	description := "fc_domain created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFCDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFCDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists("aci_fc_domain.foofc_domain", &fc_domain),
					testAccCheckAciFCDomainAttributes(description, &fc_domain),
				),
			},
			{
				Config: testAccCheckAciFCDomainConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFCDomainExists("aci_fc_domain.foofc_domain", &fc_domain),
					testAccCheckAciFCDomainAttributes(description, &fc_domain),
				),
			},
		},
	})
}

func testAccCheckAciFCDomainConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_fc_domain" "foofc_domain" {
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciFCDomainExists(name string, fc_domain *models.FCDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("FC Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No FC Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		fc_domainFound := models.FCDomainFromContainer(cont)
		if fc_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("FC Domain %s not found", rs.Primary.ID)
		}
		*fc_domain = *fc_domainFound
		return nil
	}
}

func testAccCheckAciFCDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_fc_domain" {
			cont, err := client.Get(rs.Primary.ID)
			fc_domain := models.FCDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("FC Domain %s Still exists", fc_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFCDomainAttributes(description string, fc_domain *models.FCDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != fc_domain.Description {
			return fmt.Errorf("Bad fc_domain Description %s", fc_domain.Description)
		}

		if "example" != fc_domain.Name {
			return fmt.Errorf("Bad fc_domain name %s", fc_domain.Name)
		}

		if "example" != fc_domain.Annotation {
			return fmt.Errorf("Bad fc_domain annotation %s", fc_domain.Annotation)
		}

		if "example" != fc_domain.NameAlias {
			return fmt.Errorf("Bad fc_domain name_alias %s", fc_domain.NameAlias)
		}

		return nil
	}
}
