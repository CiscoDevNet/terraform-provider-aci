---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_fabric_node_member"
sidebar_current: "docs-aci-resource-aci_fabric_node_member"
description: |-
  Manages ACI Fabric Node Member
---

# aci_fabric_node_member

Manages ACI Fabric Node Member

## API Information ##

* `Class` - fabricNodeIdentP
* `Distinguished Name` - uni/controller/nodeidentpol/nodep-{serial}

## GUI Information ##

* `Location` -Fabric -> Inventory -> Fabric Mambership

## Example Usage

```hcl
resource "aci_fabric_node_member" "example" {
  name        = "test"
  serial      = "1"
  annotation  = "example"
  description = "from terraform"
  ext_pool_id = "0"
  fabric_id   = "1"
  name_alias  = "example"
  node_id     = "201"
  node_type   = "unspecified"
  pod_id      = "1"
  role        = "unspecified"
  commission  = "yes"
}
```

## Argument Reference

- `serial` - (Required) Serial Number for the new Fabric Node Member. Type: String.
- `name` - (Required) Name of Fabric Node member. Type: String.
- `annotation` - (Optional) Specifies the annotation of a Fabric Node member. Type: String.
- `description` - (Optional) Specifies the description of a Fabric Node member. Type: String.
- `ext_pool_id` - (Optional) External pool ID for object Fabric Node member. Default value: "0". Type: String.
- `fabric_id` - (Optional) Fabric ID for the new Fabric Node Member. Default value: "1". Type: String.
- `name_alias` - (Optional) Name alias for object Fabric Node member. Type: String.
- `node_id` - (Optional) Node ID Number for the new Fabric Node Member. Allowed value range: "101" - "4000". Default value: "0". Type: String.
- `node_type` - (Optional) Node type for object Fabric Node member. Type: String.
  Allowed values: "unspecified", "remote-leaf-wan". Default value: "unspecified". Type: String.
- `pod_id` - (Optional) The pod id of the new Fabric Node Member. Allowed value range: "1" - "254". Default value: "1". Type: String.
- `role` - (Optional) Role for the new Fabric Node Member. Type: String.
  Allowed values: "unspecified", "leaf", "spine". Default value: "unspecified". Type: String.
- `commission` - (Optional) Commission a node from the switch. Type: String.
  Allowed values: "yes", "no". Default value: "yes".  Type: String. 
  - When commission is set to "no" the node will only be decommissioned and not removed from the controller.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fabric Node Member.

## Importing

An existing Fabric Node Member can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_fabric_node_member.example <Dn>
```
