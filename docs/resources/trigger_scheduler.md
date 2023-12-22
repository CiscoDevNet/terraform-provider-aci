---
subcategory: "Scheduler"
layout: "aci"
page_title: "ACI: aci_trigger_scheduler"
sidebar_current: "docs-aci-resource-aci_trigger_scheduler"
description: |-
  Manages ACI Trigger Scheduler
---

# aci_trigger_scheduler

Manages ACI Trigger Scheduler

## Example Usage

```hcl
resource "aci_trigger_scheduler" "example" {
  name  = "example"
  annotation  = "example"
  description = "from terraform"
  name_alias  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of Object Trigger Scheduler
- `annotation` - (Optional) Annotation for object Trigger Scheduler
- `name_alias` - (Optional) Name alias for object Trigger Scheduler
- `description` - (Optional) Description for object Trigger Scheduler

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Trigger Scheduler.

## Importing

An existing Trigger Scheduler can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_trigger_scheduler.example <Dn>
```
