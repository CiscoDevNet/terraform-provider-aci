---
layout: "aci"
page_title: "ACI: aci_tag"
sidebar_current: "docs-aci-data-source-tag"
description: |-
  Data source for ACI Tag
---

# aci_tag #

Data source for ACI Tag

## API Information ##

* `Class` - tagTag
* `Distinguished Named` - {parent_dn}/tagKey-{key}

## GUI Information ##

* `Location` - Under every object as Tag in recent APIC versions.

## Example Usage ##

```hcl
data "aci_tag" "example" {
  parent_dn  = aci_tenant.example.id
  key  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the object to which the Tag is attached to.
* `key` - (Required) Key of the Tag.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Tag.
* `value` - (Optional) The value of the Tag.
