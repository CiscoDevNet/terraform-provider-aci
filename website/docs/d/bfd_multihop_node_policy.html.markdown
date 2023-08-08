---
subcategory: -
layout: "aci"
page_title: "ACI: aci_bfdmultihop_node_policy"
sidebar_current: "docs-aci-data-source-bfdmultihop_node_policy"
description: |-
  Data source for ACI BFD Multihop Node Policy
---

# aci_bfdmultihop_node_policy #

Data source for ACI BFD Multihop Node Policy


## API Information ##

* `Class` - bfdMhNodePol
* `Distinguished Name` - uni/tn-{name}/bfdMhNodePol-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> BFD Multihop -> Node Policies


## Example Usage ##

```hcl
data "aci_bfdmultihop_node_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the BFD Multihop Node Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the BFD Multihop Node Policy. Type: String.
* `annotation` - (Read-Only) Annotation of the BFD Multihop Node Policy object. Type: String.
* `name_alias` - (Read-Only) Name Alias of the BFD Multihop Node Policy object. Type: String.
* `admin_state` - (Read-Only) Administrative state of the object or policy. Type: String.
* `detection_multiplier` - (Read-Only) Detection Multiplier. Type: String.
* `min_rx_intverval` - (Read-Only) Required Minimum RX Interval. Type: String.
* `min_tx_interval` - (Read-Only) Desired Minimum TX Interval. Type: String.
