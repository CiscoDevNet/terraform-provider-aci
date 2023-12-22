---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_firmware_download_task"
sidebar_current: "docs-aci-resource-aci_firmware_download_task"
description: |-
  Manages ACI Firmware Download Task
---

# aci_firmware_download_task

Manages ACI Firmware Download Task

## Example Usage

```hcl
resource "aci_firmware_download_task" "example" {
  name                             = "example"
  annotation                       = "example"
  auth_pass                        = "password"
  auth_type                        = "usePassword"
  dnld_task_flip                   = "yes"
  identity_private_key_contents    = "example"
  identity_private_key_passphrase  = "example"
  identity_public_key_contents     = "example"
  load_catalog_if_exists_and_newer = "yes"
  name_alias                       = "example"
  password                         = "SomeSecretPassword"
  polling_interval                 = "20"
  proto                            = "http"
  url                              = "foo.bar.cisco.com/download/cisco/aci/aci-msft-pkg-3.1.1i.zip"
  user                             = "admin"
}
```

## Argument Reference

- `name` - (Required) The identifying name for the outside source of images, such as an HTTP or SCP server.
- `annotation` - (Optional) Annotation for the object of firmware download task.
- `description` - (Optional) Specifies the description of a policy component.
- `auth_pass` - (Optional) The authentication type for the source.
  Allowed values: "password", "key". Default value: "password".
- `auth_type` - (Optional) The OSPF authentication type specifier.
  Allowed values: "usePassword", "useSshKeyContents". Default value: "usePassword".
- `dnld_task_flip` - (Optional) Download Task Flip flag.
  Allowed values: "yes", "no". Default value: "yes".
- `identity_private_key_contents` - (Optional) Passphrase given at the identity key creation.
- `identity_private_key_passphrase` - (Optional) Passphrase given at the identity key creation.
- `identity_public_key_contents` - (Optional) Certificate contents for data transfer. Used for credentials.
- `load_catalog_if_exists_and_newer` - (Optional) Tracks to load the contained catalog or newer.
  Allowed values: "yes", "no". Default value: "yes".
- `name_alias` - (Optional) Name alias for object firmware download task.
- `password` - (Optional) The Firmware password or key string.
- `polling_interval` - (Optional) Polling interval in minutes.
- `proto` - (Optional) The Firmware download protocol. Allowed values: "scp", "http", "usbkey", "local". Default values: "scp".
- `url` - (Optional) The firmware URL for the image(s) on the source.
- `user` - (Optional) The username for the source.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Firmware Download Task.

## Importing

An existing Firmware Download Task can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_firmware_download_task.example <Dn>
```
