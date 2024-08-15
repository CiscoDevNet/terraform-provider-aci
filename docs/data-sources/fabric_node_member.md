---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_fabric_node_member"
sidebar_current: "docs-aci-data-source-aci_fabric_node_member"
description: |-
  Data source for ACI Fabric Node Member
---

# aci_fabric_node_member

Data source for ACI Fabric Node Member

## API Information ##

* `Class` - fabricNodeIdentP
* `Distinguished Name` - uni/controller/nodeidentpol/nodep-{serial}

## GUI Information ##

* `Location` -Fabric -> Inventory -> Fabric Mambership

## Example Usage

```hcl
data "aci_fabric_node_member" "example" {
  serial  = "example"
}
```

## Argument Reference

- `serial` - (Required) Serial Number for the new Fabric Node Member. Type: String.

## Attribute Reference

- `id` - (Read-Only) Attribute id set to the Dn of the Fabric Node Member. Type: String.
- `annotation` - (Read-Only) Specifies the annotation of a Fabric Node member. Type: String.
- `ext_pool_id` - (Read-Only) External pool ID for object Fabric Node member. Type: String.
- `fabric_id` - (Read-Only) Fabric ID for the new Fabric Node Member. Type: String.
- `name_alias` - (Read-Only) Name alias for object Fabric Node member. Type: String.
- `node_id` - (Read-Only) Node ID Number for the new Fabric Node Member. Type: String.
- `node_type` - (Read-Only) Node type for object Fabric Node member. Type: String.
- `pod_id` - (Read-Only) The pod id of the new Fabric Node Member. Type: String.
- `role` - (Read-Only) Role for the new Fabric Node Member. Type: String.
- `commission` - (Read-Only) Commission a node from the switch. Type: String.
