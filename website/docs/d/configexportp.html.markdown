---
layout: "aci"
page_title: "ACI: aci_configuration_export_policy"
sidebar_current: "docs-aci-data-source-configuration_export_policy"
description: |-
  Data source for ACI Configuration Export Policy
---

# aci_configuration_export_policy #
Data source for ACI Configuration Export Policy

## Example Usage ##

```hcl
data "aci_configuration_export_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object configuration_export_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Configuration Export Policy.
* `admin_st` - (Optional) admin state of the export policy
* `annotation` - (Optional) annotation for object configuration_export_policy.
* `format` - (Optional) export data format
* `include_secure_fields` - (Optional) include_secure_fields for object configuration_export_policy.
* `max_snapshot_count` - (Optional) max_snapshot_count for object configuration_export_policy.
* `name_alias` - (Optional) name_alias for object configuration_export_policy.
* `snapshot` - (Optional) snapshot for object configuration_export_policy.
* `target_dn` - (Optional) target export object
