---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_fabric_wide_settings"
sidebar_current: "docs-aci-resource-fabric_wide_settings"
description: |-
  Manages ACI Fabric-Wide Settings Policy
---

# aci_fabric_wide_settings #
Manages ACI Fabric-Wide Settings Policy

## API Information ##
* `Class` - infraSetPol
* `Distinguished Name` - uni/infra/settings

## GUI Information ##
* `Location` - System -> System Settings -> Fabric-Wide Settings 

## Example Usage ##

```hcl
resource "aci_fabric_wide_settings" "example" {
  name = "example"
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name_alias = "example_name_alias"
  disable_ep_dampening = "yes"
  domain_validation = "yes"
  enable_mo_streaming = "yes"
  enable_remote_leaf_direct = "yes"
  enforce_subnet_check = "yes"
  leaf_opflexp_authenticate_clients = "yes"
  leaf_opflexp_use_ssl = "yes"
  opflexp_authenticate_clients = "yes"
  opflexp_ssl_protocols = "TLSv1,TLSv1.1,TLSv1.2"
  opflexp_use_ssl = "yes"
  policy_sync_node_bringup = "yes"
  reallocate_gipo = "yes"
  restrict_infra_vlan_traffic = "yes"
  unicast_xr_ep_learn_disable = "yes"
  validate_overlapping_vlans = "yes"
}
```

## NOTE ##
Users can use the resource of type `aci_fabric_wide_settings` to change the configuration of the object Fabric Wide Settings. Users cannot create more than one instance of object Fabric Wide Settings.

## Argument Reference ##
* `name` - (Optional) Name of object Fabric-Wide Settings Policy.
* `annotation` - (Optional) Annotation of object Fabric-Wide Settings Policy.
* `description` - (Optional) Description of object Fabric-Wide Settings Policy.
* `name_alias` - (Optional) Name alias of object Fabric-Wide Settings Policy.
* `disable_ep_dampening` - (Optional) Disable Ep Dampening knob of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `domain_validation` - (Optional) Validate that static path is added but no domain is associated to an EPG.
* `enable_mo_streaming` - (Optional) Enable MO streaming of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `enable_remote_leaf_direct` - (Optional) Enable remote leaf direct communication of object Fabric-Wide Settings Policy.  Allowed values are "yes" and "no". 
* `enforce_subnet_check` - (Optional) Enforce subnet check of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `leaf_opflexp_authenticate_clients` - (Optional) Require Opflexp Client Certificates for authentication for Leaf.
* `leaf_opflexp_use_ssl` - (Optional) Require SSL transport for Opflexp for Leaf.
* `opflexp_authenticate_clients` - (Optional) Opflexp Client Certificates for authentication of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `opflexp_ssl_protocols` - (Optional) SSL Opflex versions.
* `opflexp_use_ssl` - (Optional) SSL transport for Opflexp indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `policy_sync_node_bringup` - (Optional) Blacklist the Leaf frontpanel port until policy download during first time bringup.
* `reallocate_gipo` - (Optional) Reallocate gipo such that stretched and non stretched BDs have non overlapping gipos.
* `restrict_infra_vlan_traffic` - (Optional) Intra Leaf Communication traffic indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no". (Note: attribute restrict_infra_vlan_traffic is supported for version 5 and above of APIC)
* `unicast_xr_ep_learn_disable` - (Optional) Disable xrLeanrs indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".
* `validate_overlapping_vlans` - (Optional) Validate Overlapping VLANS indicator of object Fabric-Wide Settings Policy. Allowed values are "yes" and "no".


## Importing ##

An existing Fabric-WideSettings Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_wide_settings.example <Dn>
```