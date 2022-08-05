---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_match_route_destination_rule"
sidebar_current: "docs-aci-resource-match_route_destination_rule"
description: |-
  Manages ACI Match Route Destination Rule
---

# aci_match_route_destination_rule #

Manages ACI Match Route Destination Rule

## API Information ##

* `Class` - rtctrlMatchRtDest
* `Distinguished Name` - uni/tn-{name}/subj-{name}/dest-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Match Rules -> Match Prefix

## Example Usage ##

```hcl
resource "aci_match_route_destination_rule" "destination" {
  match_rule_dn  = aci_match_rule.rule.id
  ip  = "10.1.1.1"
  aggregate = "no"
  annotation = "orchestrator:terraform"
  greater_than_mask = "0"
  less_than_mask = "0"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of parent Match Rule object.
* `ip` - (Required) Ip of object Match Route Destination Rule.
* `annotation` - (Optional) Annotation of object Match Route Destination Rule.
* `aggregate` - (Optional) Aggregated Route. Allowed values are "yes", "no" and default value is "no". Type: String.
* `greater_than_mask` - (Optional) Start of Prefix Length. Allowed range is 0-128 and default value is "0".
* `less_than_mask` - (Optional) End of Prefix Length. Allowed range is 0-128 and default value is "0".
* `description` - (Optional) Description of object Match Route Destination Rule.
* `name_alias` - (Optional) Name alias of object Match Route Destination Rule.


## Importing ##

An existing Match Route Destination Rule can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_match_route_destination_rule.example <Dn>
```