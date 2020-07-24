---
layout: "aci"
page_title: "ACI: aci_access_generic"
sidebar_current: "docs-aci-resource-access_generic"
description: |-
  Manages ACI Access Generic
---

# aci_access_generic #
Manages ACI Access Generic

## Example Usage ##

```hcl

resource "aci_access_generic" "example" {
  attachable_access_entity_profile_dn   = "${aci_attachable_access_entity_profile.example.id}"
  name                                  = "example"
  annotation                            = "example"
  name_alias                            = "example"
}

```
## Argument Reference ##
* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.
* `name` - (Required) name of Object access_generic.

* `annotation` - (Optional) annotation for object access_generic.

* `name_alias` - (Optional) name_alias for object access_generic.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Generic.

## Importing ##

An existing Access Generic can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_generic.example <Dn>
```