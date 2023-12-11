---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_logical_interface_context"
sidebar_current: "docs-aci-data-source-aci_logical_interface_context"
description: |-
  Data source for ACI Logical Interface Context
---

# aci_logical_interface_context

Data source for ACI Logical Interface Context

## Example Usage

```hcl
data "aci_logical_interface_context" "example" {
  logical_device_context_dn  = aci_logical_device_context.example.id
  conn_name_or_lbl = "example"
}
```

## Argument Reference

- `logical_device_context_dn` - (Required) Distinguished name of parent Logical Device Context object.
- `conn_name_or_lbl` - (Required) The connector name or label for the Logical Interface Context.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Interface Context.
- `annotation` - (Optional) Annotation for object Logical Interface Context.
- `description` - (Optional) Description for object Logical interface context.
- `l3_dest` - (Optional) L3 dest for object Logical Interface Context.
- `name_alias` - (Optional) Name alias for object Logical Interface Context.
- `permit_log` - (Optional) Permit log for object Logical Interface Context.
