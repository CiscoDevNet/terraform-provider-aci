package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciFVDomain_Basic(t *testing.T) {
	var fvdomain models.FVDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFVDomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists("aci_epg_to_domain.check", &fvdomain),
					testAccCheckAciFVDomainAttributes(&fvdomain),
				),
			},
			{
				ResourceName:      "aci_epg_to_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciFVDomain_update(t *testing.T) {
	var domain models.FVDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciFVDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciFVDomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists("aci_epg_to_domain.check", &domain),
					testAccCheckAciFVDomainAttributes(&domain),
				),
			},
			{
				Config: testAccCheckAciFVDomainConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciFVDomainExists("aci_epg_to_domain.check", &domain),
					testAccCheckAciFVDomainAttributes(&domain),
				),
			},
		},
	})
}

func testAccCheckAciFVDomainConfig_basic() string {
	return fmt.Sprintf(`

	resource "aci_epg_to_domain" "check" {
		application_epg_dn  = "${aci_application_epg.epg2.id}"
		t_dn = "${aci_fc_domain.example.id}"
		vmm_allow_promiscuous = "accept"
		vmm_forged_transmits = "reject"
		vmm_mac_changes = "accept"
	}

	`)
}

func testAccCheckAciFVDomainExists(name string, domain *models.FVDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("Domain %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Domain dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		domainFound := models.FVDomainFromContainer(cont)
		if domainFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("Domain %s not found", rs.Primary.ID)
		}
		*domain = *domainFound
		return nil
	}
}

func testAccCheckAciFVDomainDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_epg_to_domain" {
			cont, err := client.Get(rs.Primary.ID)
			domain := models.FVDomainFromContainer(cont)
			if err == nil {
				return fmt.Errorf("Domain %s Still exists", domain.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciFVDomainAttributes(domain *models.FVDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if "none" != domain.BindingType {
			return fmt.Errorf("Bad domain binding_type %s", domain.BindingType)
		}

		if "encap" != domain.ClassPref {
			return fmt.Errorf("Bad domain class_pref %s", domain.ClassPref)
		}

		if "cos0" != domain.EpgCos {
			return fmt.Errorf("Bad domain epg_cos %s", domain.EpgCos)
		}
		return nil
	}
}
