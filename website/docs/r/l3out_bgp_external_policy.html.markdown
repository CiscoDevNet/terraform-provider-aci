---
layout: "aci"
page_title: "ACI: aci_l3out_bgp_external_policy"
sidebar_current: "docs-aci-resource-l3out_bgp_external_policy"
description: |-
  Manages ACI L3-out BGP External Policy
---

# aci_l3out_bgp_external_policy

Manages ACI L3-out BGP External Policy

## Example Usage

```hcl
resource "aci_l3out_bgp_external_policy" "example" {

  l3_outside_dn  = "${aci_l3_outside.example.id}"
  annotation  = "example"
  name_alias  = "example"

}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of parent l3 outside object.
- `annotation` - (Optional) Annotation for object L3-out BGP External Policy.
- `name_alias` - (Optional) Name alias for object L3-out BGP External Policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3-out BGP External Policy.

## Importing

An existing L3-out BGP External Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_bgp_external_policy.example <Dn>
```
