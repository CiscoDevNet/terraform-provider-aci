---
layout: "aci"
page_title: "ACI: aci_l3out_loopback_interface_profile"
sidebar_current: "docs-aci-resource-l3out_loopback_interface_profile"
description: |-
  Manages ACI L3out Loopback Interface Profile
---

# aci_l3out_loopback_interface_profile #
Manages ACI L3out Loopback Interface Profile

## Example Usage ##

```hcl
resource "aci_l3out_loopback_interface_profile" "example" {
  fabric_node_dn = "${aci_logical_node_to_fabric_node.example.id}"
  addr           = "1.2.3.5"
  description    = "from terraform"
  annotation     = "example"
  name_alias     = "example"
}
```


## Argument Reference ##

* `fabric_node_dn` - (Required) Distinguished name of parent fabric node object.
* `addr` - (Required) Address of L3out lookback interface profile.
* `description` - (Optional) Description for L3out lookback interface profile.
* `annotation` - (Optional) Annotation for L3out lookback interface profile.
* `name_alias` - (Optional) Name alias for L3out lookback interface profile.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out Loopback Interface Profile.

## Importing ##

An existing L3out Loopback Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l3out_loopback_interface_profile.example <Dn>
```