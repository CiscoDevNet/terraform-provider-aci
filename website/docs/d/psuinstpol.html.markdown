---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_power_supply_redundancy_policy"
sidebar_current: "docs-aci-data-source-power_supply_redundancy_policy"
description: |-
  Data source for ACI Power Supply Redundancy Policy
---

# aci_power_supply_redundancy_policy #

Data source for ACI Power Supply Redundancy Policy

## API Information ##

* `Class` - psuInstPol
* `Distinguished Name` - uni/fabric/psuInstP-{name}

## GUI Information ##

* `Location` - Fabric -> Fabric Policies -> Policies -> Switch -> Power Supply Redundancy

## Example Usage ##

```hcl
data "aci_power_supply_redundancy_policy" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Power Supply Redundancy Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Power Supply Redundancy Policy.
* `annotation` - (Read-Only) Annotation of the Power Supply Redundancy Policy object.
* `name_alias` - (Read-Only) Name Alias of the Power Supply Redundancy Policy object.
* `admin_rdn_m` - (Read-Only) Admin Redundancy Mode. The administrative state of the power supply policy.
