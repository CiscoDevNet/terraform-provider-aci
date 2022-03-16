---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_aaep_to_domain"
sidebar_current: "docs-aci-data-aaep-to-domain"
description: |-
  Data source for ACI Attachable Access Entity Profile - VMM, Physical or External domain interfaces.
---

# aci_aaep_to_domain #

Data source for ACI Attachable Access Entity Profile - VMM, Physical or External domain interfaces.


## API Information ##

* `Class` - infraRsDomP
* `Distinguished Name` - uni/infra/attentp-{name}/rsdomP-[{tDn}]

## GUI Information ##

* `Location` - Attachable Access Entity Profile -> Domains Associated to Interfaces



## Example Usage ##

```hcl
data "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  t_dn                                = aci_l3_domain_profile.fool3_domain_profile.id
}
```

## Argument Reference ##

* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.
* `t_dn` - (Required) TDn of object Domain.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Domain.
* `annotation` - (Optional) Annotation of object Domain.
* `t_dn` - (Optional) Target-dn. The virtual machine manager domain profile.
