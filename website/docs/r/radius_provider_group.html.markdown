---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_radius_provider_group"
sidebar_current: "docs-aci-resource-radius_provider_group"
description: |-
  Manages ACI RADIUS Provider Group
---

# aci_radius_provider_group #

Manages ACI RADIUS Provider Group

## API Information ##

* `Class` - aaaRadiusProviderGroup
* `Distinguished Named` - uni/userext/radiusext/radiusprovidergroup-{name}


## Example Usage ##

```hcl
resource "aci_radius_provider_group" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  name_alias  = "radius_provider_group_alias"
  description = "From Terraform"
}
```

## Argument Reference ##


* `name` - (Required) Name of object RADIUS Provider Group. Type: String.
* `annotation` - (Optional) Annotation of object RADIUS Provider Group. Type: String.
* `name_alias` - (Optional) Name Alias of object RADIUS Provider Group. Type: String.
* `description` - (Optional) Description of object RADIUS Provider Group. Type: String.

## Importing ##

An existing RADIUSProviderGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_radius_provider_group.example <Dn>
```