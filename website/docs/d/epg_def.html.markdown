---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_epg_def"
sidebar_current: "docs-aci-data-source-epg_def"
description: |-
  Data source for ACI EPg Def
---

# aci_e_pg_def #

Data source for ACI EPg Def

## API Information ##

* `Class` - vnsEPgDef
* `Distinguished Name` - uni/tn-{name}/GraphInst_C-[{ctrctDn}]-G-[{graphDn}]-S-[{scopeDn}]/NodeInst-{name}/LegVNode-{name}/EPgDef-{name}

## Example Usage ##

```hcl
data "aci_e_pg_def" "example" {
  name                    = "example"
  legacy_virtual_node_dn  = join("", ["uni/tn-Symmetric-PBR/GraphInst_C-[uni/tn-Symmetric-PBR/brc-intra-web]",
  "-G-[uni/tn-Symmetric-PBR/AbsGraph-FW]-S-[uni/tn-Symmetric-PBR/ctx-vrf1]/NodeInst-node1/LegVNode-0"])
}
```

## Argument Reference ##

* `legacy_virtual_node_dn` - (Required) Distinguished name of the parent Legacy Virtual Node object.
* `name` - (Required) Name of the EPgDef object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of EPgDef.
* `annotation` - Annotation of the EPgDef object.
* `name_alias` - Name Alias of the EPgDef object.
* `encap` - The VLAN encapsulation tag.
* `fabric_encap` - The VXLAN encapsulation tag.
* `delete_pbr_scenario` - The boolean value for deleting Policy Based Routing. 
* `member_type` - The type of member for the EPgDef object.
* `logical_context_dn` - Distinguished name of the Logical Interface Context.
* `router_id` - The IP address of the routing device.
