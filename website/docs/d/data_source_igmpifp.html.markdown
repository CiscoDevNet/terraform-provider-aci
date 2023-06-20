---
subcategory: -
layout: "aci"
page_title: "ACI: aci_igmp_interface_profile"
sidebar_current: "docs-aci-data-source-igmp_interface_profile"
description: |-
  Data source for IGMP ACI Interface Profile
---

# aci_interface_profile #

Data source for ACI IGMP Interface Profile


## API Information ##

* `Class` - igmpIfP
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/igmpIfP

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage ##

```hcl
data "aci_igmp_interface_profile" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent Logical Interface Profile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Interface Profile.
* `annotation` - (Optional) Annotation of the IGMP Interface Profile object. Type: String.
* `name_alias` - (Optional) Name Alias of the IGMP Interface Profile object. Type: String.
* `name` - (Optional) Name of the IGMP Interface Profile object. Type: String.
* `relation_igmp_rs_if_pol` - (Optional) Represents the relation to the IGMP Interface Policy (class igmpIfPol). Type: String.