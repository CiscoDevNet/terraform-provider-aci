---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_aaep_to_domain"
sidebar_current: "docs-aci-resource-aaep-to-domain"
description: |-
  Manages ACI Attachable Access Entity Profile - VMM, Physical or External domain interfaces.
---

# aci_aaep_to_domain #

Manages ACI Attachable Access Entity Profile  - VMM, Physical or External domain interfaces.

## API Information ##

* `Class` - infraRsDomP
* `Distinguished Name` - uni/infra/attentp-{name}/rsdomP-[{tDn}]

## GUI Information ##

* `Location` - Attachable Access Entity Profile -> Domains Associated to Interfaces


## Example Usage ##

```hcl
resource "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  t_dn                                = aci_l3_domain_profile.fool3_domain_profile.id
}
```

## Argument Reference ##

* `attachable_access_entity_profile_dn` - (Required) Distinguished name of the parent AttachableAccessEntityProfile object.
* `annotation` - (Optional) Annotation of the object Domain.

* `t_dn` - (Optional) Target-dn.The virtual machine manager domain profile.


## Importing ##

An existing Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_aaep_to_domain.example <Dn>
```