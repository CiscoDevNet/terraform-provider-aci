---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bfd_multihop_interface_policy"
sidebar_current: "docs-aci-data-source-aci_bfd_multihop_interface_policy"
description: |-
  Data source for ACI BFD Multihop Interface Policy
---

# aci_bfd_multihop_interface_policy #

Data source for ACI BFD Multihop Interface Policy


## API Information ##

* `Class` - bfdMhIfPol
* `Distinguished Name` - uni/tn-{tn_name}/bfdMhIfPol-{policy_name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> BFD Multihop -> Interface Policies



## Example Usage ##

```hcl
data "aci_bfd_multihop_interface_policy" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the  BFD Multihop Interface Policy object. Type: String.

## Attribute Reference ##
* `id` - (Read-Only) Attribute id set to the Dn of the BFD Multihop Interface Policy object. Type: String.
* `annotation` - (Read-Only) Annotation of the BFD Multihop Interface Policy object. Type: String.
* `name_alias` - (Read-Only) Name Alias of the BFD Multihop Interface Policy object. Type: String.
* `admin_state` - (Read-Only) The administrative state of the object or policy. Type: String.
* `description` - (Read-Only) Description for the BFD Multihop Interface Policy object. Type: String.
* `detection_multiplier` - (Read-Only) Detection Multiplier. Type: String.
* `min_receive_interval` - (Read-Only) Required Minimum Rx Interval.  Type: String.
* `min_transmit_interval` - (Read-Only) Desired Minimum Tx Interval. Type: String.
