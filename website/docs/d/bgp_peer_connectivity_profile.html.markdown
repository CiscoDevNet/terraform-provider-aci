---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_bgp_peer_connectivity_profile"
sidebar_current: "docs-aci-data-source-bgp_peer_connectivity_profile"
description: |-
  Data source for ACI BGP Peer Connectivity Profile
---

# aci_bgp_peer_connectivity_profile

Data source for ACI BGP Peer Connectivity Profile

## Example Usage

```hcl
// Loopback Association
data "aci_bgp_peer_connectivity_profile" "example" {
  parent_dn  = aci_logical_node_profile.example.id
  addr  = "10.0.0.1"
}

// Connected Association
data "aci_bgp_peer_connectivity_profile" "example" {
  parent_dn  = aci_l3out_path_attachment.example.id
  addr  = "10.0.0.2"
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of parent logical node profile or L3-out Path Attachment object.
- `addr` - (Required) The peer IP address.

## Attribute Reference

- `id` - Attribute id set to the Dn of the BGP Peer Connectivity Profile.
- `addr_t_ctrl` - (Optional) Ucast/Mcast Address Type AF Control.
- `description` - Description for the object of the BGP Peer Connectivity Profile.
- `allowed_self_as_cnt` - (Optional) The number of occurrences of a local Autonomous System Number (ASN).
- `annotation` - (Optional) Annotation for object bgp peer connectivity profile.
- `ctrl` - (Optional) The peer controls specify which Border Gateway Protocol (BGP) attributes are sent to a peer.
- `name_alias` - (Optional) Name alias for object bgp peer connectivity profile.
- `password` - (Optional, Sensitive) Peer password. Value is always "" to maintain confidentiality.
- `peer_ctrl` - (Optional) The peer controls.
- `private_a_sctrl` - (Optional) Remove private AS.
- `ttl` - (Optional) Specifies time to live (TTL).
- `weight` - (Optional) The weight of the fault in calculating the health score of an object. A higher weight causes a higher degradation of the health score of the affected object.
- `as_number` - (Optional) A number that uniquely identifies an autonomous system.
- `local_asn ` - (Optional) The local autonomous system number (ASN).
- `local_asn_propagate` - (Optional) The local Autonomous System Number (ASN) configuration.
