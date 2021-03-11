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
  out-_of-_band_management_e_pg_dn  = "${aci_out-_of-_band_management_e_pg.example.id}"
  tDn  = "example"
  type = "in-band"
}
```


## Argument Reference ##

* `out-_of-_band_management_e_pg_dn` - (Required) Distinguished name of parent management static node object.
* `tDn` - (Required) tDn of management static node object.
* `type` - (Required) type for management static node object.



## Attribute Reference

* `id` - Attribute id set to the Dn of management static node object.
* `addr` - (Optional) peer address for management static node object.
* `annotation` - (Optional) annotation for management static node object.
* `gw` - (Optional) gateway IP address for management static node object
* `v6_addr` - (Optional) v6 address for management static node object.
* `v6_gw` - (Optional) v6 gw for management static node object.
