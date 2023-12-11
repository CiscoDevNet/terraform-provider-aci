---
subcategory: "Node Management"
layout: "aci"
page_title: "ACI: aci_static_node_mgmt_address"
sidebar_current: "docs-aci-resource-aci_static_node_mgmt_address"
description: |-
  Manages ACI Management Static Node
---

# aci_static_node_mgmt_address

Manages ACI Management Static Node

## Example Usage

```hcl
resource "aci_static_node_mgmt_address" "example" {
  management_epg_dn = aci_node_mgmt_epg.example.id
  t_dn              = "topology/pod-1/node-1"
  type              = "out_of_band"
  description       = "from terraform"
  addr              = "10.20.30.40/20"
  annotation        = "example"
  gw                = "10.20.30.41"
  v6_addr           = "1::40/64"
  v6_gw             = "1::21"
}
```

## Argument Reference

- `management_epg_dn` - (Required) Distinguished name of parent Management static node object.
- `t_dn` - (Required) Target dn of Management static node object.
- `type` - (Required) Type of the Management static node object. Allowed values are "in_band" and "out_of_band".
  <strong>Note</strong> : for "in_band", `management_epg_dn` should be of type "in_band" and for "out_of_band", `management_epg_dn` should be of type "out_of_band".
- `addr` - (Optional) Peer address of the Management static node object. Default value: "0.0.0.0".
- `annotation` - (Optional) Annotation for Management static node object.
- `description` - (Optional) Description for Management static node object.
- `gw` - (Optional) Gateway IP address for Management static node object. Default value: "0.0.0.0".
- `v6_addr` - (Optional) V6 address for Management static node object. Default value: "::".
- `v6_gw` - (Optional) V6 gw for Management static node object. Default value: "::".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Management Static Node.

## Importing

An existing Management Static Node can be [imported][docs-import] into this resource via its Dn and type, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_static_node_mgmt_address.example <Dn>:<type>
```
