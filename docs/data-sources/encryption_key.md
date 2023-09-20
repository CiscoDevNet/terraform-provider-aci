---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_encryption_key"
sidebar_current: "docs-aci-data-source-encryption_key"
description: |-
  Data source for ACI AES Encryption Passphrase and Keys for Config Export and Import
---

# aci_encryption_key #
Data source for ACI AES Encryption Passphrase and Keys for Config Export and Import

## API Information ##
* `Class` - pkiExportEncryptionKey
* `Distinguished Name` - uni/exportcryptkey

## GUI Information ##
* `Location` - System -> System Settings -> Global AES Passphrase Encryption Settings -> Policy

## Example Usage ##
```hcl
data "aci_encryption_key" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the AES Encryption Passphrase and Keys for Config Export (and Import).
* `annotation` - (Optional) Annotation of object AES Encryption Passphrase and Keys for Config Export and Import.
* `name_alias` - (Optional) Name Alias of object AES Encryption Passphrase and Keys for Config Export and Import.
* `description` - (Optional) Description of object AES Encryption Passphrase and Keys for Config Export and Import.
* `passphrase_key_derivation_version` - (Optional) Version of the algorithm used for forward compatibility.
* `strong_encryption_enabled` - (Optional) Parameter indicating whether encryption is weak or strong. 
