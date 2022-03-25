---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_regex_community_terms"
sidebar_current: "docs-aci-data-source-aci_match_regex_community_terms"
description: |-
  Data source for ACI Match Rule Based on Community Regular Expression
---

# aci_match_regex_community_terms #

Data source for ACI Match Rule Based on Community Regular Expression


## API Information ##

* `Class` - rtctrlMatchCommRegexTerm
* `Distinguished Name` - uni/tn-{name}/subj-{name}/commrxtrm-{community_type}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match rules -> match_regex_community_terms



## Example Usage ##

```hcl
data "aci_match_regex_community_terms" "example" {
  match_rule_dn  = aci_match_rule.example.id
  community_type  = "regular"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of parent Match Rule object.
* `community_type` - (Required) Community Type of the Match Rule Based on Community Regular Expression object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Match Rule Based on Community Regular Expression.
* `annotation` - (Optional) Annotation of the Match Rule Based on Community Regular Expression object.
* `community_type` - (Optional) Community Type of the Match Rule Based on Community Regular Expression object.
* `regex` - (Optional) Regular Expression. A regular expression used to specify a pattern to match against the community string.
