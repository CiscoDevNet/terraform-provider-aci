---
layout: "aci"
page_title: "ACI: aci_firmware_download_task"
sidebar_current: "docs-aci-data-source-firmware_download_task"
description: |-
  Data source for ACI Firmware Download Task
---

# aci_firmware_download_task #
Data source for ACI Firmware Download Task

## Example Usage ##

```hcl
data "aci_firmware_download_task" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object firmware_download_task.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Firmware Download Task.
* `annotation` - (Optional) annotation for object firmware_download_task.
* `auth_pass` - (Optional) authentication type
* `auth_type` - (Optional) ospf authentication type specifier
* `dnld_task_flip` - (Optional) dnld_task_flip for object firmware_download_task.
* `identity_private_key_contents` - (Optional) identity_private_key_contents for object firmware_download_task.
* `identity_private_key_passphrase` - (Optional) identity_private_key_passphrase for object firmware_download_task.
* `identity_public_key_contents` - (Optional) identity_public_key_contents for object firmware_download_task.
* `load_catalog_if_exists_and_newer` - (Optional) tracks to load the contained catalog or newer
* `name_alias` - (Optional) name_alias for object firmware_download_task.
* `password` - (Optional) password/key string
* `polling_interval` - (Optional) polling interval
* `proto` - (Optional) download protocol
* `url` - (Optional) URL of image of source
* `user` - (Optional) username for source
