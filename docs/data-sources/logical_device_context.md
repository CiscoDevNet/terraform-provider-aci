---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_logical_device_context"
sidebar_current: "docs-aci-data-source-logical_device_context"
description: |-
  Data source for ACI Logical Device Context
---

# aci_logical_device_context

Data source for ACI Logical Device Context

## Example Usage

```hcl

data "aci_logical_device_context" "check" {
  tenant_dn         = aci_tenant.tenentcheck.id
  ctrct_name_or_lbl = "example"
  graph_name_or_lbl = "example"
  node_name_or_lbl  = "example"
}

```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `ctrct_name_or_lbl` - (Required) Ctrct name or label of Object logical device context.
- `graph_name_or_lbl` - (Required) Graph name or label of Object logical device context.
- `node_name_or_lbl` - (Required) Node name or label of Object logical device context.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Device Context.
- `annotation` - (Optional) Annotation for object logical device context.
- `description` - (Optional) Description for object logical device context.
- `context` - (Optional) Context for object logical device context.
- `name_alias` - (Optional) Name alias for object logical device context.
