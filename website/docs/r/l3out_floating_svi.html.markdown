---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_floating_svi"
sidebar_current: "docs-aci-resource-l3out_floating_svi"
description: |-
  Manages ACI L3out Floating SVI
---

# aci_l3out_floating_svi

Manages ACI L3out Floating SVI

## Example Usage

```hcl
resource "aci_l3out_floating_svi" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  node_dn                      = "topology/pod-1/node-201"
  encap                        = "vlan-20"
  addr                         = "10.20.30.40/16"
  annotation                   = "example"
  description                  = "from terraform"
  autostate                    = "enabled"
  encap_scope                  = "ctx"
  if_inst_t                    = "ext-svi"
  ipv6_dad                     = "disabled"
  ll_addr                      = "::"
  mac                          = "12:23:34:45:56:67"
  mode                         = "untagged"
  mtu                          = "580"
  target_dscp                  = "CS1"
  relation_l3ext_rs_dyn_path_att {
    tdn = data.aci_physical_domain.dom.id
    floating_address = "10.21.0.254/24"
    forged_transmit = "Disabled"
    mac_change = "Disabled"
    promiscuous_mode = "Disabled"
  }
}
```

## Argument Reference

* `logical_interface_profile_dn` - (Required) Distinguished name of the parent logical interface profile object.
* `node_dn` - (Required) Distinguished name of the node of the L3out floating SVI object.
* `encap` - (Required) Port encapsulation of the L3out floating SVI object.
* `addr` - (Optional) Peer address of the L3out floating SVI object. Default value: "0.0.0.0".
* `annotation` - (Optional) Annotation of the L3out floating SVI object.
* `description` - (Optional) Description of the L3out floating SVI object.
* `autostate` - (Optional) Autostate of the L3out floating SVI object. Allowed values are "disabled" and "enabled". Default value is "disabled".
* `encap_scope` - (Optional) Encap scope of the L3out floating SVI object. Allowed values are "ctx" and "local". Default value is "local".
* `if_inst_t` - (Optional) Interface type of the L3out floating SVI object. Allowed values are "ext-svi", "l3-port", "sub-interface" and "unspecified". Default value is "unspecified".
* `ipv6_dad` - (Optional) IPv6 dad of the L3out floating SVI object. Allowed values are "disabled" and "enabled". Default value is "enabled".
* `ll_addr` - (Optional) Link local address of the L3out floating SVI object. Default value: "::".
* `mac` - (Optional) MAC address of the L3out floating SVI object.
* `mode` - (Optional) BGP domain mode of the L3out floating SVI object. Allowed values are "native", "regular" and "untagged". Default value is "regular".
* `mtu` - (Optional) Administrative MTU port on the aggregated interface of the L3out floating SVI object. Range of allowed values is "576" to "9216". Default value is "inherit".
* `target_dscp` - (Optional) Target DSCP of the L3out floating SVI object. Allowed values are "AF11", "AF12", "AF13", "AF21", "AF22", "AF23", "AF31", "AF32", "AF33", "AF41", "AF42", "AF43", "CS0", "CS1", "CS2", "CS3", "CS4", "CS5", "CS6", "CS7", "EF", "VA" and "unspecified". Default value is "unspecified".
* `relation_l3ext_rs_dyn_path_att` - (Optional) A block representing the relation to a Domain (class infraDomP or vmmDomP). Type: Block.
  * `tdn` - (Required) The distinguished name of the target.
  * `floating_address` - (Optional) Floating address of the target.
  * `forged_transmit` - (Optional) A configuration option that allows virtual machines to send frames with a mac address that is different from the one specified in the virtual-nic adapter configuration. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `mac_change` - (Optional) The status of the mac address change support of the port groups in an external VMM controller, such as a vCenter. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `promiscuous_mode` - (Optional) The status of promiscuous mode support status of the port groups in an external VMM controller, such as a vCenter. This needs to be turned on only for service devices in the cloud, not for Enterprise AVE service deployments. Allowed values are "Disabled" and "Enabled". Default value is "Disabled".
  * `enhanced_lag_policy_dn` - (Optional) The distinguished name of the target enhanced lag policy (class lacpEnhancedLagPol).

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out Floating SVI.

## Importing

An existing L3out Floating SVI can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_floating_svi.example <Dn>
```
