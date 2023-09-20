---
layout: "aci"
page_title: "ACI: aci_tag"
sidebar_current: "docs-aci-resource-tag"
description: |-
  Manages ACI Tag
---

# aci_tag #

Manages ACI Tag

## API Information ##

* `Class` - tagTag
* `Distinguished Name` - {parent_dn}/tagKey-{key}

## GUI Information ##

* `Location` - Under every object as Policy Tags in the Operational tab in recent APIC versions.

## Example Usage ##

```hcl
resource "aci_tag" "example" {
  parent_dn  = aci_tenant.example.id
  key  = "example"
  value = "example-value"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the object to which the Tag is attached to.
* `key` - (Required) The key of the Tag.
* `value` - (Required) The value of the Tag.


## Importing ##

An existing Tag can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tag.example <Dn>
```