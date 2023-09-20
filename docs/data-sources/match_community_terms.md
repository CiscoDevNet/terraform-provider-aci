---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_community_terms"
sidebar_current: "docs-aci-data-source-match_community_terms"
description: |-
  Data source for ACI Match Community Term
---

# aci_match_community_terms #

Data source for ACI Match Community Term


## API Information ##

* `Class` - rtctrlMatchCommTerm
* `Distinguished Name` - uni/tn-{name}/subj-{name}/commtrm-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match rules -> Match Community Terms



## Example Usage ##

```hcl
data "aci_match_community_terms" "example" {
  match_rule_dn  = aci_match_rule.example.id
  name  = "example"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of the parent Match Rule object.
* `name` - (Required) Name of the Match Community Term object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Match Community Term.
* `annotation` - (Optional) Annotation of the Match Community Term object.
* `name_alias` - (Optional) Name Alias of the Match Community Term object.
* `match_community_factors` - (Optional) Match Community Factor object.Type: Block.
  * `community` - (Required) The community of the Match Community Factor object. Type: String.
  * `scope` - (Optional) The scope of the Match Community Factor object. Allowed values are "transitive", "non-transitive", and default value is "transitive". Type: String.
  * `community` - (Optional) The community of the Match Community Factor object. Type: String.
  * `description` - (Optional) The description of the Match Community Factor object.

