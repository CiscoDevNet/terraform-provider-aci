---
subcategory: "AAA"
layout: "aci"
page_title: "ACI: aci_aaa_domain_relationship"
sidebar_current: "docs-aci-data-source-aaa_domain_relationship"
description: |-
  Data source for ACI AAA Domain Relationship for Parent Object
---

# aci_aaa_domain_relationship #

Data source for ACI AAA Domain Relationship for Parent Object


## API Information ##

* `Class` - aaaDomainRef
* `Distinguished Name` - {parent_dn}/domain-{name}

## GUI Information ##

* `Location` - Security Domain list under the Parent Object 



## Example Usage ##

```hcl
data "aci_aaa_domain_relationship" "example" {
  parent_dn     = aci_tenant.example.id
  aaa_domain_dn = aci_aaa_domain.foosecurity_domain.id
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of parent object.
* `aaa_domain_dn` - (Required) Distinguished name of the AAA Security Domain for Parent Object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the AAA Security Domain for Parent Object.
* `annotation` - (Optional) Annotation of object AAA Security Domain for Parent Object.
* `name_alias` - (Optional) Name Alias of object AAA Security Domain for Parent Object.
