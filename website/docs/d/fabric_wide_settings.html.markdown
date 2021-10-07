---
layout: "aci"
page_title: "ACI: aci_fabric_wide_settings"
sidebar_current: "docs-aci-data-source-fabric_wide_settings"
description: |-
  Data source for ACI Fabric-Wide Settings Policy
---

# aci_fabric_wide_settings #
Data source for ACI Fabric-Wide Settings Policy


## API Information ##
* `Class` - infraSetPol
* `Distinguished Named` - uni/infra/settings

## GUI Information ##
* `Location` - System -> System Settings -> Fabric-Wide Settings 


## Example Usage ##

```hcl
data "aci_fabric_wide_settings" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Fabric-Wide Settings Policy.
* `name` - (Optional) Name of object Fabric-Wide Settings Policy.
* `annotation` - (Optional) Annotation of object Fabric-Wide Settings Policy.
* `description` - (Optional) Description of object Fabric-Wide Settings Policy.
* `name_alias` - (Optional) Name alias of object Fabric-Wide Settings Policy.
* `disable_ep_dampening` - (Optional) Disable Ep Dampening knob of object Fabric-Wide Settings Policy. 
* `enable_mo_streaming` - (Optional) Enable MO streaming of object Fabric-Wide Settings Policy.  
* `enable_remote_leaf_direct` - (Optional) Enable remote leaf direct communication of object Fabric-Wide Settings Policy.
* `enforce_subnet_check` - (Optional) Enforce subnet check of object Fabric-Wide Settings Policy.  
* `opflexp_authenticate_clients` - (Optional) Opflexp Client Certificates for authentication of object Fabric-Wide Settings Policy.  
* `opflexp_use_ssl` - (Optional) SSL transport for Opflexp indicator of object Fabric-Wide Settings Policy. 
* `restrict_infra_vlan_traffic` - (Optional) Intra Leaf Communication traffic indicator of object Fabric-Wide Settings Policy. 
* `unicast_xr_ep_learn_disable` - (Optional) Disable xrLeanrs indicator of object Fabric-Wide Settings Policy. 
* `validate_overlapping_vlans` - (Optional) Validate Overlapping VLANS indicator of object Fabric-Wide Settings Policy.