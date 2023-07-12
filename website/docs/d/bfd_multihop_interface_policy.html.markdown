---
layout: "aci"
page_title: "ACI: aci_bfd_multihop_interface_policy"
sidebar_current: "docs-aci-data-source-bfd_multihop_interface_policy"
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
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the  BFD Multihop Interface Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the BFD Multihop Interface Policy object.
* `annotation` - (Optional) Annotation of the BFD Multihop Interface Policy object.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Interface Policy object.
* `admin_state` - (Optional) Enable Disable sessions. The administrative state of the object or policy.
* `detection_multiplier` - (Optional) Detection Multiplier. Detection multiplier.
* `min_rx_intvl` - (Optional) Required Minimum RX Interval. Required minimum rx interval.
* `min_transmit_interval` - (Optional) Desired Minimum TX Interval. Desired minimum tx interval.
