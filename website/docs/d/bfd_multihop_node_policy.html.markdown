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

* `Location` - 


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
* `annotation` - (Optional) Annotation of the BFD Multihop Node Policy object. Type: String.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Node Policy object. Type: String.
* `admin_state` - (Optional) Enable Disable sessions. The administrative state of the object or policy. Type: String.
* `detection_multiplier` - (Optional) Detection Multiplier. Allowed range is 1-50 and default value is "3". Type: String.
* `min_rx_intverval` - (Optional) Required Minimum RX Interval. Allowed range is 250-999 and default value is "250". Type: String.
* `min_tx_interval` - (Optional) Desired Minimum TX Interval. Allowed range is 250-999 and default value is "250". Type: String.
