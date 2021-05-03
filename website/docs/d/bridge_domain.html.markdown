---
layout: "aci"
page_title: "ACI: aci_bridge_domain"
sidebar_current: "docs-aci-data-source-bridge_domain"
description: |-
  Data source for ACI Bridge Domain
---

# aci_bridge_domain

Data source for ACI Bridge Domain

## Example Usage

```hcl
data "aci_bridge_domain" "dev_bd" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "foo_bd"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) name of Object bridge_domain.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Bridge Domain.
- `optimize_wan_bandwidth` - (Optional) Flag to enable OptimizeWanBandwidth between sites.
- `annotation` - (Optional) Annotation for object bridge domain.
- `description` - (Optional) Description for object bridge domain.
- `arp_flood` - (Optional) A property to specify whether ARP flooding is enabled. If flooding is disabled, unicast routing will be performed on the target IP address.
- `ep_clear` - (Optional) Represents the parameter used by the node (i.e. Leaf) to clear all EPs in all leaves for this BD.
- `ep_move_detect_mode` - (Optional) The End Point move detection option uses the Gratuitous Address Resolution Protocol (GARP). A gratuitous ARP is an ARP broadcast-type of packet that is used to verify that no other device on the network has the same IP address as the sending device.
- `host_based_routing` - (Optional) Enables advertising host routes out of l3outs of this BD.
- `intersite_bum_traffic_allow` - (Optional) Control whether BUM traffic is allowed between sites.
- `intersite_l2_stretch` - (Optional) Flag to enable l2Stretch between sites.
- `ip_learning` - (Optional) Endpoint Dataplane Learning.
- `ipv6_mcast_allow` - (Optional) Flag to indicate multicast IpV6 is allowed or not.
- `limit_ip_learn_to_subnets` - (Optional) Limits IP address learning to the bridge domain subnets only. Every BD can have multiple subnets associated with it. By default, all IPs are learned.
- `ll_addr` - (Optional) Override of system generated ipv6 link-local address.
- `mac` - (Optional) The MAC address of the bridge domain (BD) or switched virtual interface (SVI). Every BD by default takes the fabric-wide default MAC address. You can override that address with a different one. By default the BD will take a 00:22:BD:F8:19:FF mac address.
- `mcast_allow` - (Optional) Flag to indicate if multicast is enabled for IpV4 addresses.
- `multi_dst_pkt_act` - (Optional) The multiple destination forwarding method for L2 Multicast, Broadcast, and Link Layer traffic types.
- `name_alias` - (Optional) Name alias for object bridge_domain.
- `bridge_domain_type` - (Optional) The specific type of the object or component.
- `unicast_route` - (Optional) The forwarding method based on predefined forwarding criteria (IP or MAC address).
- `unk_mac_ucast_act` - (Optional) The forwarding method for unknown layer 2 destinations.
- `unk_mcast_act` - (Optional) The parameter used by the node (i.e. a leaf) for forwarding data for an unknown multicast destination.
- `v6unk_mcast_act` - (Optional) M-cast action for object bridge_domain.
- `vmac` - (Optional) Virtual MAC address of the BD/SVI. This is used when the BD is extended to multiple sites using l2 Outside. Only allowed values is "not-applicable".
