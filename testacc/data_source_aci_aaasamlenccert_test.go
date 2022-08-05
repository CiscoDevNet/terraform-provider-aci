package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAciKeypairforSAMLEncryptionDataSource_Basic(t *testing.T) {
	dataSourceName := "data.aci_saml_certificate.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccKeypairforSAMLEncryptionConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "regenerate"),
					resource.TestCheckResourceAttrSet(dataSourceName, "expiry_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate_validty"),
					resource.TestCheckResourceAttrSet(dataSourceName, "certificate"),
				),
			},
			{
				Config:      CreateAccKeypairforSAMLEncryptionDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: CreateAccKeypairforSAMLEncryptionConfigDataSource(),
			},
		},
	})
}

func CreateAccKeypairforSAMLEncryptionConfigDataSource() string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with required arguments only")
	resource := fmt.Sprintf(`
	data "aci_saml_certificate" "test" {}
	`)
	return resource
}

func CreateAccKeypairforSAMLEncryptionDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing saml_certificate Data Source with random attribute")
	resource := fmt.Sprintf(`
	data "aci_saml_certificate" "test" {
		%s = "%s"
	}
	`, key, value)
	return resource
}
