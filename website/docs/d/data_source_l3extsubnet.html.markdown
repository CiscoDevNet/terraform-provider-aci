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
  ip                                    = "10.0.3.28/27"
}
```
## Argument Reference ##
* `external_network_instance_profile_dn` - (Required) Distinguished name of parent ExternalNetworkInstanceProfile object.
* `ip` - (Required) ip of Object subnet.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Subnet.
* `aggregate` - (Optional) Aggregate Routes for Subnet.
* `annotation` - (Optional) annotation for object subnet.
* `name_alias` - (Optional) name_alias for object subnet.
* `scope` - (Optional) The domain applicable to the capability.
