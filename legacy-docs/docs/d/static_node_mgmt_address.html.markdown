---
subcategory: "Node Management"
layout: "aci"
page_title: "ACI: aci_static_node_mgmt_address"
sidebar_current: "docs-aci-data-source-aci_static_node_mgmt_address"
description: |-
  Data source for ACI Management Static Node

---

# aci_static_node_mgmt_address #
Data source for ACI Management Static Node

## Example Usage ##

```hcl
data "aci_static_node_mgmt_address" "example" {
  management_epg_dn = aci_node_mgmt_epg.example.id
  t_dn              = aci_static_node_mgmt_address.example.t_dn
  type              = "out_of_band"
}
```


## Argument Reference ##

* `management_epg_dn` - (Required) Distinguished name of parent Management static node object.
* `t_dn` - (Required) Target dn of Management static node object.
* `type` - (Required) Type for Management static node object. Allowed values are "in_band" and "out_of_band".
Note := for "in_band", `management_epg_dn` should be of type "in_band" and for "out_of_band", `management_epg_dn` should be of type "out_of_band".



## Attribute Reference

* `id` - Attribute id set to the Dn of Management static node object.
* `addr` - Peer address for Management static node object.
* `annotation` - Annotation for Management static node object.
* `description` - Description for Management static node object.
* `gw` - Gateway IP address for Management static node object.
* `v6_addr` - V6 address for Management static node object.
* `v6_gw` - V6 gw for Management static node object.
