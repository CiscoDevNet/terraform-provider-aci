---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_community_terms"
sidebar_current: "docs-aci-resource-match_community_terms"
description: |-
  Manages ACI Match Community Term
---

# aci_match_community_terms #

Manages ACI Match Community Term

## API Information ##

* `Class` - rtctrlMatchCommTerm
* `Distinguished Name` - uni/tn-{name}/subj-{name}/commtrm-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match rules -> Match Community Terms


## Example Usage ##

```hcl
resource "aci_match_community_terms" "example" {
  match_rule_dn  = aci_match_rule.example.id
  name  = "example"
  annotation = "orchestrator:terraform"

}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of the parent Match Rule object.
* `name` - (Required) Name of the Match Community Term object.
* `annotation` - (Optional) Annotation of the Match Community Term object.



## Importing ##

An existing Match Community Term can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_match_community_terms.example <Dn>
```