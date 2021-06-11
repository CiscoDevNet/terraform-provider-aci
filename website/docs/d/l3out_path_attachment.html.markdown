---
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment"
sidebar_current: "docs-aci-data-source-l3out_path_attachment"
description: |-
  Data source for ACI L3-out Path Attachment
---

# aci_l3out_path_attachment

Data source for ACI L3-out Path Attachment

## Example Usage

```hcl
data "aci_l3out_path_attachment" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
  target_dn  = "example"
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `target_dn` - (Required) The logical interface identifier.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3 out Path Attachment.
- `addr` - (Optional) The IP address of the path attached to the layer 3 outside profile.
- `description` - (Optional) Description for object L3 out Path Attachment.
- `annotation` - (Optional) Annotation for object L3 out Path Attachment.
- `autostate` - (Optional) Autostate for object L3 out Path Attachment.
- `encap` - (Optional) The encapsulation of the path attached to the layer 3 outside profile.
- `encap_scope` - (Optional) The encapsulation scope for object L3 out Path Attachment.
- `if_inst_t` - (Optional) Interface type.
- `ipv6_dad` - (Optional) IPv6 DAD for object L3 out Path Attachment.
- `ll_addr` - (Optional) The override of the system generated IPv6 link local address.
- `mac` - (Optional) The MAC address of the path attached to the layer 3 outside profile.
- `mode` - (Optional) BGP Domain mode.
- `mtu` - (Optional) The maximum transmit unit of the external network.
- `target_dscp` - (Optional) The target differentiated service code point (DSCP) of the path attached to the layer 3 outside profile.
