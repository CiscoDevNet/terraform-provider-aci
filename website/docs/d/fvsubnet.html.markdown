---
layout: "aci"
page_title: "ACI: aci_subnet"
sidebar_current: "docs-aci-data-source-subnet"
description: |-
  Data source for ACI Subnet
---

# aci_subnet #
Data source for ACI Subnet

## Example Usage ##

```hcl

data "aci_subnet" "dev_subnet" {
  parent_dn         = "${aci_bridge_domain.example.id}"
  ip                = "10.0.3.28/27"
}

```


## Argument Reference ##
* `parent_dn` - (Required) Distinguished name of parent object.
* `ip` - (Required) The IP address and mask of the default gateway.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Subnet.
* `annotation` - (Optional) annotation for object subnet.
* `ctrl` - (Optional) The subnet control state. The control can be specific protocols applied to the subnet such as IGMP Snooping.
* `name_alias` - (Optional) name_alias for object subnet.
* `preferred` - (Optional) Indicates if the subnet is preferred (primary) over the available alternatives. Only one preferred subnet is allowed.
* `scope` - (Optional) The network visibility of the subnet.
* `virtual` - (Optional) Treated as virtual IP address. Used in case of BD extended to multiple sites.