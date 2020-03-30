package aci

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAciX509Certificate_Basic(t *testing.T) {
	var x509_certificate models.X509Certificate
	description := "x509_certificate created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciX509CertificateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists("aci_x509_certificate.foox509_certificate", &x509_certificate),
					testAccCheckAciX509CertificateAttributes(description, &x509_certificate),
				),
			},
			{
				ResourceName:      "aci_x509_certificate",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAciX509Certificate_update(t *testing.T) {
	var x509_certificate models.X509Certificate
	description := "x509_certificate created while acceptance testing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAciX509CertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAciX509CertificateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists("aci_x509_certificate.foox509_certificate", &x509_certificate),
					testAccCheckAciX509CertificateAttributes(description, &x509_certificate),
				),
			},
			{
				Config: testAccCheckAciX509CertificateConfig_basic(description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciX509CertificateExists("aci_x509_certificate.foox509_certificate", &x509_certificate),
					testAccCheckAciX509CertificateAttributes(description, &x509_certificate),
				),
			},
		},
	})
}

func testAccCheckAciX509CertificateConfig_basic(description string) string {
	return fmt.Sprintf(`

	resource "aci_x509_certificate" "foox509_certificate" {
		  local_user_dn  = "${aci_local_user.example.id}"
		description = "%s"
		
		name  = "example"
		  annotation  = "example"
		  data  = "example"
		  name_alias  = "example"
		}
	`, description)
}

func testAccCheckAciX509CertificateExists(name string, x509_certificate *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("X509 Certificate %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No X509 Certificate dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		x509_certificateFound := models.X509CertificateFromContainer(cont)
		if x509_certificateFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("X509 Certificate %s not found", rs.Primary.ID)
		}
		*x509_certificate = *x509_certificateFound
		return nil
	}
}

func testAccCheckAciX509CertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {

		if rs.Type == "aci_x509_certificate" {
			cont, err := client.Get(rs.Primary.ID)
			x509_certificate := models.X509CertificateFromContainer(cont)
			if err == nil {
				return fmt.Errorf("X509 Certificate %s Still exists", x509_certificate.DistinguishedName)
			}

		} else {
			continue
		}
	}

	return nil
}

func testAccCheckAciX509CertificateAttributes(description string, x509_certificate *models.X509Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if description != x509_certificate.Description {
			return fmt.Errorf("Bad x509_certificate Description %s", x509_certificate.Description)
		}

		if "example" != x509_certificate.Name {
			return fmt.Errorf("Bad x509_certificate name %s", x509_certificate.Name)
		}

		if "example" != x509_certificate.Annotation {
			return fmt.Errorf("Bad x509_certificate annotation %s", x509_certificate.Annotation)
		}

		if "example" != x509_certificate.Data {
			return fmt.Errorf("Bad x509_certificate data %s", x509_certificate.Data)
		}

		if "example" != x509_certificate.NameAlias {
			return fmt.Errorf("Bad x509_certificate name_alias %s", x509_certificate.NameAlias)
		}

		return nil
	}
}
