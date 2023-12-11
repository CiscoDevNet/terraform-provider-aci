---
subcategory: Fabric Access Policies
layout: "aci"
page_title: "ACI: aci_interface_config"
sidebar_current: "docs-aci-resource-aci_interface_config"
description: |-
  Manages ACI Access and Fabric Ports is only supported for ACI 5.2(7)+
---

# aci_interface_config #

Manages ACI Access and Fabric Ports is only supported for ACI 5.2(7)+

## API Information ##

* `Class` - infraPortConfig
* `Distinguished Name` - uni/infra/portconfnode-{node}-card-{card}-port-{port}-sub-{subPort}
* `Class` - fabricPortConfig
* `Distinguished Name` - uni/fabric/portconfnode-{node}-card-{card}-port-{port}-sub-{subPort}


## GUI Information ##

* `Location` - Fabric -> Access Policies -> Interface Configuration


## Example Usage ##

```hcl
# Access Port Configuration
resource "aci_interface_config" "access_port_config_1001" {
  node         = 1001
  interface    = "1/1"
  policy_group = aci_leaf_access_port_policy_group.leaf_access_port.id # Policy Group and Breakout cannot be configured togater
}

# Breakout an Access Port Configuration
resource "aci_interface_config" "access_port_config_1001_brkout" {
  node      = 1001
  interface = "1/1"
  breakout  = "100g-4x" # Policy Group and Breakout cannot be configured togater
}

# Fabric Port Configuration
resource "aci_interface_config" "fabric_port_config" {
  node      = 1003
  interface = "2/2/2"
  port_type = "fabric"
}
```

## Argument Reference ##

* `node` - (Required) The Node ID of the Port Configuration object. Type: Integer.
* `interface` - (Required) The Interface address of the Port Configuration object. The format of the interface value should be 1/1/1 (card/port_id/sub_port) or 1/1 (card/port_id). Type: String.
* `port_type` - (Optional) The Type of the Port Configuration object. Allowed values are "access", "fabric". Default value is "access". Type: String.
* `role` - (Optional) The Role of the Port Configuration object. Allowed values are "leaf", "spine". Default value is "leaf". Type: String.
* `policy_group` - (Optional) The Distinguished Name of the Policy Group being associated with the Port Configuration object. The Policy Group and Breakout cannot be configured simultaneously. Type: String.
* `breakout` - (Optional) The Breakout Map of the Port Configuration object. Allowed values are "100g-2x", "100g-4x", "10g-4x", "25g-4x", "50g-8x", "none". The Policy Group and Breakout cannot be configured simultaneously. Type: String.
* `admin_state` - (Optional) The Admin State of the Port Configuration object. Allowed values are "up", "down". Default value is "up". Type: String.
* `pc_member` - (Optional) The Distinguished Name of the Port Channel Member being associated with the Port Configuration object. Type: String.
* `description` - (Optional) The Description of the Port Configuration object. Type: String.
* `annotation` - (Optional) The Annotation of the Port Configuration object. Type: String.
* `name_alias` - (Optional) The Name Alias of the Port Configuration object. Type: String.
* `operational_associated_group` - (Read-Only) The Distinguished Name of the Operational Associated Group.
* `operational_associated_sub_group` - (Read-Only) The Distinguished Name of the Operational Associated Sub Group. This is only applicable when `port_type` is "access".
* `port_dn` - (Read-Only) The Distinguished Name of the Port. This DN can be used as a reference to the static path resources.
* `pc_port_dn` - (Read-Only) The Distinguished Name of the Port Channel Port. This contains the VPC Path DN when two switches in a VPC pair have multiple ports that are in the same Policy Group. This DN can be used as a reference to the static path resources and is only applicable when `port_type` is "access".

## Importing ##

An existing Port Configuration can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_interface_config.example <Dn>
```