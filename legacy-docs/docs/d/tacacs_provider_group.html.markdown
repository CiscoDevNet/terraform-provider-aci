---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_provider_group"
sidebar_current: "docs-aci-data-source-tacacs_provider_group"
description: |-
  Data source for ACI TACACS + Provider Group
---

# aci_tacacs_provider_group #
Data source for ACI TACACS+ Provider Group


## API Information ##
* `Class` - aaaTacacsPlusProviderGroup
* `Distinguished Name` - uni/userext/tacacsext/tacacsplusprovidergroup-{name}

## Example Usage ##
```hcl
data "aci_tacacs_provider_group" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) Name of object TACACS + Provider Group.er Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the TACACS+ Provider Group.
* `annotation` - (Optional) Annotation of object TACACS + Provider Group.
* `description` - (Optional) Description of object TACACS + Provider Group.
* `name_alias` - (Optional) Name alias of object TACACS + Provider Group.
