package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciL2Domain_Basic(t *testing.T) {
	var l2_domain models.L2Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2DomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists("aci_l2_domain.fool2_domain", &l2_domain),
					testAccCheckAciL2DomainAttributes(&l2_domain),
				),
			},
		},
	})
}

func TestAccAciL2Domain_update(t *testing.T) {
	var l2_domain models.L2Domain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciL2DomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciL2DomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists("aci_l2_domain.fool2_domain", &l2_domain),
					testAccCheckAciL2DomainAttributes(&l2_domain),
				),
			},
			{
				Config: testAccCheckAciL2DomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciL2DomainExists("aci_l2_domain.fool2_domain", &l2_domain),
					testAccCheckAciL2DomainAttributes(&l2_domain),
				),
			},
		},
	})
}

func testAccCheckAciL2DomainConfig_basic() string {
	return fmt.Sprintf(`

	resource "aci_l2_domain" "fool2_domain" {
		name  = "example"
		annotation  = "example"
		name_alias  = "example"
		}
	`)
}

func testAccCheckAciL2DomainExists(name string, l2_domain *models.L2Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("L2 Domain Profile %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No L2 Domain Profile dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		l2_domainFound := models.L2DomainFromContainer(cont)
		if l2_domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("L2 Domain Profile %s not found", rs.Primary.ID)
		}
		*l2_domain = *l2_domainFound
		return nil
	}
}

func testAccCheckAciL2DomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_l2_domain" {
			cont, err := client.Get(rs.Primary.ID)
			l2_domain := models.L2DomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("L2 Domain Profile %s Still exists", l2_domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciL2DomainAttributes(l2_domain *models.L2Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if "example" != l2_domain.Name {
			return fmt.Errorf("Bad l2_domain name %s", l2_domain.Name)
		}

		if "example" != l2_domain.Annotation {
			return fmt.Errorf("Bad l2_domain annotation %s", l2_domain.Annotation)
		}

		if "example" != l2_domain.NameAlias {
			return fmt.Errorf("Bad l2_domain name_alias %s", l2_domain.NameAlias)
		}

		return nil
	}
}
