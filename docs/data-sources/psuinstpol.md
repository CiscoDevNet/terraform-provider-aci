---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_power_supply_redundancy_policy"
sidebar_current: "docs-aci-data-source-aci_power_supply_redundancy_policy"
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
* `id` - (Read-Only) Attribute id set to the Dn of the Power Supply Redundancy Policy. Type: String
* `annotation` - (Read-Only) Annotation of the Power Supply Redundancy Policy object. Type: String
* `name_alias` - (Read-Only) Name Alias of the Power Supply Redundancy Policy object. Type: String
* `administrative_state` - (Read-Only) Admin Redundancy Mode. The administrative state of the power supply policy. Type: String
