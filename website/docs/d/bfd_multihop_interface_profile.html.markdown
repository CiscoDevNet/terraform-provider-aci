---
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

* `logical_interface_profile_dn` - (Required) Distinguished name of parent LogicalInterfaceProfile object.

## Attribute Reference ##
* `id`                     - Attribute id set to the Dn of the BFD Multihop Interface Profile.
* `annotation`             - (Read-Only) Annotation of the BFD Multihop Interface Profile object.
* `name_alias`             - (Read-Only) Name Alias of the BFD Multihop Interface Profile object.
* `key`                    - (Read-Only) Authentication Key.
* `key_id`                 - (Read-Only) Authentication Key ID.
* `interface_profile_type` - (Read-Only) Authentication Type.
