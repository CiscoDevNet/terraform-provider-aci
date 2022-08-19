---
subcategory: "System Settings"
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
* `Distinguished Name` - uni/infra/settings

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
* `domain_validation` - (Optional) Validate that the correct physical domain is added before associating a new static path to an EPG.
* `enable_mo_streaming` - (Optional) Enable MO streaming of object Fabric-Wide Settings Policy.  
* `enable_remote_leaf_direct` - (Optional) Enable remote leaf direct communication of object Fabric-Wide Settings Policy.
* `enforce_subnet_check` - (Optional) Enforce subnet check of object Fabric-Wide Settings Policy.  
* `leaf_opflexp_authenticate_clients` - (Optional) Require Opflexp Client Certificates for authentication for Leaf.
* `leaf_opflexp_use_ssl` - (Optional) Require SSL transport for Opflexp for Leaf.
* `opflexp_authenticate_clients` - (Optional) Opflexp Client Certificates for authentication of object Fabric-Wide Settings Policy.  
* `opflexp_ssl_protocols` - (Optional) SSL Opflex versions.
* `opflexp_use_ssl` - (Optional) SSL transport for Opflexp indicator of object Fabric-Wide Settings Policy. 
* `policy_sync_node_bringup` - (Optional) Blacklist the Leaf frontpanel port until policy download during first time bringup.
* `reallocate_gipo` - (Optional) Reallocate gipo such that stretched and non stretched BDs have non overlapping gipos.
* `restrict_infra_vlan_traffic` - (Optional) Intra Leaf Communication traffic indicator of object Fabric-Wide Settings Policy. 
* `unicast_xr_ep_learn_disable` - (Optional) Disable xrLeanrs indicator of object Fabric-Wide Settings Policy. 
* `validate_overlapping_vlans` - (Optional) Validate Overlapping VLANS indicator of object Fabric-Wide Settings Policy.