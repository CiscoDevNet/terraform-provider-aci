---
layout: "aci"
page_title: "ACI: aci_tacacs_accounting_destination"
sidebar_current: "docs-aci-resource-tacacs_accounting_destination"
description: |-
  Manages ACI TACACS Accounting Destination
---

# aci_tacacs_accounting_destination #
Manages ACI TACACS Accounting Destination

## API Information ##
* `Class` - tacacsTacacsDest
* `Distinguished Named` - uni/fabric/tacacsgroup-{name}/tacacsdest-{host}-port-{port}

## GUI Information ##
* `Location` - Admin -> External Data Collectors -> Monitoring Destinations -> TACACS -> TACACS Destinations


## Example Usage ##
```hcl
resource "aci_tacacs_accounting_destination" "example" {
  tacacs_monitoring_destination_group_dn  = aci_tacacs_accounting.example.id
  host = "cisco.com"
  port = "49"
  annotation = "orchestrator:terraform"
  auth_protocol = "pap"
  key = "example_key_value"
  description = "from terraform"
}
```

## Argument Reference ##
* `tacacs_accounting_dn` - (Required) Distinguished name of parent TACACS Accounting object..
* `host` - (Required) Host or IP address of object TACACS Accounting Destination.
* `port` - (Required) Port of object TACACS Accounting Destination. Allowed Range: "1" - "65535". Default value: "49".
* `annotation` - (Optional) Annotation of object TACACS Accounting Destination.
* `name_alias` - (Optional) Name Alias of object TACACS Accounting Destination.
* `name` - (Optional) Name of object TACACS Accounting Destination.
* `auth_protocol` - (Optional) Authentication Protocol of object TACACS Accounting Destination. Allowed values are "chap", "mschap" and "pap". Default value is "pap". Type: String.
* `key` - (Optional) The key or password used to uniquely identify object TACACS Accounting Destination.
* `description` - (Optional) Description of object TACACS Accounting Destination.
* `relation_file_rs_a_remote_host_to_epg` - (Optional) Represents the relation to a Attachable Target Group (class fvATg). A source relation to the endpoint group through which the remote host is reachable. Type: String.
* `relation_file_rs_a_remote_host_to_epp` - (Optional) Represents the relation to a Relation to Remote Host  Reachability EPP (class fvAREpP). A source relation to the abstract representation of the resolvable endpoint profile. Type: String.


## Importing ##
An existing TACACS Accounting Destination can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tacacs_accounting_destination.example <Dn>
```