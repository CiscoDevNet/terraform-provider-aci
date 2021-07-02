package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciPhysicalDomain_Basic(t *testing.T) {
	var physical_domain models.PhysicalDomain
	name_alias := "physical_domain"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPhysicalDomainConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists("aci_physical_domain.foophysical_domain", &physical_domain),
					testAccCheckAciPhysicalDomainAttributes(name_alias, &physical_domain),
				),
			},
		},
	})
}

func TestAccAciPhysicalDomain_update(t *testing.T) {
	var physical_domain models.PhysicalDomain
	name_alias := "physical_domain"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciPhysicalDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciPhysicalDomainConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists("aci_physical_domain.foophysical_domain", &physical_domain),
					testAccCheckAciPhysicalDomainAttributes(name_alias, &physical_domain),
				),
			},
			{
				Config: testAccCheckAciPhysicalDomainConfig_basic(name_alias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciPhysicalDomainExists("aci_physical_domain.foophysical_domain", &physical_domain),
					testAccCheckAciPhysicalDomainAttributes(name_alias, &physical_domain),
				),
			},
		},
	})
}

func testAccCheckAciPhysicalDomainConfig_basic(name_alias string) string {
	return fmt.Sprintf(`

	resource "aci_physical_domain" "foophysical_domain" {
			name  = "example"
		  	annotation  = "example"
		  	name_alias  = "%s"
		}
	`, name_alias)
}

func testAccCheckAciPhysicalDomainExists(name string, physical_domain *models.PhysicalDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Physical Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Physical Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		physical_domainFound := models.PhysicalDomainFromContainer(cont)
		if physical_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Physical Domain %s not found", rs.Primary.ID)
		}
		*physical_domain = *physical_domainFound
		return nil
	}
}

func testAccCheckAciPhysicalDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_physical_domain" {
			cont, err := client.Get(rs.Primary.ID)
			physical_domain := models.PhysicalDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Physical Domain %s Still exists", physical_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciPhysicalDomainAttributes(name_alias string, physical_domain *models.PhysicalDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// if description != physical_domain.Description {
		// 	return fmt.Errorf("Bad physical_domain Description %s", physical_domain.Description)
		// }

		if "example" != physical_domain.Name {
			return fmt.Errorf("Bad physical_domain name %s", physical_domain.Name)
		}

		if "example" != physical_domain.Annotation {
			return fmt.Errorf("Bad physical_domain annotation %s", physical_domain.Annotation)
		}

		if name_alias != physical_domain.NameAlias {
			return fmt.Errorf("Bad physical_domain name_alias %s", physical_domain.NameAlias)
		}

		return nil
	}
}
