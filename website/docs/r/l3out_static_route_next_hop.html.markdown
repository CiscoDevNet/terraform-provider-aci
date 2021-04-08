---
layout: "aci"
page_title: "ACI: aci_l3out_static_route_next_hop"
sidebar_current: "docs-aci-resource-l3out_static_route_next_hop"
description: |-
  Manages ACI L3out Static Route Next Hop
---

# aci_l3out_static_route_next_hop

Manages ACI L3out Static Route Next Hop

## Example Usage

```hcl
resource "aci_l3out_static_route_next_hop" "example" {

  static_route_dn  = "${aci_l3out_static_route.example.id}"
  nh_addr  = "10.0.0.1"
  annotation  = "example"
  name_alias  = "example"
  pref = "unspecified"
  nexthop_profile_type = "prefix"

}
```

## Argument Reference

- `static_route_dn` - (Required) Distinguished name of parent static route object.
- `nh_addr` - (Required) The nexthop IP address for the static route to the outside network.
- `annotation` - (Optional) Annotation for object l3out static route next hop.
- `description` - (Optional) Description for object l3out static route next hop.
- `name_alias` - (Optional) Name alias for object l3out static route next hop.
- `pref` - (Optional) Administrative preference value for this route.  
  Allowed values: "unspecified". Default value: "unspecified".
- `nexthop_profile_type` - (Optional) Component type.  
  Allowed values: "none", "prefix". Default value: "prefix".

- `relation_ip_rs_nexthop_route_track` - (Optional) Relation to class fvTrackList. Cardinality - N_TO_ONE. Type - String.
- `relation_ip_rs_nh_track_member` - (Optional) Relation to class fvTrackMember. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out Static Route Next Hop.

## Importing

An existing L3out Static Route Next Hop can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_static_route_next_hop.example <Dn>
```
