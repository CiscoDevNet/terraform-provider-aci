---
layout: "aci"
page_title: "ACI: aci_access_generic"
sidebar_current: "docs-aci-data-source-access_generic"
description: |-
  Data source for ACI Access Generic
---

# aci_access_generic #
Data source for ACI Access Generic

## Example Usage ##

```hcl

data "aci_access_generic" "example" {
  attachable_access_entity_profile_dn   = "${aci_attachable_access_entity_profile.example.id}"
  name                                  = "example"
}

```

## Argument Reference ##

* `attachable_access_entity_profile_dn` - (Required) Distinguished name of parent AttachableAccessEntityProfile object.
* `name` - (Required) name of Object access_generic.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Generic.
* `annotation` - (Optional) annotation for object access_generic.
* `name_alias` - (Optional) name_alias for object access_generic.
