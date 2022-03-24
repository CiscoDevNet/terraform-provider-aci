---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_aaep_to_domain"
sidebar_current: "docs-aci-data-aaep-to-domain"
description: |-
  Data source for ACI Attachable Access Entity Profile (AAEP) to domain (VMM, Physical or External domain) relationships.
---

# aci_aaep_to_domain #

Data source for ACI Attachable Access Entity Profile (AAEP) to domain (VMM, Physical or External domain) relationships.


## API Information ##

* `Class` - infraRsDomP
* `Distinguished Name` - uni/infra/attentp-{aaep_name}/rsdomP-[{domain_dn}]

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> AAEP



## Example Usage ##

```hcl
data "aci_aaep_to_domain" "foo_aaep_to_domain" {
  attachable_access_entity_profile_dn = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
  domain_dn                           = aci_l3_domain_profile.fool3_domain_profile.id
}
```

## Argument Reference ##

* `attachable_access_entity_profile_dn` - (Required) Distinguished name of the parent Attachable Access Entity Profile object.
* `domain_dn` - (Required) The Distinguished name of the domain object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Attachable AccessEntity Profile to Domain Relationship object.
* `annotation` - (Optional) Annotation of the Attachable AccessEntity Profile to Domain Relationship object.