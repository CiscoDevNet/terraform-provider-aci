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

  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
  target_dn  = "topology/pod-1/paths-101/pathep-[eth1/1]"
  if_inst_t = "ext-svi"
  addr  = "0.0.0.0"
  annotation  = "example"
  autostate = "disabled"
  encap  = "vlan-1"
  encap_scope = "ctx"
  ipv6_dad = "disabled"
  ll_addr  = "example"
  mac  = "0F:0F:0F:0F:FF:FF"
  mode = "native"
  mtu = "inherit"
  target_dscp = "AF11"

}
```

## Argument Reference

- `logical_interface_profile_dn` - (Required) Distinguished name of parent logical interface profile object.
- `target_dn` - (Required) The logical interface identifier.
- `if_inst_t` - (Required) Interface type.  
  Allowed values: "ext-svi", "l3-port", "sub-interface", "unspecified".
- `addr` - (Optional) The IP address of the path attached to the layer 3 outside profile.

- `annotation` - (Optional) Annotation for object L3-out Path Attachment.

- `autostate` - (Optional) Autostate for object L3-out Path Attachment.
  Allowed values: "disabled", "enabled". Default value: "disabled".
- `encap` - (Optional) The encapsulation of the path attached to the layer 3 outside profile.

- `encap_scope` - (Optional) The encapsulation scope for object L3-out Path Attachment.
  Allowed values: "ctx", "local". Default value: "local".

- `ipv6_dad` - (Optional) IPv6-Dad for object L3-out Path Attachment.
  Allowed values: "disabled", "enabled". Default value: "enabled".
- `ll_addr` - (Optional) The override of the system generated IPv6 link-local address.

- `mac` - (Optional) The MAC address of the path attached to the layer 3 outside profile.

- `mode` - (Optional) BGP Domain mode.  
  Allowed values: "native", "regular", "untagged". Default value: "regular".
- `mtu` - (Optional) The maximum transmit unit of the external network.
  Allowed value: "inherit".

- `target_dscp` - (Optional) The target differentiated service code point (DSCP) of the path attached to the layer 3 outside profile. Default value: "unspecified".
  Allowed values: "AF11", "AF12", "AF13", "AF21", "AF22", "AF23", "AF31", "AF32", "AF33", "AF41", "AF42", "AF43", "CS0", "CS1", "CS2", "CS3", "CS4", "CS5", "CS6", "CS7", "EF", "VA", "unspecified". Default value: "unspecified".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3-out Path Attachment.

## Importing

An existing L3-out Path Attachment can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_path_attachment.example <Dn>
```
