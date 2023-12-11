---
subcategory: "Fabric Policies"
layout: "aci"
page_title: "ACI: aci_fabric_node_control"
sidebar_current: "docs-aci-data-source-aci_fabric_node_control"
description: |-
  Data source for ACI Fabric Node Control
---

# aci_fabric_node_control #
Data source for ACI Fabric Node Control

## API Information ##
* `Class` - fabricNodeControl
* `Distinguished Name` - uni/fabric/nodecontrol-{name}

## GUI Information ##
* `Location` - Fabric -> Fabric Policies -> Policies -> Monitoring -> Fabric Node Controls

## Example Usage ##
```hcl
data "aci_fabric_node_control" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object Fabric Node Control.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Fabric Node Control.
* `annotation` - (Optional) Annotation of object Fabric Node Control.
* `name_alias` - (Optional) Name Alias of object Fabric Node Control.
* `control` - (Optional) Fabric node control bitmask of object Fabric Node Control. 
* `feature_sel` - (Optional) Feature Selection of object Fabric Node Control.
* `description` - (Optional) Description of object Fabric Node Control.
