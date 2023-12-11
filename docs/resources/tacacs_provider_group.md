---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_tacacs_provider_group"
sidebar_current: "docs-aci-resource-aci_tacacs_provider_group"
description: |-
  Manages ACI TACACS + Provider Group
---

# aci_tacacs_provider_group #
Manages ACI TACACS + Provider Group

## API Information ##
* `Class` - aaaTacacsPlusProviderGroup
* `Distinguished Name` - uni/userext/tacacsext/tacacsplusprovidergroup-{name}


## Example Usage ##

```hcl
resource "aci_tacacs_provider_group" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
}
```

## Argument Reference ##
* `name` - (Required) Name of object TACACS + Provider Group.
* `annotation` - (Optional) Annotation of object TACACS + Provider Group.
* `description` - (Optional) Description of object TACACS + Provider Group.
* `name_alias` - (Optional) Name alias of object TACACS + Provider Group.

## Importing ##
An existing TACACS + ProviderGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tacacs_provider_group.example <Dn>
```