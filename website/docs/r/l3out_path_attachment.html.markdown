---
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment"
sidebar_current: "docs-aci-resource-l3out_path_attachment"
description: |-
  Manages ACI L3-out Path Attachment
---

# aci_l3out_path_attachment

Manages ACI L3-out Path Attachment

## Example Usage

```hcl
resource "aci_l3out_path_attachment" "example" {
  logical_interface_profile_dn  = aci_logical_interface_profile.example.id
  target_dn  = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t = "ext-svi"
  description = "from terraform"
  addr  = "10.20.30.40/16"
  annotation  = "example"
  autostate = "disabled"
  encap  = "vlan-1"
  encap_scope = "ctx"
  ipv6_dad = "disabled"
  ll_addr  = "::"
  mac  = "0F:0F:0F:0F:FF:FF"
  mode = "native"
  mtu = "inherit"
  target_dscp = "AF11"
}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `target_dn` - (Required) The logical interface identifier for object L3 out Path Attachment.
- `if_inst_t` - (Required) Interface type for object L3 out Path Attachment. Allowed values: "ext-svi", "l3-port", "sub-interface", "unspecified".
- `addr` - (Optional) The IP address of the path attached to the layer 3 outside profile. Default value is "0.0.0.0".
- `description` - (Optional) Description for object L3 out Path Attachment.
- `annotation` - (Optional) Annotation for object L3 out Path Attachment.
- `autostate` - (Optional) Autostate for object L3 out Path Attachment.
  Allowed values: "disabled", "enabled". Default value is "disabled".
- `encap` - (Optional) The encapsulation of the path attached to the layer 3 outside profile. Encap should be set to "unknown" if the value of if_inst_t is "l3-port". Default value is "unknown".
- `encap_scope` - (Optional) The encapsulation scope for object L3 out Path Attachment. Allowed values: "ctx", "local". Default value is "local".
- `ipv6_dad` - (Optional) IPv6 DAD for object L3 out Path Attachment.
  Allowed values: "disabled", "enabled". Default value is "enabled".
- `ll_addr` - (Optional) The override of the system generated IPv6 link-local address. Default value is "::".
- `mac` - (Optional) The MAC address of the path attached to the layer 3 outside profile. Default value is "00:22:BD:F8:19:FF".
- `mode` - (Optional) BGP Domain mode. Allowed values: "native", "regular", "untagged". Default value is "regular".
- `mtu` - (Optional) The maximum transmit unit of the external network. Default value is "inherit".
- `target_dscp` - (Optional) The target differentiated service code point (DSCP) of the path attached to the layer 3 outside profile. Default value: "unspecified". Allowed values: "AF11", "AF12", "AF13", "AF21", "AF22", "AF23", "AF31", "AF32", "AF33", "AF41", "AF42", "AF43", "CS0", "CS1", "CS2", "CS3", "CS4", "CS5", "CS6", "CS7", "EF", "VA", "unspecified". Default value: "unspecified".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3-out Path Attachment.

## Importing

An existing L3-out Path Attachment can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_path_attachment.example <Dn>
```
