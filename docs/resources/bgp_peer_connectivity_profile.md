---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_bgp_peer_connectivity_profile"
sidebar_current: "docs-aci-resource-bgp_peer_connectivity_profile"
description: |-
  Manages ACI BGP Peer Connectivity Profile
---

# aci_bgp_peer_connectivity_profile

Manages ACI BGP Peer Connectivity Profile

## Example Usage

```hcl
// Loopback Association
resource "aci_bgp_peer_connectivity_profile" "example" {
  parent_dn           = aci_logical_node_profile.example.id
  addr                = "10.0.0.1"
  description         = "from terraform"
  addr_t_ctrl         = ["af-mcast", "af-ucast"]
  allowed_self_as_cnt = "3"
  annotation          = "example"
  ctrl                = ["allow-self-as"]
  name_alias          = "example"
  password            = "example"
  peer_ctrl           = ["bfd"]
  private_a_sctrl     = ["remove-all", "remove-exclusive"]
  ttl                 = "1"
  weight              = "1"
  as_number           = "1"
  local_asn           = "15"
  local_asn_propagate = "dual-as"
  admin_state         = "enabled"

  relation_bgp_rs_peer_to_profile {
    direction = "import"
    target_dn = "uni/tn-tenant01/prof-test"
  }
  relation_bgp_rs_peer_to_profile {
    direction = "export"
    target_dn = "uni/tn-tenant01/prof-data"
  }
}

// Connected Association
resource "aci_bgp_peer_connectivity_profile" "example" {
  parent_dn           = aci_l3out_path_attachment.example.id
  addr                = "10.0.0.2"
  description         = "from terraform"
  addr_t_ctrl         = "af-mcast,af-ucast"
  allowed_self_as_cnt = "3"
  annotation          = "example"
  ctrl                = "allow-self-as"
  name_alias          = "example"
  password            = "example"
  peer_ctrl           = "bfd"
  private_a_sctrl     = "remove-all,remove-exclusive"
  ttl                 = "1"
  weight              = "1"
  as_number           = "1"
  local_asn           = "15"
  local_asn_propagate = "dual-as"
  admin_state         = "enabled"

  relation_bgp_rs_peer_to_profile {
    direction = "import"
    target_dn = "uni/tn-tenant01/prof-test"
  }
  relation_bgp_rs_peer_to_profile {
    direction = "export"
    target_dn = "uni/tn-tenant01/prof-data"
  }
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of parent logical node profile or L3-out Path Attachment object.
- `logical_node_profile_dn` - **Deprecated** (Required if parent_dn is not used) Distinguished name of parent logical node profile or L3-out Path Attachment object.
- `addr` - (Required) The peer IP address.
- `addr_t_ctrl` - (Optional) Ucast/Mcast Addr Type AF Control. (Multiple Comma-Delimited values are allowed. E.g., "af-mcast,af-ucast"). Apply "" to clear all the values.  
  Allowed values: "af-mcast", "af-ucast". Default value: "af-ucast".
- `admin_state` - (Optional) The administrative state of the object or policy. Allowed values are "disabled", "enabled", and default value is "enabled". Type: String.
- `allowed_self_as_cnt` - (Optional) The number of occurrences of a local Autonomous System Number (ASN). Default value: "3".
- `description` - (Optional) Description for object bgp peer connectivity profile.
- `annotation` - (Optional) Annotation for object bgp peer connectivity profile.
- `ctrl` - (Optional) The peer controls specify which Border Gateway Protocol (BGP) attributes are sent to a peer.Allowed values: "allow-self-as", "as-override", "dis-peer-as-check", "nh-self", "send-com", "send-ext-com".
- `name_alias` - (Optional) Name alias for object bgp peer connectivity profile.
- `password` - (Optional, Sensitive) Peer password. If `password` is set, the peer password will reset when terraform configuration is applied.
- `peer_ctrl` - (Optional) The peer controls. Allowed values: "bfd", "dis-conn-check". 
- `private_a_sctrl` - (Optional) Remove private AS. Allowed values: "remove-all", "remove-exclusive", "replace-as".
- `ttl` - (Optional) Specifies time to live (TTL). Default value: "1".
- `weight` - (Optional) The weight of the fault in calculating the health score of an object. A higher weight causes a higher degradation of the health score of the affected object. Default value: "0".
- `as_number` - (Optional) A number that uniquely identifies an autonomous system.
- `local_asn ` - (Optional) The local autonomous system number (ASN).
- `local_asn_propagate` - (Optional) The local Autonomous System Number (ASN) configuration.
  Allowed values: "dual-as", "no-prepend", "none", "replace-as". Default value: "none".
- `relation_bgp_rs_peer_pfx_pol` - (Optional) Relation to class bgpPeerPfxPol. Cardinality - N_TO_ONE. Type - String.
- `relation_bgp_rs_peer_to_profile` - (Optional) A block representing the relation to a Route Control Profile (class rtctrlProfile). Type: Block.
  * `direction` - (Optional) The connector direction. Allowed values are "export", "import", and default value is "import". Type: String.
  * `target_dn` - (Required) The distinguished name of the target Route Map for Route Control Profile. Type: String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Peer Connectivity Profile.

## Importing

An existing BGP Peer Connectivity Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_peer_connectivity_profile.example <Dn>
```
