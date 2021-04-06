---
layout: "aci"
page_title: "ACI: aci_l3out_route_tag_policy"
sidebar_current: "docs-aci-data-source-l3out_route_tag_policy"
description: |-
  Data source for ACI L3out Route Tag Policy
---

# aci_l3out_route_tag_policy

Data source for ACI L3out Route Tag Policy

## Example Usage

```hcl
data "aci_l3out_route_tag_policy" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object l3out route tag policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out Route Tag Policy.
- `annotation` - (Optional) Annotation for object L3out route tag policy.
- `name_alias` - (Optional) Name alias for object L3out route tag policy.
- `tag` - (Optional) Tagged number.
