---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_epg_def"
sidebar_current: "docs-aci-data-source-epg_def"
description: |-
  Data source for ACI EPg Def
---

# aci_epg_def #

Data source for ACI EPg Def

## API Information ##

* `Class` - vnsEPgDef
* `Distinguished Name` - uni/tn-{tenant_name}/GraphInst_C-[{ctrctDn}]-G-[{graphDn}]-S-[{scopeDn}]/NodeInst-{node_name}/LegVNode-{virtual_node_name}/EPgDef-{name}

## Example Usage ##

```hcl
data "aci_epg_def" "example" {
  logical_context_dn = "uni/tn-Symmetric-PBR/ldevCtx-c-intra-web-g-FW-n-node1/lIfCtx-c-consumer"
}
```

## Argument Reference ##
* `logical_context_dn` - (Required) Distinguished name of the Logical Interface Context.
* 
## Attribute Reference ##
* `id` - Attribute id set to the Dn of EPgDef.
* `legacy_virtual_node_dn` - Distinguished name of the parent Legacy Virtual Node object.
* `name` - Name of the EPgDef object.
* `annotation` - Annotation of the EPgDef object.
* `name_alias` - Name Alias of the EPgDef object.
* `encap` - The VLAN encapsulation tag.
* `fabric_encap` - The VXLAN encapsulation tag.
* `delete_pbr_scenario` - The boolean value for deleting Policy Based Routing. 
* `member_type` - The type of member for the EPgDef object.
* `router_id` - The IP address of the routing device.
