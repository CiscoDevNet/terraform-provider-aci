---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_radius_provider_group"
sidebar_current: "docs-aci-data-source-aci_radius_provider_group"
description: |-
  Data source for ACI RADIUS Provider Group
---

# aci_radius_provider_group #

Data source for ACI RADIUS Provider Group


## API Information ##

* `Class` - aaaRadiusProviderGroup
* `Distinguished Name` - uni/userext/radiusext/radiusprovidergroup-{name}


## Example Usage ##

```hcl
data "aci_radius_provider_group" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of object RADIUS Provider Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the RADIUS Provider Group.
* `annotation` - (Optional) Annotation of object RADIUS Provider Group.
* `name_alias` - (Optional) Name Alias of object RADIUS Provider Group.
* `description` - (Optional) Description of object RADIUS Provider Group.