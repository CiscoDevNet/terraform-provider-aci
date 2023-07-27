---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_bfd_multihop_interface_profile"
sidebar_current: "docs-aci-data-source-aci_bfd_multihop_interface_profile"
description: |-
  Data source for ACI BFD Multihop Interface Profile
---

# aci_bfd_multihop_interface_profile #

Data source for ACI BFD Multihop Interface Profile


## API Information ##

* `Class` - bfdMhIfP
* `Distinguished Name` - uni/tn-{tn_name}/out-{l3out_name}/lnodep-{ln_name}/lifp-{lifp_name}/bfdMhIfP

## GUI Information ##

* `Location` -  Tenant -> Networking -> L3Out -> Logical Node Profiles -> Logical Interface Profiles 



## Example Usage ##

```hcl
data "aci_bfd_multihop_interface_profile" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of parent Logical Interface Profile object.  Type: String.

## Attribute Reference ##
* `id`                           - (Read-Only) Attribute id set to the Dn of the BFD Multihop Interface Profile. Type: String.
* `annotation`                   - (Read-Only) Annotation of the BFD Multihop Interface Profile object. Type: String.
* `name_alias`                   - (Read-Only) Name Alias of the BFD Multihop Interface Profile object. Type: String.
* `name`                         - (Read-Only) Name of the BFD Multihop Interface Profile object. Type: String.
* `key`                          - (Read-Only) Authentication Key. Type: String.
* `key_id`                       - (Read-Only) Authentication Key ID. Type: String.
* `interface_profile_type`       - (Read-Only) Authentication Type. Type: String.
* `relation_bfd_rs_mh_if_pol`    - (Read-Only) Represents the relation to the BFD interface policy (class bfdMhIfPol). Type: String.
