---
subcategory: -
layout: "aci"
page_title: "ACI: aci_igmp_interface_profile"
sidebar_current: "docs-aci-resource-igmp_interface_profile"
description: |-
  Manages ACI IGMP Interface Profile
---

# aci_igmp_interface_profile #

Manages ACI IGMP Interface Profile

## API Information ##

* `Class` - igmpIfP
* `Distinguished Name` - uni/tn-{name}/out-{name}/lnodep-{name}/lifp-{name}/igmpIfP

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage ##

```hcl
resource "aci_igmp_interface_profile" "igmp_interface" {
  logical_interface_profile_dn  = aci_logical_interface_profile.foological_interface_profile.id
  name = "igmp1"
}
```

## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent Logical Interface Profile object.
* `annotation` - (Optional) Annotation of the IGMP Interface Profile object. Type: String.
* `name_alias` - (Optional) Name Alias of the IGMP Interface Profile object. Type: String.
* `name` - (Optional) Name of the IGMP Interface Profile object. Type: String.
* `relation_igmp_rs_if_pol` - (Optional) Represents the relation to the IGMP Interface Policy (class igmpIfPol). Type: String.

## Importing ##

An existing IGMP Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_igmp_interface_profile.example <Dn>
```