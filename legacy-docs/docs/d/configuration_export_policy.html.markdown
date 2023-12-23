---
subcategory: "Import/Export"
layout: "aci"
page_title: "ACI: aci_configuration_export_policy"
sidebar_current: "docs-aci-data-source-aci_configuration_export_policy"
description: |-
  Data source for ACI Configuration Export Policy
---

# aci_configuration_export_policy

Data source for ACI Configuration Export Policy

## Example Usage

```hcl
data "aci_configuration_export_policy" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of Object configuration export policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Configuration Export Policy.
- `admin_st` - (Optional) Admin state of the export policy
- `annotation` - (Optional) Annotation for object configuration export policy.
- `format` - (Optional) Export data format.
- `description` - (Optional) Description for object configuration export policy.
- `include_secure_fields` - (Optional) Include secure fields for object configuration export policy.
- `max_snapshot_count` - (Optional) Max snapshot count for object configuration export policy.
- `name_alias` - (Optional) Name alias for object configuration export policy.
- `snapshot` - (Optional) Snapshot for object configuration export policy.
- `target_dn` - (Optional) Target export object.The distinguished name of the object to be exported.
