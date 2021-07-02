---
layout: "aci"
page_title: "ACI: aci_access_generic"
sidebar_current: "docs-aci-resource-access_generic"
description: |-
  Manages ACI Access Generic
---

# aci_access_generic

Manages ACI Access Generic

## Example Usage

```hcl

resource "aci_access_generic" "example" {
  attachable_access_entity_profile_dn   = aci_attachable_access_entity_profile.example.id
  name                                  = "default"
  annotation                            = "example"
  description                           = "from terraform"
  name_alias                            = "example"
}

```

## Argument Reference

- `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent Attachable Access Entity Profile.
- `name` - (Required) The name of the user defined function object. Name must be "default".
- `annotation` - (Optional) Specifies the annotation of a policy component.
- `description` - (Optional) Specifies the description of a policy component.
- `name_alias` - (Optional) Specifies the alias name of a policy component.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Generic.

## Importing

An existing Access Generic can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_access_generic.example <Dn>
```
