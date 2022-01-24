package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/terraform-providers/terraform-provider-aci/aci"
)

func TestAccAciAESEncryptionPassphraseandKeysforConfigExportImportDataSource_Basic(t *testing.T) {
	resourceName := "aci_encryption_key.test"
	dataSourceName := "data.aci_encryption_key.test"
	randomParameter := acctest.RandStringFromCharSet(10, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(10)
	pkiExportEncryptionKey, err := aci.GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(sharedAciClient(), "uni/exportcryptkey")
	if err != nil {
		t.Errorf("reading initial config of pkiExportEncryptionKey")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfigDataSource(),
				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name_alias", resourceName, "name_alias"),
					resource.TestCheckResourceAttrPair(dataSourceName, "passphrase_key_derivation_version", resourceName, "passphrase_key_derivation_version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "strong_encryption_enabled", resourceName, "strong_encryption_enabled"),
				),
			},
			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportDataSourceUpdate(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},

			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportDataSourceUpdatedResource("annotation", "orchestrator:terraform-testacc"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "annotation", resourceName, "annotation"),
				),
			},
			{
				Config: restoreAESEncryptionPassphraseandKeysforConfigExportImport(pkiExportEncryptionKey),
			},
		},
	})
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfigDataSource() string {
	fmt.Println("=== STEP  testing encryption_key Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}

	data "aci_encryption_key" "test" {
	
		depends_on = [ aci_encryption_key.test ]
	}
	`)
	return resource
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportDSWithInvalidName(string) string {
	fmt.Println("=== STEP  testing encryption_key Data Source with required arguments only")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}

	data "aci_encryption_key" "test" {
	
		depends_on = [ aci_encryption_key.test ]
	}
	`)
	return resource
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportDataSourceUpdate(key, value string) string {
	fmt.Println("=== STEP  testing encryption_key Data Source with random attribute")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
	}

	data "aci_encryption_key" "test" {
	
		%s = "%s"
		depends_on = [ aci_encryption_key.test ]
	}
	`, key, value)
	return resource
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportDataSourceUpdatedResource(key, value string) string {
	fmt.Println("=== STEP  testing encryption_key Data Source with updated resource")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		%s = "%s"
	}

	data "aci_encryption_key" "test" {
	
		depends_on = [ aci_encryption_key.test ]
	}
	`, key, value)
	return resource
}
