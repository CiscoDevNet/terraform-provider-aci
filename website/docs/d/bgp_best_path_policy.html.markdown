---
layout: "aci"
page_title: "ACI: aci_bgp_best_path_policy"
sidebar_current: "docs-aci-data-source-bgp_best_path_policy"
description: |-
  Data source for ACI BGP Best Path Policy
---

# aci_bgp_best_path_policy

Data source for ACI BGP Best Path Policy

## Example Usage

```hcl
data "aci_bgp_best_path_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object BGP Best Path Policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the BGP Best Path Policy.
- `annotation` - (Optional) Annotation for object BGP Best Path Policy.
- `ctrl` - (Optional) The control state.
- `description` - Description for the object of the BGP Best Path Policy.
- `name_alias` - (Optional) Name alias for object BGP Best Path Policy.
