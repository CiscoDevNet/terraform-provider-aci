---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_floating_svi"
sidebar_current: "docs-aci-data-source-l3out_floating_svi"
description: |-
  Data source for ACI L3out Floating SVI
---

# aci_l3out_floating_svi

Data source for ACI L3out Floating SVI

## Example Usage

```hcl
data "aci_l3out_floating_svi" "check" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  node_dn                      = "topology/pod-1/node-201"
  encap                        = "vlan-20"
}
```

## Argument Reference

* `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
* `node_dn` - (Required) Node DN of L3out floating SVI object.
* `encap` - (Required) Port encapsulation for L3out floating SVI object.

## Attribute Reference

* `id` - Attribute id set to the Dn of the Logical Interface Profile.
* `addr` - (Optional) Peer address for L3out floating SVI object.
* `annotation` - (Optional) Annotation for L3out floating SVI object.
* `description` - (Optional) Description for L3out floating SVI object.
* `autostate` - (Optional) Autostate for L3out floating SVI object.
* `encap_scope` - (Optional) Encap scope for L3out floating SVI object.
* `if_inst_t` - (Optional) Interface type for L3out floating SVI object.
* `ipv6_dad` - (Optional) IPv6 dad for L3out floating SVI object.
* `ll_addr` - (Optional) Link local address address for L3out floating SVI object.
* `mac` - (Optional) MAC address for L3out floating SVI object.
* `mode` - (Optional) BGP domain mode for L3out floating SVI object.
* `mtu` -( Optional) Administrative MTU port on the aggregated interface for L3out floating SVI object.
* `target_dscp` - (Optional) Target DSCP for L3out floating SVI object.
* `relation_l3ext_rs_dyn_path_att` - (Optional) A block representing the relation to a Domain (class infraDomP or vmmDomP). Type: Block.
  * `tdn` - (Required) The distinguished name of the target.
  * `floating_address` - (Optional) Floating address of the target.
  * `forged_transmit` - (Optional) A configuration option that allows virtual machines to send frames with a mac address that is different from the one specified in the virtual-nic adapter configuration. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `mac_change` - (Optional) The status of the mac address change support for port groups in an external VMM controller, such as a vCenter. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `promiscuous_mode` - (Optional) The status of promiscuous mode support status for port groups in an external VMM controller, such as a vCenter. This needs to be turned on only for service devices in the cloud, not for Enterprise AVE service deployments. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `enhanced_lag_policy_tdn` - (Optional) The distinguished name of the target enhanced lag policy (class lacpEnhancedLagPol).