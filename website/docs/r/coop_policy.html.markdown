---
layout: "aci"
page_title: "ACI: aci_coop_policy"
sidebar_current: "docs-aci-resource-coop_policy"
description: |-
  Manages ACI COOP Policy
---

# aci_coop_policy #

Manages ACI COOP Policy

## API Information ##

* `Class` - coopPol
* `Distinguished Named` - uni/fabric/pol-{name}

## GUI Information ##

* `Location` - System -> System Settings -> COOP Group -> Policy


## Example Usage ##

```hcl
resource "aci_coop_policy" "example" {

  annotation  = "orchestrator:terraform"
  type        = "compatible"
  name_alias  = "alias_coop_policy"
  description = "From Terraform"

}
```

## NOTE ##
User can use resource of type aci_coop_policy to change configuration of object COOP Group Policy. User cannot create more than one instances of object COOP Group Policy.

## Argument Reference ##

* `annotation` - (Optional) Annotation of object COOP Group Policy.
* `type` - (Optional) Authentication type.The specific type of the object or component. Allowed values are "compatible", "strict". Type: String.
* `name_alias` - (Optional) Name Alias of object COOP Group Policy. Type: String.
* `description` - (Optional) Description of object COOP Group Policy. Type: String.


## Importing ##

An existing COOPGroupPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_coop_policy.example <Dn>
```