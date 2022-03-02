---
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

* `Location` - 



## Example Usage ##

```hcl
data "aci_aaa_domain_relationship" "example" {
  parent_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of parent object.
* `name` - (Required) Name of the AAA Security Domain for Parent Object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the AAA Security Domain for Parent Object.
* `annotation` - (Optional) Annotation of object AAA Security Domain for Parent Object.
* `name_alias` - (Optional) Name Alias of object AAA Security Domain for Parent Object.
