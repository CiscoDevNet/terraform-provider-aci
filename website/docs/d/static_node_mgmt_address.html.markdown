---
layout: "aci"
page_title: "ACI: aci_static_node_mgmt_address"
sidebar_current: "docs-aci-data-source-static_node_mgmt_address"
description: |-
  Data source for ACI Management Static Node

---

# aci_static_node_mgmt_address #
Data source for ACI Management Static Node

## Example Usage ##

```hcl
data "aci_static_node_mgmt_address" "example" {
  management_epg_dn = "${aci_node_mgmt_epg.example.id}"
  t_dn              = "${aci_static_node_mgmt_address.example.t_dn}"
  type              = "out_of_band"
}
```


## Argument Reference ##

* `management_epg_dn` - (Required) distinguished name of parent management static node object.
* `t_dn` - (Required) target dn of management static node object.
* `type` - (Required) type for management static node object.
Note := for "in_band", `management_epg_dn` should be of type "in_band" and for "out_of_band", `management_epg_dn` should be of type "out_of_band".



## Attribute Reference

* `id` - attribute id set to the Dn of management static node object.
* `addr` - peer address for management static node object.
* `annotation` - annotation for management static node object.
* `gw` - gateway IP address for management static node object
* `v6_addr` - v6 address for management static node object.
* `v6_gw` - v6 gw for management static node object.
