---
layout: "aci"
page_title: "ACI: aci_annotation"
sidebar_current: "docs-aci-resource-annotation"
description: |-
  Manages ACI Annotation
---

# aci_annotation #

Manages ACI Annotation

## API Information ##

* `Class` - tagAnnotation
* `Distinguished Name` - {parent_dn}/annotationKey-[{key}]

## GUI Information ##

* `Location` - Under every object as Annotations in recent APIC versions.


## Example Usage ##

```hcl
resource "aci_annotation" "example" {
  parent_dn  = aci_tenant.example.id
  key  = "example"
  value = "example_value"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the object to which the Annotation is attached to.
* `key` - (Required) The key of the Annotation.
* `value` - (Required) The value of the Annotation.


## Importing ##

An existing Annotation can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_annotation.example <Dn>
```