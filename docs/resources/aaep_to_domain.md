---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_aaep_to_domain"
sidebar_current: "docs-aci-resource-aci_aaep-to-domain"
description: |-
  Manages the ACI Attachable Access Entity Profile (AAEP) to domain (VMM, Physical or External domain) relationship.
---

# aci_aaep_to_domain #

Manages the ACI Attachable Access Entity Profile (AAEP) to domain (VMM, Physical or External domain) relationship.

## API Information ##

* `Class` - infraRsDomP
* `Distinguished Name` - uni/infra/attentp-{aaep_name}/rsdomP-[{domain_dn}]

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> AAEP


## Example Usage ##

```hcl
resource "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  domain_dn                           = aci_l3_domain_profile.fool3_domain_profile.id
}
```

## Argument Reference ##

* `attachable_access_entity_profile_dn` - (Required) Distinguished name of the parent Attachable Access Entity Profile object.
* `annotation` - (Optional) Annotation of the Attachable AccessEntity Profile to Domain Relationship object.
* `domain_dn` - (Required) The Distinguished name of the domain object.


## Importing ##

An existing Attachable Access Entity Profile to Domain Relationship object can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_aaep_to_domain.example <Dn>
```