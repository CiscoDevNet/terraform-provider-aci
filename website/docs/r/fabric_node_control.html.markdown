---
layout: "aci"
page_title: "ACI: aci_fabric_node_control"
sidebar_current: "docs-aci-resource-fabric_node_control"
description: |-
  Manages ACI Fabric Node Control
---

# aci_fabric_node_control #
Manages ACI Fabric Node Control

## API Information ##
* `Class` - fabricNodeControl
* `Distinguished Named` - uni/fabric/nodecontrol-{name}

## GUI Information ##
* `Location` - Fabric -> Fabric Policies -> Policies -> Monitoring -> Fabric Node Controls

## Example Usage ##
```hcl
resource "aci_fabric_node_control" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  control = "Dom"
  feature_sel = "telemetry"
  name_alias = "example_name_alias"
  description = "from terraform"
}
```

## Argument Reference ##
* `name` - (Required) Name of object Fabric Node Control.
* `annotation` - (Optional) Annotation of object Fabric Node Control.
* `control` - (Optional) Fabric node control bitmask of object Fabric Node Control. Allowed value is "Dom".
* `feature_sel` - (Optional) Feature Selection of object Fabric Node Control. Allowed values are "analytics", "netflow" and "telemetry". Default value is "telemetry". 
* `description` - (Optional) Description of object Fabric Node Control.
* `name_alias` - (Optional) Name Alias of object Fabric Node Control.


## Importing ##
An existing FabricNodeControl can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_node_control.example <Dn>
```