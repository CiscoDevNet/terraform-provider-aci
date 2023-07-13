---
subcategory: -
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
  name  = "example"
  admin_rdn_m = "comb"
  annotation = "example"
  name_alias = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Power Supply Redundancy Policy object.
* `annotation` - (Optional) Annotation of the Power Supply Redundancy Policy object.
* `name_alias` - (Optional) Name Alias of the Power Supply Redundancy Policy object.
* `admin_rdn_m` - (Optional) Admin Redundancy Mode. The administrative state of the power supply policy. Allowed values are "comb", "insrc-rdn", "n-rdn", "not-supp", "ps-rdn", "rdn", "sinin-rdn", "unknown", and default value is "comb". Type: String.

## Importing ##

An existing PowerSupplyRedundancyPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_power_supply_redundancy_policy.example <Dn>
```