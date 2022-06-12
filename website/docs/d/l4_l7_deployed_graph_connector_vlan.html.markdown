---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_deployed_graph_connector_vlan"
sidebar_current: "docs-aci-data-source-l4_l7_deployed_graph_connector_vlan"
description: |-
  Data Source for ACI L4-L7 Deployed Graph Connector VLAN
---

# aci_epg_def #

Data Source for ACI L4-L7 Deployed Graph Connector VLAN

## API Information ##

* `Class` - vnsEPgDef
* `Distinguished Name` - uni/tn-{tenant_name}/GraphInst_C-[{ctrctDn}]-G-[{graphDn}]-S-[{scopeDn}]/NodeInst-{node_name}/LegVNode-{virtual_node_name}/EPgDef-{name}

## GUI Information ##

* `Location` -  Tenant -> {tenant_name} -> Services -> L4-L7 -> Deployed Graph Instance -> {graph_name} -> Function Connector

## Example Usage ##

```hcl
data "aci_l4_l7_deployed_graph_connector_vlan" "example" {
  logical_context_dn = "uni/tn-Symmetric-PBR/ldevCtx-c-intra-web-g-FW-n-node1/lIfCtx-c-consumer"
}
```

## Argument Reference ##
* `logical_context_dn` - (Required) Distinguished name of the Logical Interface Context.
* 
## Attribute Reference ##
* `id` - Attribute id set to the Dn of EPgDef.
* `name` - Name of the EPgDef object.
* `annotation` - Annotation of the EPgDef object.
* `name_alias` - Name Alias of the EPgDef object.
* `encap` - The VLAN encapsulation tag.
* `fabric_encap` - The VXLAN encapsulation tag.
* `delete_pbr_scenario` - The boolean value for deleting Policy Based Routing. 
* `member_type` - The type of member for the EPgDef object.
* `router_id` - The IP address of the routing device.
