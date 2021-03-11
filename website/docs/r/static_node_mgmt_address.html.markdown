---
layout: "aci"
page_title: "ACI: aci_static_node_mgmt_address"
sidebar_current: "docs-aci-resource-static_node_mgmt_address"
description: |-
  Manages ACI Management Static Node
---

# aci_static_node_mgmt_address #
Manages ACI Management Static Node

## Example Usage ##

```hcl
resource "aci_static_node_mgmt_address" "example" {
  management_epg_dn  = "${aci_node_mgmt_epg.example.id}"
  tDn  = "example"
  type = "in-band"
  addr  = "example"
  annotation  = "example"
  gw  = "example"
  v6_addr  = "example"
  v6_gw  = "example"
}
```


## Argument Reference ##

* `management_epg_dn` - (Required) Distinguished name of parent management static node object.
* `tDn` - (Required) tDn of management static node object.
* `type` - (Required) type of the management static node object.
* `addr` - (Optional) peer address of the management static node object
* `annotation` - (Optional) annotation for management static node object.
* `gw` - (Optional) gateway IP address for management static node object
* `v6_addr` - (Optional) v6 address for management static node object.
* `v6_gw` - (Optional) v6 gw for management static node object.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Management Static Node.

## Importing ##

An existing Management Static Node can be [imported][docs-import] into this resource via its Dn and type, using the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_static_node_mgmt_address.example <Dn>:<type>
```