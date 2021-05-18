---
layout: "aci"
page_title: "ACI: aci_l3out_static_route"
sidebar_current: "docs-aci-data-source-l3out_static_route"
description: |-
  Data source for ACI L3out Static Route
---

# aci_l3out_static_route

Data source for ACI L3out Static Route

## Example Usage

```hcl
data "aci_l3out_static_route" "example" {
  fabric_node_dn  = "${aci_logical_node_to_fabric_node.example.id}"
  ip  = "example"
}
```

## Argument Reference

- `fabric_node_dn` - (Required) Distinguished name of parent fabric node object.
- `ip` - (Required) The static route IP address assigned to the outside network.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out Static Route.
- `aggregate` - (Optional) Aggregated Route for object l3out static route.
- `annotation` - (Optional) Annotation for object l3out static route.
- `description` - (Optional) Description for object l3out static route.
- `name_alias` - (Optional) Name alias for object l3out static route.
- `pref` - (Optional) The administrative preference value for this route. This value is useful for resolving routes advertised from different protocols.
- `rt_ctrl` - (Optional) Route control for object l3out static route.
