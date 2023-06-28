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
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/bfdMhIfP

## GUI Information ##

* `Location` -  Tenant -> Networking -> L3Out -> Logical Node Profiles -> Logical Interface Profiles 



## Example Usage ##

```hcl
data "aci_bfd_multihop_interface_profile" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of parent LogicalInterfaceProfile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the BFD Multihop Interface Profile.
* `annotation` - (Optional) Annotation of the BFD Multihop Interface Profile object.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Interface Profile object.
* `key` - (Optional) Authentication Key. Authentication key
* `key_id` - (Optional) Authentication Key ID. Authentication key id
* `interface_profile_type` - (Optional) Authentication Type. Authentication type
