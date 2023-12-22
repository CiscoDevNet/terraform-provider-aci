---
subcategory: "System Settings"
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
* `Distinguished Name` - uni/exportcryptkey

## GUI Information ##
* `Location` - System -> System Settings -> Global AES Passphrase Encryption Settings -> Policy

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
* Users can use the resource of type `aci_encryption_key` to change the configuration of the object AES Encryption Passphrase and Keys for Config Export and Import. Users cannot create more than one instance of object AES Encryption Passphrase and Keys for Config Export and Import.
* If `passphrase` is set then `clear_encryption_key` should be `no`. Also, `clear_encryption_key` and `strong_encryption_enabled` both should not be set `yes` simultaneously.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object AES Encryption Passphrase and Keys for Config Export and Import.
* `passphrase` - (Optional) Parameter to set the passphrase of object AES Encryption Passphrase and Keys for Config Export and Import. Length of `passphrase` should be between 16 - 32 characters. 
* `strong_encryption_enabled` - (Optional) Parameter indicating whether encryption is weak or strong. This parameter can be set if and only if the `passphrase` is set. Allowed values are "yes" and "no". Type: String. 
* `clear_encryption_key` - (Optional) Parameter to clear the encryption key, if configured. Allowed values are "yes" and "no". Default value is "no". Type: String.
* `passphrase_key_derivation_version` - (Optional) Version of the algorithm used for forward compatibility. Allowed value is "v1". Default value is "v1".
* `description` - (Optional) Description of object AES Encryption Passphrase and Keys for Config Export and Import.
* `name_alias` - (Optional) Name Alias of object AES Encryption Passphrase and Keys for Config Export and Import.


## Importing ##

An existing AES Encryption Passphrase and Keys for Config Export and Import can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_encryption_key.example <Dn>
```
