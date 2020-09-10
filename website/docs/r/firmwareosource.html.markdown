---
layout: "aci"
page_title: "ACI: aci_firmware_download_task"
sidebar_current: "docs-aci-resource-firmware_download_task"
description: |-
  Manages ACI Firmware Download Task
---

# aci_firmware_download_task #
Manages ACI Firmware Download Task

## Example Usage ##

```hcl
resource "aci_firmware_download_task" "example" {


  name  = "example"
  annotation  = "example"
  auth_pass  = "example"
  auth_type  = "example"
  dnld_task_flip  = "example"
  identity_private_key_contents  = "example"
  identity_private_key_passphrase  = "example"
  identity_public_key_contents  = "example"
  load_catalog_if_exists_and_newer  = "example"
  name_alias  = "example"
  password  = "example"
  polling_interval  = "example"
  proto  = "example"
  url  = "example"
  user  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_download_task.
* `annotation` - (Optional) annotation for object firmware_download_task.
* `auth_pass` - (Optional) authentication type.
Allowed values: "password", "key"
* `auth_type` - (Optional) ospf authentication type specifier.
Allowed values: "usePassword", "useSshKeyContents"
* `dnld_task_flip` - (Optional) dnld_task_flip for object firmware_download_task.
* `identity_private_key_contents` - (Optional) identity_private_key_contents for object firmware_download_task.
* `identity_private_key_passphrase` - (Optional) identity_private_key_passphrase for object firmware_download_task.
* `identity_public_key_contents` - (Optional) identity_public_key_contents for object firmware_download_task.
* `load_catalog_if_exists_and_newer` - (Optional) tracks to load the contained catalog or newer.
Allowed values: "yes", "no"
* `name_alias` - (Optional) name_alias for object firmware_download_task.
* `password` - (Optional) password/key string
* `polling_interval` - (Optional) polling interval
* `proto` - (Optional) download protocol.
Allowed values: "scp", "http", "usbkey", "local"
* `url` - (Optional) URL of image of source
* `user` - (Optional) username for source



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Firmware Download Task.

## Importing ##

An existing Firmware Download Task can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_firmware_download_task.example <Dn>
```