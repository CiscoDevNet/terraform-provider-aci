---
layout: "aci"
page_title: "ACI: aci_l3out_loopback_interface_profile"
sidebar_current: "docs-aci-data-source-l3out_loopback_interface_profile"
description: |-
  Data source for ACI Loop Back Interface Profile
---

# aci_l3out_loopback_interface_profile #
Data source for ACI Loop Back Interface Profile

## Example Usage ##

```hcl
data "aci_l3out_loopback_interface_profile" "example" {
  fabric_node_dn = aci_logical_node_to_fabric_node.example.id
  addr           = "1.2.3.5"
}
```

## Argument Reference ##

* `fabric_node_dn` - (Required) Distinguished name of parent Fabric Node object.
* `addr` - (Required) Address of L3out lookback interface profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L3out lookback interface profile.
* `description` - Description for L3out lookback interface profile.
* `annotation` - Annotation for L3out lookback interface profile.
* `name_alias` - Name alias for L3out lookback interface profile.
