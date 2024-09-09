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

- `serial` - (Required) Serial Number of the Fabric Node Member. Type: String.

## Attribute Reference

- `id` - (Read-Only) Attribute ID set to the DN of the Fabric Node Member. Type: String.
- `annotation` - (Read-Only) Specifies the annotation of a Fabric Node Member. Type: String.
- `ext_pool_id` - (Read-Only) External pool ID of object Fabric Node Member. Type: String.
- `fabric_id` - (Read-Only) Fabric ID of the Fabric Node Member. Type: String.
- `name_alias` - (Read-Only) Name alias of object Fabric Node Member. Type: String.
- `node_id` - (Read-Only) Node ID Number of the Fabric Node Member. Type: String.
- `node_type` - (Read-Only) Node type of object Fabric Node Member. Type: String.
- `pod_id` - (Read-Only) Pod ID of the Fabric Node Member. Type: String.
- `role` - (Read-Only) Role of the Fabric Node Member. Type: String.
- `commission` - (Read-Only) Commission a node from the switch to make it an active member of the fabric. Type: String.