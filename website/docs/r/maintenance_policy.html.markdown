---
layout: "aci"
page_title: "ACI: aci_maintenance_policy"
sidebar_current: "docs-aci-resource-maintenance_policy"
description: |-
  Manages ACI Maintenance Policy
---

# aci_maintenance_policy

Manages ACI Maintenance Policy

## Example Usage

```hcl
resource "aci_maintenance_policy" "example" {
  name           = "mnt_policy"
  admin_st       = "triggered"
  description    = "from terraform"
  annotation     = "example"
  graceful       = "yes"
  ignore_compat  = "yes"
  internal_label = "example"
  name_alias     = "example"
  notif_cond     = "notifyOnlyOnFailures"
  run_mode       = "pauseOnlyOnFailures"
  version        = "n9000-15.0(1k)"
  version_check_override = "trigger-immediate"
}
```

## Argument Reference

- `name` - (Required) The name for the maintenance policy. The maintenance policy determines the pre-defined action to take when there is a disruptive change made to the service profile associated with the node or node group.
- `admin_st` - (Optional) The administrative state of the executable policies. It will trigger an immediate upgrade for nodes if adminst is set to triggered. Once upgrade is done, value is reset back to untriggered.
  Allowed values: "untriggered", "triggered". Default value is "untriggered"
- `annotation` - (Optional) Annotation for object maintenance policy.
- `description` - (Optional) Description for object maintenance policy.
- `graceful` - (Optional) Whether the system will bring down the nodes gracefully during an upgrade, which reduces traffic lost. Allowed values: "yes", "no". Default value is "no".
- `ignore_compat` - (Optional) A property that specifies whether compatibility checks should be ignored when applying the firmware policy. Allowed values: "yes", "no". Default value is "no".
- `internal_label` - (Optional) The firmware label. This is for internal use only.
- `name_alias` - (Optional) Name alias for object maintenance policy.
- `notif_cond` - (Optional) Specifies under what pause condition will admin be notified via email/text as configured. This notification mechanism is independent of events/faults. Allowed values: "notifyOnlyOnFailures", "notifyAlwaysBetweenSets", "notifyNever". Default value is "notifyOnlyOnFailures".
- `run_mode` - (Optional) Specifies whether to proceed automatically to next set of nodes once a set of nodes have gone through maintenance successfully. Allowed values: "pauseOnlyOnFailures","pauseAlwaysBetweenSets", "pauseNever". Default value is "pauseOnlyOnFailures".
- `version` - (Optional) The version of the compatibility catalog.
- `version_check_override` - (Optional) The version check override. This is a directive to ignore the version check for the next install. The version check, which occurs during a maintenance window, checks to see if the desired version matches the running version. If the versions do not match, the install is performed. If the versions do match, the install is not performed. The version check override is a one-time override that performs the install whether or not the versions match. Allowed values: "trigger-immediate", "trigger", "triggered","untriggered". Default value is "untriggered".

- `relation_maint_rs_pol_scheduler` - (Optional) Relation to class trigSchedP. Cardinality - N_TO_ONE. Type - String.
- `relation_maint_rs_pol_notif` - (Optional) Relation to class maintUserNotif. Cardinality - N_TO_ONE. Type - String.
- `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Maintenance Policy.

## Importing

An existing Maintenance Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_maintenance_policy.example <Dn>
```
