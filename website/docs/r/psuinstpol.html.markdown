---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_power_supply_redundancy_policy"
sidebar_current: "docs-aci-resource-power_supply_redundancy_policy"
description: |-
  Manages ACI Power Supply Redundancy Policy
---

# aci_power_supply_redundancy_policy #

Manages ACI Power Supply Redundancy Policy

## API Information ##

* `Class` - psuInstPol
* `Distinguished Name` - uni/fabric/psuInstP-{name}

## GUI Information ##

* `Location` - Fabric -> Fabric Policies -> Policies -> Switch -> Power Supply Redundancy

## Example Usage ##

```hcl
resource "aci_power_supply_redundancy_policy" "example" {
  name                  = "example"
  administrative_state  = "comb"
  annotation            = "example"
  name_alias            = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Power Supply Redundancy Policy object. Type: String
* `annotation` - (Optional) Annotation of the Power Supply Redundancy Policy object. Type: String
* `name_alias` - (Optional) Name Alias of the Power Supply Redundancy Policy object. Type: String
* `administrative_state` - (Optional) The administrative state of the power supply policy. Allowed values are "comb" (Combined), "insrc-rdn" (Input source redundancy), "n-rdn" (Non redundant), "not-supp" (Not supported), "ps-rdn" (N+1 Redundancy), "rdn" (N+N Redundancy), "sinin-rdn" (Single input redundancy), "unknown", and default value is "comb". Type: String.

## Importing ##

An existing PowerSupplyRedundancyPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_power_supply_redundancy_policy.example <Dn>
```