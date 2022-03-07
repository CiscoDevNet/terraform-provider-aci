---
subcategory: "AAA"
layout: "aci"
subcategory: "AAA"
page_title: "ACI: aci_aaa_domain_relationship"
sidebar_current: "docs-aci-resource-aaa_domain_relationship"
description: |-
  Manages ACI AAA Domain Relationship for Parent Object
---

# aci_aaa_domain_relationship #

Manages ACI AAA Domain Relationship for Parent Object

## API Information ##

* `Class` - aaaDomainRef
* `Distinguished Name` - {parent_dn}/domain-{name}

## GUI Information ##

* `Location` Security Domain list under the Parent Object 


## Example Usage ##

```hcl
resource "aci_aaa_domain_relationship" "example" {
  parent_dn     = aci_tenant.example.id
  aaa_domain_dn = aci_aaa_domain.foosecurity_domain.id
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of parent object.
* `aaa_domain_dn` - (Required) Distinguished name of the AAA Security Domain for Parent Object.
* `annotation` - (Optional) Annotation of the AAA Security Domain for Parent Object.


## Importing ##

An existing AAA Security Domain Relationship object can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_aaa_domain_relationship.example <Dn>
```