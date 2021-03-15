---
layout: "aci"
page_title: "ACI: aci_l3out_ospf_interface_profile"
sidebar_current: "docs-aci-data-source-l3out_ospf_interface_profile"
description: |-
  Data source for ACI L3out OSPF Interface Profile
---

# aci_l3out_ospf_interface_profile #
Data source for ACI L3out OSPF Interface Profile

## Example Usage ##

```hcl
data "aci_l3out_ospf_interface_profile" "example" {
  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
}
```


## Argument Reference ##

* `logical_interface_profile_dn` - (Required) distinguished name of parent logical interface profile object.



## Attribute Reference

* `id` - attribute id set to the Dn of the Interface Profile.
* `annotation` - annotation for L3out OSPF interface profile object.
* `description` - description for L3out OSPF interface profile object.
* `auth_key_id` - authentication key id for L3out OSPF interface profile object.
* `auth_type` - ospf authentication type for L3out OSPF interface profile object.
* `name_alias` - name_alias for L3out OSPF interface profile object.
