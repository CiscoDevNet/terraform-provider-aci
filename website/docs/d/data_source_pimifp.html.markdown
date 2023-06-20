---
subcategory: -
layout: "aci"
page_title: "ACI: aci_pim_interface_profile"
sidebar_current: "docs-aci-data-source-pim_interface_profile"
description: |-
  Data source for PIM ACI Interface Profile
---

# aci_pim_interface_profile #

Data source for ACI PIM Interface Profile


## API Information ##

* `Class` - pimIfP
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/pimifp

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage ##

```hcl
data "aci_pim_interface_profile" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent Logical Interface Profile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the PIM Interface Profile.
* `annotation` - (Optional) Annotation of the PIM Interface Profile object. Type: String.
* `name_alias` - (Optional) Name Alias of the PIM Interface Profile object. Type: String.
* `name` - (Optional) Name of the PIM Interface Profile object. Type: String.
* `relation_pim_rs_if_pol` - (Optional) Represents the relation to the PIM Interface Policy (class pimIfPol). Type: String.