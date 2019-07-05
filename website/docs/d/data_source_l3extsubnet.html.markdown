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
data "aci_l3_ext_subnet" "example" {

  external_network_instance_profile_dn  = "${aci_external_network_instance_profile.example.id}"

  ip  = "example"
}
```
## Argument Reference ##
* `external_network_instance_profile_dn` - (Required) Distinguished name of parent ExternalNetworkInstanceProfile object.
* `ip` - (Required) ip of Object subnet.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Subnet.
* `aggregate` - (Optional) aggregate for object subnet.
* `annotation` - (Optional) annotation for object subnet.
* `ip` - (Optional) ip address
* `name_alias` - (Optional) name_alias for object subnet.
* `scope` - (Optional) capability domain
