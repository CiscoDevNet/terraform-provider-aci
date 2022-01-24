---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_rule"
sidebar_current: "docs-aci-resource-match_rule"
description: |-
  Manages ACI Match Rule
---

# aci_match_rule #

Manages ACI Match Rule

## API Information ##

* `Class` - rtctrlSubjP
* `Distinguished Name` - uni/tn-{name}/subj-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match Rules

## Example Usage ##

```hcl
resource "aci_match_rule" "rule" {
  tenant_dn  = aci_tenant.terraform_tenant.id
  name  = "match_rule"
  annotation = "orchestrator:terraform"

}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of object Match Rule.
* `annotation` - (Optional) Annotation of object Match Rule.
* `name_alias` - (Optional) Name Alias of object Match Rule.
* `description` - (Optional) Description of object Match Rule.



## Importing ##

An existing Match Rule can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_match_rule.example <Dn>
```