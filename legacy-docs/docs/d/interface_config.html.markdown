---
subcategory: Fabric Access Policies
layout: "aci"
page_title: "ACI: aci_interface_config"
sidebar_current: "docs-aci-data-source-interface_config"
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
data "aci_interface_config" "access_port_config_1001" {
  node         = 1001
  interface    = "1/1"
}

data "aci_interface_config" "access_port_config_1003" {
  node         = 1003
  interface    = "2/2/2"
  port_type    = "fabric"
}
```

## Argument Reference ##
* `node` - (Required) The Node ID of the Port Configuration object. Type: Integer.
* `interface` - (Required) The Interface address of the Port Configuration object. The format of the interface value should be 1/1/1 (card/port_id/sub_port) or 1/1 (card/port_id). Type: String.
* `port_type` - (Optional) The Type of the Port Configuration object. Allowed values are "access", "fabric". Default value is "access". Type: String.

## Attribute Reference ##
* `id` - The Attribute ID set to the Dn of the Port Configuration.
* `role` - (Read-Only) The Role of the Port Configuration object. Type: String.
* `policy_group` - (Read-Only) The Distinguished Name of the Policy Group being associated with the Port Configuration object. Type: String.
* `breakout` - (Read-Only) The Breakout Map of the Port Configuration object. Type: String.
* `admin_state` - (Read-Only) The Admin State of the Port Configuration object. Type: String.
* `pc_member` - (Read-Only) The Distinguished Name of the Port Channel Member being associated with the Port Configuration object. Type: String.
* `description` - (Read-Only) The Description of the Port Configuration object. Type: String.
* `annotation` - (Read-Only) The Annotation of the Port Configuration object. Type: String.
* `name_alias` - (Read-Only) The Name Alias of the Port Configuration object. Type: String.
* `operational_associated_group` - (Read-Only) The Distinguished Name of the Operational Associated Group.
* `operational_associated_sub_group` - (Read-Only) The Distinguished Name of the Operational Associated Sub Group. This is only applicable when `port_type` is "access".
* `port_dn` - (Read-Only) The Distinguished Name of the Port. This DN can be used as a reference to the static path resources.
* `pc_port_dn` - (Read-Only) The Distinguished Name of the Port Channel Port. This contains the VPC Path DN when two switches in a VPC pair have multiple ports that are in the same Policy Group. This DN can be used as a reference to the static path resources and is only applicable when `port_type` is "access".

