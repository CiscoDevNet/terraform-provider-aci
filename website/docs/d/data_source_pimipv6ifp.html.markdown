---
subcategory: -
layout: "aci"
page_title: "ACI: aci_pim_ipv6_interface_profile"
sidebar_current: "docs-aci-data-source-pim_ipv6_interface_profile"
description: |-
  Data source for PIM IPv6 ACI Interface Profile
---

# "aci_pim_ipv6_interface_profile" #

Data source for ACI PIM IPv6 Interface Profile

## API Information ##

* `Class` - pimIPV6IfP
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/pimipv6ifp

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage ##

```hcl
data "aci_pim_ipv6_interface_profile" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent Logical Interface Profile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the PIM IPv6 Interface Profile.
* `annotation` - (Optional) Annotation of the PIM IPv6 Interface Profile object. Type: String.
* `name_alias` - (Optional) Name Alias of the PIM IPv6 Interface Profile object. Type: String.
* `name` - (Optional) Name of the PIM IPv6 Interface Profile object. Type: String.
* `relation_pim_ipv6_rs_if_pol` - (Optional) Represents the relation to the PIM IPv6 Interface Policy (class pimIfPol). Type: String.