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

- `name` - (Required) Name of Object Trigger Scheduler

## Attribute Reference

- `id` - Attribute id set to the Dn of the Trigger Scheduler.
- `annotation` - (Optional) Annotation for object Trigger Scheduler
- `name_alias` - (Optional) Name alias for object Trigger Scheduler
- `description` - (Optional) Description for object Trigger Scheduler
