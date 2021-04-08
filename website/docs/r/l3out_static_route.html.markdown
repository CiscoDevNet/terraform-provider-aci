---
layout: "aci"
page_title: "ACI: aci_l3out_static_route"
sidebar_current: "docs-aci-resource-l3out_static_route"
description: |-
  Manages ACI L3out Static Route
---

# aci_l3out_static_route

Manages ACI L3out Static Route

## Example Usage

```hcl
resource "aci_l3out_static_route" "example" {

  fabric_node_dn  = "${aci_logical_node_to_fabric_node.example.id}"
  ip  = "10.0.0.1"
  aggregate = "no"
  annotation  = "example"
  name_alias  = "example"
  pref  = "example"
  rt_ctrl = "bfd"

}
```

## Argument Reference

- `fabric_node_dn` - (Required) Distinguished name of parent fabric node object.
- `ip` - (Required) The static route IP address assigned to the outside network.
- `aggregate` - (Optional) Aggregated Route for object l3out static route.
  Allowed values: "no", "yes". Default value: "no".
- `annotation` - (Optional) Annotation for object l3out static route.
- `description` - (Optional) Description for object l3out static route.
- `name_alias` - (Optional) Name alias for object l3out static route.
- `pref` - (Optional) The administrative preference value for this route. This value is useful for resolving routes advertised from different protocols. Default value: "1".
- `rt_ctrl` - (Optional) Route control for object l3out static route.
  Allowed values: "bfd", "unspecified". Default value: "unspecified".

- `relation_ip_rs_route_track` - (Optional) Relation to class fvTrackList. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out Static Route.

## Importing0

An existing L3out Static Route can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_static_route.example <Dn>
```
