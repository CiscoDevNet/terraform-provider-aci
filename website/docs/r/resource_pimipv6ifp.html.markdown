---
subcategory: -
layout: "aci"
page_title: "ACI: aci_pim_ipv6_interface_profile"
sidebar_current: "docs-aci-resource-pim_ipv6_interface_profile"
description: |-
  Manages ACI PIM IPv6 Interface Profile
---

# aci_pim_ipv6_interface_profile #

Manages ACI PIM IPv6 Interface Profile

## API Information ##

* `Class` - pimIPV6IfP
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/pimipv6ifp

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage ##

```hcl
resource "aci_pim_ipv6_interface_profile" "pimipv6_interface" {
  logical_interface_profile_dn  = aci_logical_interface_profile.foological_interface_profile.id
  name = "pimipv61"
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent Logical Interface Profile object.
* `annotation` - (Optional) Annotation of the PIM IPv6 Interface Profile object. Type: String.
* `name_alias` - (Optional) Name Alias of the PIM IPv6 Interface Profile object. Type: String.
* `name` - (Optional) Name of the PIM IPv6 Interface Profile object. Type: String.
* `relation_pim_ipv6_rs_if_pol` - (Optional) Represents the relation to the PIM IPv6 Interface Policy (class pimIfPol). Type: String.

## Importing ##

An existing PIM IPv6 Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_pim_ipv6_interface_profile.example <Dn>
```