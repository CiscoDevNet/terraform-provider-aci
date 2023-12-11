---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_maintenance_policy"
sidebar_current: "docs-aci-data-source-aci_maintenance_policy"
description: |-
  Data source for ACI Maintenance Policy
---

# aci_maintenance_policy

Data source for ACI Maintenance Policy

## Example Usage

```hcl
data "aci_maintenance_policy" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) The name for the maintenance policy. The maintenance policy determines the pre-defined action to take when there is a disruptive change made to the service profile associated with the node or node group.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Maintenance Policy.
- `admin_st` - (Optional) The administrative state of the executable policies. It will trigger an immediate upgrade for nodes if adminst is set to triggered. Once upgrade is done, value is reset back to untriggered.
- `annotation` - (Optional) Annotation for object maintenance policy.
- `description` - (Optional) Description for object maintenance policy.
- `graceful` - (Optional) Whether the system will bring down the nodes gracefully during an upgrade, which reduces traffic lost.
- `ignore_compat` - (Optional) A property that specifies whether compatibility checks should be ignored when applying the firmware policy.
- `internal_label` - (Optional) The firmware label. This is for internal use only.
- `name_alias` - (Optional) Name alias for object maintenance policy.
- `notif_cond` - (Optional) Specifies under what pause condition will admin be notified via email/text as configured. This notification mechanism is independent of events/faults.
- `run_mode` - (Optional) Specifies whether to proceed automatically to next set of nodes once a set of nodes have gone through maintenance successfully.
- `version` - (Optional) The version of the compatibility catalog.
- `version_check_override` - (Optional) The version check override. This is a directive to ignore the version check for the next install. The version check, which occurs during a maintenance window, checks to see if the desired version matches the running version. If the versions do not match, the install is performed. If the versions do match, the install is not performed. The version check override is a one-time override that performs the install whether or not the versions match.
