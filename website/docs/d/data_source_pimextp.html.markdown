---
subcategory: -
layout: "aci"
page_title: "ACI: aci_pim_external_profile"
sidebar_current: "docs-aci-data-source-pim_external_profile"
description: |-
  Data source for ACI External Profile
---

# aci_pim_external_profile #

Data source for ACI PIM External Profile

## API Information ##

* `Class` - pimExtP
* `Distinguished Name` - uni/tn-{name}/out-{name}/pimextp

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs

## Example Usage ##

```hcl
data "aci_pim_external_profile" "example" {
  l3_outside_dn  = aci_l3_outside.example.id
}
```

## Argument Reference ##

* `l3_outside_dn` - (Required) Distinguished name of the parent L3-Outside object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the PIM External Profile.
* `annotation` - (Optional) Annotation of the PIM External Profile object.
* `name_alias` - (Optional) Name Alias of the PIM External Profile object.
* `enabled_af` - (Optional) Enable Multicast Address Families.
