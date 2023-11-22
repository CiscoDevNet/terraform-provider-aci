---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_bfd_multihop_interface_profile"
sidebar_current: "docs-aci-resource-bfd_multihop_interface_profile"
description: |-
  Manages ACI BFD Multihop Interface Profile
---

# aci_bfd_multihop_interface_profile #

Manages ACI BFD Multihop Interface Profile

## API Information ##

* `Class` - bfdMhIfP
* `Distinguished Name` - uni/tn-{tn_name}/out-{l3out_name}/lnodep-{ln_name}/lifp-{lifp_name}/bfdMhIfP

## GUI Information ##

* `Location` -  Tenant -> Networking -> L3Out -> Logical Node Profiles -> Logical Interface Profiles 


## Example Usage ##

```hcl
resource "aci_bfd_multihop_interface_profile" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  key                          = "SomeAuthKey"
  key_id                       = "1"
  interface_profile_type       = "none"
  bfd_rs_mh_if_pol             = aci_resource.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent BFD Multihop Logical Interface Profile object. Type: String.
* `annotation`                   - (Optional) Annotation of the BFD Multihop Interface Profile object. Type: String.
* `name_alias`                   - (Optional) Name Alias of the BFD Multihop Interface Profile object. Type: String.
* `name`                         - (Optional) Name of the BFD Multihop Interface Profile object. Type: String.
* `key`                          - (Optional) Authentication Key. Type: String.
* `key_id`                       - (Optional) Authentication Key ID. Allowed range is 1-255 and default value is "1". Type: String.
* `interface_profile_type`       - (Optional) Authentication Type. Allowed values are "none", "sha1", and default value is "none". Type: String.
* `relation_bfd_rs_mh_if_pol`    - (Optional) Represents the relation to the BFD interface policy (class bfdMhIfPol). Type: String.

## Importing ##

An existing BFD Multihop Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bfd_multihop_interface_profile.example <Dn>
```