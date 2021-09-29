---
layout: "aci"
page_title: "ACI: aci_encryption_key"
sidebar_current: "docs-aci-data-source-encryption_key"
description: |-
  Data source for ACI Encryption Key
---

# aci_encryption_key #
Data source for ACI Encryption Key


## API Information ##
* `Class` - pkiExportEncryptionKey
* `Distinguished Named` - uni/exportcryptkey

## GUI Information ##
* `Location` - 

## Example Usage ##
```hcl
data "aci_encryption_key" "example" {}
```

## Argument Reference ##


## Attribute Reference ##
* `id` - Attribute id set to the Dn of the AES Encryption Passphrase and Keys for Config Export (and Import).
* `annotation` - (Optional) Annotation of object AES Encryption Passphrase and Keys for Config Export (and Import).
* `name_alias` - (Optional) Name Alias of object AES Encryption Passphrase and Keys for Config Export (and Import).
* `clear_encryption_key` - (Optional) Pushbutton property to clear the encryption key, if configured. Setting this property to true will trigger the clearing of all fields in this mo,
             set the strongEncryptionEnabled policy to False and keyConfigured to False. There is no
             method to recover the previous passphrase before the clear operation.
* `passphrase` - (Optional) passphrase. The encryption parameters cannot be modified by a client request - only via a passphrase changeSetting this passphrase to blank/empty will trigger the clearing of all fields in this mo,
             set the strongEncryptionEnabled policy to False and keyConfigured to False. There is no
             method to recover the previous passphrase before the clear operation.
* `passphrase_key_derivation_version` - (Optional) passphraseKeyDerivationVersion. Version of the algorithm used - used for forward compatibility
* `strong_encryption_enabled` - (Optional) Strong Encryption Enabled for configuration export and import. Toggle to choose between weak and strong encryption - this flag can be set to True
           only when keyConfigured=True
