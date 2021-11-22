---
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
* `Distinguished Named` - uni/tn-{name}/subj-{name}/dest-[{ip}]

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_match_route_destination_rule" "destination" {
  match_rule_dn  = aci_match_rule.rule.id
  ip  = "10.1.1.1"
  aggregate = "no"
  annotation = "orchestrator:terraform"
  from_pfx_len = "0"
  to_pfx_len = "0"
}
```

## Argument Reference ##

* `match_rule_dn` - (Required) Distinguished name of parent Match Rule object.
* `ip` - (Required) Ip of object Match Route Destination Rule.
* `annotation` - (Optional) Annotation of object Match Route Destination Rule.
* `aggregate` - (Optional) Aggregated Route. Allowed values are "false", "true" and default value is "false". Type: String.

* `from_pfx_len` - (Optional) Start of Prefix Length. Allowed range is 0-128 and default value is "0".
* `ip` - (Optional) Match IP Address.
* `to_pfx_len` - (Optional) End of Prefix Length. Allowed range is 0-128 and default value is "0".


## Importing ##

An existing Match Route Destination Rule can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_match_route_destination_rule.example <Dn>
```