---
layout: "aci"
page_title: "ACI: aci_l3out_floating_svi"
sidebar_current: "docs-aci-data-source-l3out_floating_svi"
description: |-
  Data source for ACI L3out Floating SVI
---

# aci_l3out_floating_svi #
Data source for ACI L3out Floating SVI

## Example Usage ##

```hcl
data "aci_l3out_floating_svi" "example" {
  logical_interface_profile_dn  = "${aci_l3out_floating_svi.example.id}"
  nodeDn  = "example"
  encap  = "example"
}
```


## Argument Reference ##

* `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
* `node_dn` - (Required) Node DN of L3out floating SVI object.
* `encap` - (Required) Encap of L3out floating SVI object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Logical Interface Profile.
* `addr` - Peer address for L3out floating SVI object.
* `annotation` - Annotation for L3out floating SVI object.
* `description` - Description for L3out floating SVI object.
* `autostate` - Autostate for L3out floating SVI object.
* `encap` - Port encapsulation for L3out floating SVI object.
* `encap_scope` - Encap scope for L3out floating SVI object.
* `if_inst_t` - Interface type for L3out floating SVI object.
* `ipv6_dad` - IPv6 dad for L3out floating SVI object.
* `ll_addr` - LL address for L3out floating SVI object.
* `mac` - MAC address for L3out floating SVI object.
* `mode` - BGP domain mode for L3out floating SVI object.
* `mtu` - Administrative MTU port on the aggregated interface for L3out floating SVI object.
* `target_dscp` - Target DSCP for L3out floating SVI object.
* `userdom` - Userdom for L3out floating SVI object.
