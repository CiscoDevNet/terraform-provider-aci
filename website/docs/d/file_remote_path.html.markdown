---
subcategory: "Import/Export"
layout: "aci"
page_title: "ACI: aci_file_remote_path"
sidebar_current: "docs-aci-data-source-file_remote_path"
description: |-
  Data source for ACI Remote Path of a File
---

# aci_file_remote_path #
Data source for ACI Remote Path of a File

## API Information ##
* `Class` - fileRemotePath
* `Distinguished Name` - uni/fabric/path-{name}

## GUI Information ##
* `Location` - Admin -> Import/Export -> Remote Locations -> Create Remote Location

## Example Usage ##

```hcl
data "aci_file_remote_path" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object File Remote Path.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the File Remote Path.
* `annotation` - (Optional) Annotation of object File Remote Path.
* `name_alias` - (Optional) Name Alias of object File Remote Path.
* `auth_type` - (Optional) Authentication Type Choice.
* `host` - (Optional) Hostname or IP for export destination of object File Remote Path.
* `protocol` - (Optional) Transfer protocol to be used for data export of object File Remote Path.
* `remote_path` - (Optional) Path where data will reside in the destination of object File Remote Path.
* `remote_port` - (Optional) Remote port for data export destination of object File Remote Path.
* `user_name` - (Optional) Username to be used to transfer data to destination of object File Remote Path.
* `description` - (Optional) Description of object File Remote Path.

