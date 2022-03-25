---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_regex_community_terms"
sidebar_current: "docs-aci-resource-aci_match_regex_community_terms"
description: |-
  Manages ACI Match Rule Based on Community Regular Expression
---

# aci_match_regex_community_terms #

Manages ACI Match Rule Based on Community Regular Expression

## API Information ##

* `Class` - rtctrlMatchCommRegexTerm
* `Distinguished Name` - uni/tn-{name}/subj-{name}/commrxtrm-{community_type}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match rules -> match_regex_community_terms


## Example Usage ##

```hcl
resource "aci_match_regex_community_terms" "example" {
  match_rule_dn  = aci_match_rule.example.id
  annotation = "orchestrator:terraform"
  community_type = "regular"
  regex = ".*"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of the parent Match Rule object.
* `annotation` - (Optional) Annotation of the Match Rule Based on Community Regular Expression object.
* `community_type` - (Optional) Community Type of the Match Rule Based on Community Regular Expression object. Allowed values are "extended", "regular", and default value is "regular". Type: String.
* `regex` - (Optional) Regular Expression.A regular expression used to specify a pattern to match against the community string.


## Importing ##

An existing Match Rule Based on Community Regular Expression can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_match_regex_community_terms.example <Dn>
```