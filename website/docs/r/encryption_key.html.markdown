---
layout: "aci"
page_title: "ACI: aci_aes_encryption_passphraseand_keysfor_config_export(and_import)"
sidebar_current: "docs-aci-resource-aes_encryption_passphraseand_keysfor_config_export(and_import)"
description: |-
  Manages ACI AES Encryption Passphrase and Keys for Config Export (and Import)
---

# aci_aes_encryption_passphraseand_keysfor_config_export(and_import) #

Manages ACI AES Encryption Passphrase and Keys for Config Export (and Import)

## API Information ##

* `Class` - pkiExportEncryptionKey
* `Distinguished Named` - uni/exportcryptkey

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_aes_encryption_passphraseand_keysfor_config_export(and_import)" "example" {

  annotation = "orchestrator:terraform"
  clear_encryption_key = "false"

  passphrase = 
  passphrase_key_derivation_version = "0"
  strong_encryption_enabled = "false"
}
```

## Argument Reference ##



* `annotation` - (Optional) Annotation of object AES Encryption Passphrase and Keys for Config Export (and Import).

* `clear_encryption_key` - (Optional) Pushbutton property to clear the encryption key, if configured.Setting this property to true will trigger the clearing of all fields in this mo,
             set the strongEncryptionEnabled policy to False and keyConfigured to False. There is no
             method to recover the previous passphrase before the clear operation. Allowed values are "no", "yes", and default value is "false". Type: String.
* `passphrase` - (Optional) passphrase.The encryption parameters cannot be modified by a client request - only via a passphrase changeSetting this passphrase to blank/empty will trigger the clearing of all fields in this mo,
             set the strongEncryptionEnabled policy to False and keyConfigured to False. There is no
             method to recover the previous passphrase before the clear operation.
* `passphrase_key_derivation_version` - (Optional) passphraseKeyDerivationVersion.Version of the algorithm used - used for forward compatibility Allowed values are "v1", and default value is "0". Type: String.
* `strong_encryption_enabled` - (Optional) Strong Encryption Enabled for configuration export and import.Toggle to choose between weak and strong encryption - this flag can be set to True
           only when keyConfigured=True Allowed values are "no", "yes", and default value is "false". Type: String.


## Importing ##

An existing AESEncryptionPassphraseandKeysforConfigExport(andImport) can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_aes_encryption_passphraseand_keysfor_config_export(and_import).example <Dn>
```