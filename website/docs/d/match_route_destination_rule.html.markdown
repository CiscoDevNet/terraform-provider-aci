---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_route_destination_rule"
sidebar_current: "docs-aci-data-source-match_route_destination_rule"
description: |-
  Data source for ACI Match Route Destination Rule
---

# aci_match_route_destination_rule #

Data source for ACI Match Route Destination Rule


## API Information ##

* `Class` - rtctrlMatchRtDest
* `Distinguished Named` - uni/tn-{name}/subj-{name}/dest-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match Rules -> Match Prefix



## Example Usage ##

```hcl
data "aci_match_route_destination_rule" "example" {
  match_rule_dn  = aci_match_rule.example.id
  ip  = "example"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of parent MatchRule object.
* `ip` - (Required) ip of object Match Route Destination Rule.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Match Route Destination Rule.
* `annotation` - (Optional) Annotation of object Match Route Destination Rule.
* `name_alias` - (Optional) Name Alias of object Match Route Destination Rule.
* `aggregate` - (Optional) Aggregated Route. Aggregated Route
* `greater_than_mask` - (Optional) Start of Prefix Length. Prefix list range
* `ip` - (Optional) Match IP Address. null
* `less_than_mask` - (Optional) End of Prefix Length. 
