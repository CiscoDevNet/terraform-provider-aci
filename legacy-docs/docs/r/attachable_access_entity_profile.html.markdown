---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_attachable_access_entity_profile"
sidebar_current: "docs-aci-resource-aci_attachable_access_entity_profile"
description: |-
  Manages ACI Attachable Access Entity Profile
---

# aci_attachable_access_entity_profile #

Manages ACI Attachable Access Entity Profile

## Example Usage ##

```hcl
resource "aci_attachable_access_entity_profile" "example" {
  description = "AAEP description"
  name        = "demo_entity_prof"
  annotation  = "tag_entity"
  name_alias  = "alias_entity"
}
```

## Argument Reference ##

* `name` - (Required) Name of Object attachable access entity profile.
* `annotation` - (Optional) Annotation for object attachable access entity profile.
* `name_alias` - (Optional) Name alias for object attachable access entity profile.
* `description` - (Optional) Description for object attachable access entity profile.
* `relation_infra_rs_dom_p` - (Optional) Relation to class infraADomP. Cardinality - N_TO_M. Type - [Set of String].

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Attachable Access Entity Profile.

## Importing ##

An existing Attachable Access Entity Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_attachable_access_entity_profile.example <Dn>
```
