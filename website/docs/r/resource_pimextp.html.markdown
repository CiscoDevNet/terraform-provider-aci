---
subcategory: -
layout: "aci"
page_title: "ACI: aci_pim_external_profile"
sidebar_current: "docs-aci-resource-pim_external_profile"
description: |-
  Manages ACI PIM External Profile
---

# aci_external_profile #

Manages ACI PIM External Profile

## API Information ##

* `Class` - pimExtP
* `Distinguished Name` - uni/tn-{name}/out-{name}/pimextp

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs

## Example Usage ##

```hcl
resource "aci_pim_external_profile" "example" {
  l3_outside_dn  = aci_l3outside.example.id
  enabled_af = ["ipv4-mcast"]
}
```

## Argument Reference ##

* `l3outside_dn` - (Required) Distinguished name of the parent L3-Outside object.
* `annotation` - (Optional) Annotation of the PIM External Profile object.
* `name_alias` - (Optional) Name Alias of the PIM External Profile object.
* `enabled_af` - (Optional) Enable Multicast Address Families. Allowed values are "ipv4-mcast", "ipv6-mcast", and default value is "ipv4-mcast". Type: List.

## Importing ##

An existing PIM External Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_pim_external_profile.example <Dn>
```