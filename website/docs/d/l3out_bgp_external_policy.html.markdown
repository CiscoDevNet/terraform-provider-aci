---
layout: "aci"
page_title: "ACI: aci_l3out_bgp_external_policy"
sidebar_current: "docs-aci-data-source-l3out_bgp_external_policy"
description: |-
  Data source for ACI L3-out BGP External Policy
---

# aci_l3out_bgp_external_policy

Data source for ACI L3-out BGP External Policy

## Example Usage

```hcl
data "aci_l3out_bgp_external_policy" "example" {
  l3_outside_dn  = "${aci_l3_outside.example.id}"
}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of parent l3 outside object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3-out BGP External Policy.
- `annotation` - (Optional) Annotation for object l3out BGP external policy.
- `name_alias` - (Optional) Name alias for object l3out BGP external policy.
