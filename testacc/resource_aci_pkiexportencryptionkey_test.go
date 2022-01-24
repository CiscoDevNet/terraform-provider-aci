package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/terraform-providers/terraform-provider-aci/aci"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAciAESEncryptionPassphraseandKeysforConfigExportImport_Basic(t *testing.T) {
	var encryption_key_default models.AESEncryptionPassphraseandKeysforConfigExportImport
	var encryption_key_updated models.AESEncryptionPassphraseandKeysforConfigExportImport
	resourceName := "aci_encryption_key.test"
	pkiExportEncryptionKey, err := aci.GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(sharedAciClient(), "uni/exportcryptkey")
	if err != nil {
		t.Errorf("reading initial config of pkiExportEncryptionKey")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportExists(resourceName, &encryption_key_default),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
					resource.TestCheckResourceAttrSet(resourceName, "name_alias"),
					resource.TestCheckResourceAttr(resourceName, "clear_encryption_key", "no"),
					resource.TestCheckResourceAttr(resourceName, "passphrase_key_derivation_version", "v1"),
					resource.TestCheckResourceAttrSet(resourceName, "strong_encryption_enabled"),
				),
			},
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfigWithOptionalValues(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportExists(resourceName, &encryption_key_updated),
					resource.TestCheckResourceAttr(resourceName, "annotation", "orchestrator:terraform_testacc"),
					resource.TestCheckResourceAttr(resourceName, "description", "created while acceptance testing"),
					resource.TestCheckResourceAttr(resourceName, "name_alias", "test_encryption_key"),
					resource.TestCheckResourceAttr(resourceName, "clear_encryption_key", "no"),
					resource.TestCheckResourceAttr(resourceName, "passphrase", "abcdefghijklmnop"),
					resource.TestCheckResourceAttr(resourceName, "passphrase_key_derivation_version", "v1"),
					resource.TestCheckResourceAttr(resourceName, "strong_encryption_enabled", "yes"),
					testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportIdEqual(&encryption_key_default, &encryption_key_updated),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"clear_encryption_key", "passphrase"},
			},
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("strong_encryption_enabled", "no"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "strong_encryption_enabled", "no"),
					testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportIdEqual(&encryption_key_default, &encryption_key_updated),
				),
			},
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("clear_encryption_key", "yes"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "clear_encryption_key", "yes"),
					testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportIdEqual(&encryption_key_default, &encryption_key_updated),
				),
			},
			{
				Config: restoreAESEncryptionPassphraseandKeysforConfigExportImport(pkiExportEncryptionKey),
			},
		},
	})
}

func TestAccAciAESEncryptionPassphraseandKeysforConfigExportImport_Negative(t *testing.T) {
	randomParameter := acctest.RandStringFromCharSet(5, "abcdefghijklmnopqrstuvwxyz")
	randomValue := acctest.RandString(5)
	pkiExportEncryptionKey, err := aci.GetRemoteAESEncryptionPassphraseandKeysforConfigExportImport(sharedAciClient(), "uni/exportcryptkey")
	if err != nil {
		t.Errorf("reading initial config of pkiExportEncryptionKey")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:	  func(){ testAccPreCheck(t) },
		ProviderFactories:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfig(),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("description", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("annotation", acctest.RandString(129)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},
			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("name_alias", acctest.RandString(64)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("clear_encryption_key", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("passphrase", acctest.RandString(15)),
				ExpectError: regexp.MustCompile(`smaller than minimum required`),
			},
			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("passphrase", acctest.RandString(33)),
				ExpectError: regexp.MustCompile(`failed validation for value '(.)+'`),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("passphrase_key_derivation_version", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr("strong_encryption_enabled", randomValue),
				ExpectError: regexp.MustCompile(`expected(.)+ to be one of (.)+, got(.)+`),
			},

			{
				Config:      CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr(randomParameter, randomValue),
				ExpectError: regexp.MustCompile(`An argument named (.)+ is not expected here.`),
			},
			{
				Config: restoreAESEncryptionPassphraseandKeysforConfigExportImport(pkiExportEncryptionKey),
			},
		},
	})
}

func restoreAESEncryptionPassphraseandKeysforConfigExportImport(pkiExportEncryptionKey *models.AESEncryptionPassphraseandKeysforConfigExportImport) string {
	var resource string
	if pkiExportEncryptionKey.KeyConfigured == "yes" {
		resource = fmt.Sprintf(`
		resource "aci_encryption_key" "test" {
			description = "%s"
			annotation = "%s"
			name_alias = "%s"
			clear_encryption_key = "no"
			passphrase = "abcdefghijklmnop"
			passphrase_key_derivation_version = "%s"
			strong_encryption_enabled = "%s"
		}`, pkiExportEncryptionKey.Description, pkiExportEncryptionKey.Annotation, pkiExportEncryptionKey.NameAlias, pkiExportEncryptionKey.PassphraseKeyDerivationVersion, pkiExportEncryptionKey.StrongEncryptionEnabled)
	} else {
		resource = fmt.Sprintf(`
		resource "aci_encryption_key" "test" {
			description = "%s"
			annotation = "%s"
			name_alias = "%s"
			clear_encryption_key = "yes"
		}
		`, pkiExportEncryptionKey.Description, pkiExportEncryptionKey.Annotation, pkiExportEncryptionKey.NameAlias)
	}
	return resource
}

func testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportExists(name string, encryption_key *models.AESEncryptionPassphraseandKeysforConfigExportImport) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("AES Encryption Passphraseand KeysforConfig Export Import %s not found", name)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No AES Encryption Passphraseand KeysforConfig Export Import dn was set")
		}

		client := testAccProvider.Meta().(*client.Client)

		cont, err := client.Get(rs.Primary.ID)
		if err != nil {
			return err
		}

		encryption_keyFound := models.AESEncryptionPassphraseandKeysforConfigExportImportFromContainer(cont)
		if encryption_keyFound.DistinguishedName != rs.Primary.ID {
			return fmt.Errorf("AES Encryption Passphraseand KeysforConfig Export Import %s not found", rs.Primary.ID)
		}
		*encryption_key = *encryption_keyFound
		return nil
	}
}

func testAccCheckAciAESEncryptionPassphraseandKeysforConfigExportImportIdEqual(m1, m2 *models.AESEncryptionPassphraseandKeysforConfigExportImport) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if m1.DistinguishedName != m2.DistinguishedName {
			return fmt.Errorf("encryption_key DNs are not equal")
		}
		return nil
	}
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfig() string {
	fmt.Println("=== STEP  Basic: testing encryption_key creation with required parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {}
	`)

	return resource
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportConfigWithOptionalValues() string {
	fmt.Println("=== STEP  Basic: testing encryption_key creation with optional parameters")
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		description = "created while acceptance testing"
		annotation = "orchestrator:terraform_testacc"
		name_alias = "test_encryption_key"
		clear_encryption_key = "no"
		passphrase = "abcdefghijklmnop"
		passphrase_key_derivation_version = "v1"
		strong_encryption_enabled = "yes"
		
	}
	`)

	return resource
}

func CreateAccAESEncryptionPassphraseandKeysforConfigExportImportUpdatedAttr(attribute, value string) string {
	fmt.Printf("=== STEP  testing encryption_key attribute: %s = %s \n", attribute, value)
	resource := fmt.Sprintf(`
	
	resource "aci_encryption_key" "test" {
	
		%s = "%s"
	}
	`, attribute, value)
	return resource
}
