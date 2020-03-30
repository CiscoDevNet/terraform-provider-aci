---
layout: "aci"
page_title: "ACI: aci_maintenance_policy"
sidebar_current: "docs-aci-data-source-maintenance_policy"
description: |-
  Data source for ACI Maintenance Policy
---

# aci_maintenance_policy #
Data source for ACI Maintenance Policy

## Example Usage ##

```hcl
data "aci_maintenance_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object maintenance_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Maintenance Policy.
* `admin_st` - (Optional) maintenance policy admin state
* `annotation` - (Optional) annotation for object maintenance_policy.
* `graceful` - (Optional) graceful for object maintenance_policy.
* `ignore_compat` - (Optional) whether compatibility check required
* `internal_label` - (Optional) firmware label
* `name_alias` - (Optional) name_alias for object maintenance_policy.
* `notif_cond` - (Optional) when to send notifications to the admin
* `run_mode` - (Optional) maintenance policy run mode
* `version` - (Optional) compatibility catalog version
* `version_check_override` - (Optional) version check override
