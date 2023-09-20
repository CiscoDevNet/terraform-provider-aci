---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_rule"
sidebar_current: "docs-aci-data-source-match_rule"
description: |-
  Data source for ACI Match Rule
---

# aci_match_rule #

Data source for ACI Match Rule


## API Information ##

* `Class` - rtctrlSubjP
* `Distinguished Name` - uni/tn-{name}/subj-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match Rules



## Example Usage ##

```hcl
data "aci_match_rule" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of object Match Rule.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Match Rule.
* `annotation` - (Optional) Annotation of object Match Rule.
* `name_alias` - (Optional) Name Alias of object Match Rule.
