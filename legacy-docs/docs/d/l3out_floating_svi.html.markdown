---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_floating_svi"
sidebar_current: "docs-aci-data-source-aci_l3out_floating_svi"
description: |-
  Data source for ACI L3out Floating SVI
---

# aci_l3out_floating_svi

Data source for ACI L3out Floating SVI

## API Information ##

* `Class` - l3extVirtualLIfP
* `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/lnodep-{lnodep}/lifp-{lifp}/vlifp-[nodeDn]-[encap]

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles -> Floating SVI

## Example Usage

```hcl
data "aci_l3out_floating_svi" "check" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  node_dn                      = "topology/pod-1/node-201"
  encap                        = "vlan-20"
}
```

## Argument Reference

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent logical interface profile object.
* `node_dn` - (Required) Node DN of the L3out floating SVI object.
* `encap` - (Required) Port encapsulation of the L3out floating SVI object.

## Attribute Reference

* `id` - Attribute id set to the Dn of the Logical Interface Profile.
* `addr` - (Read-Only) Peer address of the L3out floating SVI object.
* `annotation` - (Read-Only) Annotation of the L3out floating SVI object.
* `description` - (Read-Only) Description of the L3out floating SVI object.
* `autostate` - (Read-Only) Autostate of the L3out floating SVI object.
* `encap_scope` - (Read-Only) Encap scope of the L3out floating SVI object.
* `if_inst_t` - (Read-Only) Interface type of the L3out floating SVI object.
* `ipv6_dad` - (Read-Only) IPv6 dad of the L3out floating SVI object.
* `ll_addr` - (Read-Only) Link local address address of the L3out floating SVI object.
* `mac` - (Read-Only) MAC address of the L3out floating SVI object.
* `mode` - (Read-Only) BGP domain mode of the L3out floating SVI object.
* `mtu` - (Read-Only) Administrative MTU port on the aggregated interface of the L3out floating SVI object.
* `target_dscp` - (Read-Only) Target DSCP of the L3out floating SVI object.
* `relation_l3ext_rs_dyn_path_att` - (Read-Only) A block representing the relation to a Domain (class infraDomP or vmmDomP). Type: Block.
  * `tdn` - (Read-Only) The distinguished name of the target.
  * `floating_address` - (Read-Only) Floating address of the target.
  * `forged_transmit` - (Read-Only) A configuration option that allows virtual machines to send frames with a mac address that is different from the one specified in the virtual-nic adapter configuration.
  * `mac_change` - (Read-Only) The status of the mac address change support of the port groups in an external VMM controller, such as a vCenter.
  * `promiscuous_mode` - (Read-Only) The status of the promiscuous mode support status of the port groups in an external VMM controller, such as a vCenter. This needs to be turned on only for service devices in the cloud, not for Enterprise AVE service deployments.
  * `enhanced_lag_policy_dn` - (Read-Only) The distinguished name of the target enhanced lag policy (class lacpEnhancedLagPol).
  * `encap` - (Read-Only) Access port encapsulation (VLAN) of the target (format: vlan-101).