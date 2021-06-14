---
layout: "aci"
page_title: "ACI: aci_bgp_route_summarization"
sidebar_current: "docs-aci-data-source-bgp_route_summarization"
description: |-
  Data source for ACI BGP Route Summarization
---

# aci_bgp_route_summarization

Data source for ACI BGP Route Summarization

## Example Usage

```hcl
data "aci_bgp_route_summarization" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object BGP route summarization.

## Attribute Reference

- `id` - Attribute id set to the Dn of the BGP route summarization.
- `annotation` - (Optional) Annotation for object BGP route summarization.
- `attrmap` - (Optional) Summary attribute map.
- `ctrl` - (Optional) The control state.
- `name_alias` - (Optional) Name alias for object BGP route summarization.
