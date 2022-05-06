---
subcategory: "Import/Export"
layout: "aci"
page_title: "ACI: aci_file_remote_path"
sidebar_current: "docs-aci-resource-file_remote_path"
description: |-
  Manages ACI Remote Path of a File
---

# aci_file_remote_path #
Manages ACI Remote Path of a File

## API Information ##
* `Class` - fileRemotePath
* `Distinguished Name` - uni/fabric/path-{name}

## GUI Information ##
* `Location` - Admin -> Import/Export -> Remote Locations -> Create Remote Location

## Example Usage ##
```hcl
resource "aci_file_remote_path" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  auth_type = "usePassword"
  host = "cisco.com"
  protocol = "sftp"
  remote_path = "/example_remote_path"
  remote_port = "0"
  user_name = "example_user_name"
  user_passwd = "password"
  name_alias = "example_name_alias"
  description = "from terraform"
}
```

## Argument Reference ##
* `name` - (Required) Name of object File Remote Path.
* `host` - (Required) Hostname or IP for export destination of object File Remote Path.
* `name_alias` - (Optional) Name alias for object File Remote Path.
* `annotation` - (Optional) Annotation of object File Remote Path.
* `auth_type` - (Optional) Authentication Type Choice. Allowed values are "usePassword" and "useSshKeyContents". Default value is "usePassword". Type: String.
* `identity_private_key_contents` - (Optional) SSH Private Key File contents for datatransfer. Must be set if `auth_type` is equal to "useSshKeyContents".
* `identity_private_key_passphrase` - (Optional)  Passphrase given at the identity key creation. Should be set if and only if `identity_private_key_contents` is set.
* `protocol` - (Optional) Transfer protocol to be used for data export of object File Remote Path .Allowed values are "ftp", "scp" and "sftp". Default value is "sftp". Type: String. Value "ftp" cannot be set if `auth_type` is equal to "useSshKeyContents".
* `remote_path` - (Optional) Path where data will reside in the destination of object File Remote Path(The first character of remote_path should be '/').
* `remote_port` - (Optional) Remote port for data export destination of object File Remote Path. Range: "0" - "65535". Default value is "0".
* `user_name` - (Optional) Username to be used to transfer data to destination of object File Remote Path.
* `user_passwd` - (Optional) Password to be used to transfer data to destination of object File Remote Path. Must be set if `auth_type` is equal to "usePassword".
* `relation_file_rs_a_remote_host_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the remote host is reachable. Type: String.
* `relation_file_rs_a_remote_host_to_epp` - (Optional) Represents the relation to a Relation to Remote Host  Reachability EPP (class fvAREpP). A source relation to the abstract representation of the resolvable endpoint profile. Type: String.
* `description` - (Optional) Description of object File Remote Path.

## Importing ##

An existing File Remote Path can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_file_remote_path.example <Dn>
```