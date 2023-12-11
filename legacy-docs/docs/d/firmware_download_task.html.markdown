---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_firmware_download_task"
sidebar_current: "docs-aci-data-source-aci_firmware_download_task"
description: |-
  Data source for ACI Firmware Download Task
---

# aci_firmware_download_task

Data source for ACI Firmware Download Task

## Example Usage

```hcl
data "aci_firmware_download_task" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) The identifying name for the outside source of images, such as an HTTP or SCP server.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Firmware Download Task.
- `annotation` - (Optional) Annotation for the object of firmware download task.
- `description` - (Optional) Specifies the description of a policy component.
- `auth_pass` - (Optional) The authentication type for the source.
- `auth_type` - (Optional) The OSPF authentication type specifier.
- `dnld_task_flip` - (Optional) Download Task Flip flag.
- `identity_private_key_contents` - (Optional) Passphrase given at the identity key creation.
- `identity_private_key_passphrase` - (Optional) Passphrase given at the identity key creation.
- `identity_public_key_contents` - (Optional) Certificate contents for data transfer. Used for credentials.
- `load_catalog_if_exists_and_newer` - (Optional) Tracks to load the contained catalog or newer.
- `name_alias` - (Optional) Name alias for object firmware download task.
- `password` - (Optional) The Firmware password or key string.
- `polling_interval` - (Optional) Polling interval in minutes.
- `proto` - (Optional) The Firmware download protocol.
- `url` - (Optional) The firmware URL for the image(s) on the source.
- `user` - (Optional) The username for the source.
