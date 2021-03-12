---
layout: "aci"
page_title: "ACI: aci_l3out_ospf_interface_profile"
sidebar_current: "docs-aci-resource-l3out_ospf_interface_profile"
description: |-
  Manages ACI L3out OSPF Interface Profile
---

# aci_l3out_ospf_interface_profile #
Manages ACI L3out OSPF Interface Profile

## Example Usage ##

```hcl
resource "aci_l3out_ospf_interface_profile" "example" {
  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
  annotation  = "example"
  auth_key  = "example"
  auth_key_id  = "example"
  auth_type  = "example"
  name_alias  = "example"
}
```


## Argument Reference ##

* `logical_interface_profile_dn` - (Required) distinguished name of parent logical interface profile object.
* `annotation` - (Optional) annotation for L3out OSPF interface profile object.
* `auth_key` - (Optional) ospf authentication key for L3out OSPF interface profile object.
* `auth_key_id` - (Optional) authentication key id for L3out OSPF interface profile object.
* `auth_type` - (Optional) ospf authentication type for L3out OSPF interface profile object.
* `name_alias` - (Optional) name_alias for L3out OSPF interface profile object.

* `relation_ospf_rs_if_pol` - (Optional) Relation to class ospfIfPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out OSPF Interface Profile.

## Importing ##

An existing L3out OSPF Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l3out_ospf_interface_profile.example <Dn>
```