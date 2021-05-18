---
layout: "aci"
page_title: "ACI: aci_trigger_scheduler"
sidebar_current: "docs-aci-data-source-trigger_scheduler"
description: |-
  Data source for ACI Trigger Scheduler
---

# aci_trigger_scheduler

Data source for ACI Trigger Scheduler

## Example Usage

```hcl
data "aci_trigger_scheduler" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) name of Object trigger_scheduler.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Trigger Scheduler.
- `annotation` - (Optional) annotation for object trigger_scheduler.
- `name_alias` - (Optional) name_alias for object trigger_scheduler.
- `description` - (Optional) Description for object trigger_scheduler.
