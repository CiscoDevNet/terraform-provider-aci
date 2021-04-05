---
layout: "aci"
page_title: "ACI: aci_l3out_bgp_protocol_profile"
sidebar_current: "docs-aci-resource-l3out_bgp_protocol_profile"
description: |-
  Manages ACI L3out BGP Protocol Profile
---

# aci_l3out_bgp_protocol_profile

Manages ACI L3out BGP Protocol Profile

## Example Usage

```hcl
resource "aci_l3out_bgp_protocol_profile" "example" {

  logical_node_profile_dn  = "${aci_logical_node_profile.example.id}"
  annotation  = "example"
  name_alias  = "example"

}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of parent logical node profile object.
- `annotation` - (Optional) Annotation for object L3out BGP Protocol Profile.
- `name_alias` - (Optional) Name alias for object L3out BGP Protocol Profile.
- `relation_bgp_rs_bgp_node_ctx_pol` - (Optional) Relation to class bgpCtxPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out BGP Protocol Profile.

## Importing

An existing L3out BGP Protocol Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_bgp_protocol_profile.example <Dn>
```
