---
layout: "aci"
page_title: "ACI: aci_logical_interface_context"
sidebar_current: "docs-aci-data-source-logical_interface_context"
description: |-
  Data source for ACI Logical Interface Context
---

# aci_logical_interface_context

Data source for ACI Logical Interface Context

## Example Usage

```hcl
data "aci_logical_interface_context" "example" {

  logical_device_context_dn  = "${aci_logical_device_context.example.id}"
  conn_name_or_lbl = "example"
}
```

## Argument Reference

- `logical_device_context_dn` - (Required) Distinguished name of parent Logical Device Context object.
- `conn_name_or_lbl` - (Required) The connector name or label for the logical interface context.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Interface Context.
- `annotation` - (Optional) Annotation for object logical_interface_context.
- `l3_dest` - (Optional) l3_dest for object logical_interface_context.
- `name_alias` - (Optional) Name_alias for object logical_interface_context.
- `permit_log` - (Optional) permit_log for object logical_interface_context.
