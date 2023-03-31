---
subcategory: Fabric Access Policies
layout: "aci"
page_title: "ACI: aci_interface_config"
sidebar_current: "docs-aci-data-source-interface_config"
description: |-
  Manages ACI Access and Fabric Ports is only supported for ACI 5.2(5)+
---

# aci_interface_config #

Manages ACI Access and Fabric Ports is only supported for ACI 5.2(5)+

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
  interface    = "1/1/1"
  port_type    = "access"
}
```

## Argument Reference ##
* `node` - (Required) Node ID of the Port Configuration object. Type: Integer.
* `interface` - (Required) Interface address of the Port Configuration object. The format of the interface value should be 1/1/1 (card/port_id/sub_port) or 1/1 (card/port_id). Type: String.
* `port_type` - (Required) Type of the Port Configuration object. Allowed values are "access", "fabric". Type: String.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Port Configuration.
* `role` - (Optional) Role of the Port Configuration object. Type: String.
* `policy_group` - (Optional) The distinguished name of the Policy Group being associated with the Port Configuration object. Type: String.
* `breakout` - (Optional) Breakout Map of the Port Configuration object. Type: String.
* `admin_state` - (Optional) Admin State of the Port Configuration object. Type: String.
* `pc_member` - (Optional) The distinguished name of the Port Channel Member being associated with the Port Configuration object. Type: String.
* `description` - (Optional) Description of the Port Configuration object. Type: String.
* `annotation` - (Optional) Annotation of the Port Configuration object. Type: String.
* `name_alias` - (Optional) Name Alias of the Port Configuration object. Type: String.
