---
layout: "aci"
page_title: "ACI: aci_annotation"
sidebar_current: "docs-aci-data-source-annotation"
description: |-
  Data source for ACI Annotation
---

# aci_annotation #

Data source for ACI Annotation

## API Information ##

* `Class` - tagAnnotation
* `Distinguished Name` - {parent_dn}/annotationKey-[{key}]

## GUI Information ##

* `Location` - Under every object as Annotations in recent APIC versions. 

## Example Usage ##

```hcl
data "aci_annotation" "example" {
  parent_dn  = aci_tenant.example.id
  key  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the object to which the Annotation is attached to.
* `key` - (Required) Key of the Annotation.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Annotation.
* `value` - (Optional) The value of the Annotation.
