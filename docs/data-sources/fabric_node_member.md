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

## Example Usage

```hcl
data "aci_fabric_node_member" "example" {
  serial  = "example"
}
```

## Argument Reference

- `serial` - (Required) serial of Object fabric_node_member.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Fabric Node Member.
- `annotation` - (Optional) annotation for object fabric_node_member.
- `ext_pool_id` - (Optional) ext_pool_id for object fabric_node_member.
- `fabric_id` - (Optional) place holder for a value
- `name_alias` - (Optional) name_alias for object fabric_node_member.
- `node_id` - (Optional) node id
- `node_type` - (Optional) node_type for object fabric_node_member.
- `pod_id` - (Optional) pod id
- `role` - (Optional) system role type
- `serial` - (Optional) serial number
