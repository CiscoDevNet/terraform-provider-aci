---
layout: "aci"
page_title: "ACI: aci_encryption_key"
sidebar_current: "docs-aci-resource-aci_encryption_key"
description: |-
  Manages ACI AES Encryption Passphrase and Keys for Config Export and Import
---

# aci_encryption_key #
Manages ACI AES Encryption Passphrase and Keys for Config Export and Import

## API Information ##
* `Class` - pkiExportEncryptionKey
* `Distinguished Named` - uni/exportcryptkey

## GUI Information ##
* `Location` - System -> System Settings -> Global AES Encryption Settings for all Configuration Import and Export

## Example Usage ##
```hcl
resource "aci_encryption_key" "example" {
  description = "from terraform"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
  clear_encryption_key = "no"
  passphrase = "example_passphrase"
  passphrase_key_derivation_version = "v1"
  strong_encryption_enabled = "yes"
}
```

## NOTE ##
User can use resource of type aci_encryption_key to change configuration of object AES Encryption Passphrase and Keys for Config Export and Import. User cannot create more than one instances of object AES Encryption Passphrase and Keys for Config Export and Import.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object AES Encryption Passphrase and Keys for Config Export and Import.
* `clear_encryption_key` - (Optional) Parameter to clear the encryption key, if configured. Allowed values are "yes" and "no". Type: String.
* `passphrase` - (Optional) Parameter to set the passphrase of object AES Encryption Passphrase and Keys for Config Export and Import
* `passphrase_key_derivation_version` - (Optional) Version of the algorithm used for forward compatibility. Allowed value is "v1". Default value is "v1".
* `strong_encryption_enabled` - (Optional) Parameter indicating whether encryption is weak or strong. This parameter can be set if and only if `passphrase` is set. Allowed values are "yes" and "no". Type: String.
* `description` - (Optional) Description of object AES Encryption Passphrase and Keys for Config Export and Import.
* `name_alias` - (Optional) Name Alias of object AES Encryption Passphrase and Keys for Config Export and Import.


## Importing ##

An existing AES Encryption Passphrase and Keys for Config Export and Import can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_encryption_key.example <Dn>
```