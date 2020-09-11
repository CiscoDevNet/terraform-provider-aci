---
layout: "aci"
page_title: "ACI: aci_maintenance_policy"
sidebar_current: "docs-aci-resource-maintenance_policy"
description: |-
  Manages ACI Maintenance Policy
---

# aci_maintenance_policy #
Manages ACI Maintenance Policy

## Example Usage ##

```hcl
resource "aci_maintenance_policy" "example" {


  name  = "example"
  admin_st  = "example"
  annotation  = "example"
  graceful  = "example"
  ignore_compat  = "example"
  internal_label  = "example"
  name_alias  = "example"
  notif_cond  = "example"
  run_mode  = "example"
  version  = "example"
  version_check_override  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object maintenance_policy.
* `admin_st` - (Optional) maintenance policy admin state.
Allowed values: "untriggered", "triggered"
* `annotation` - (Optional) annotation for object maintenance_policy.
* `graceful` - (Optional) graceful for object maintenance_policy.
Allowed values: "yes", "no"
* `ignore_compat` - (Optional) whether compatibility check required.
Allowed values: "yes", "no"
* `internal_label` - (Optional) firmware label
* `name_alias` - (Optional) name_alias for object maintenance_policy.
* `notif_cond` - (Optional) when to send notifications to the admin.
Allowed values: "notifyOnlyOnFailures","notifyAlwaysBetweenSets", "notifyNever"
* `run_mode` - (Optional) maintenance policy run mode.
Allowed values: "pauseOnlyOnFailures","pauseAlwaysBetweenSets", "pauseNever"
* `version` - (Optional) compatibility catalog version
* `version_check_override` - (Optional) version check override.
Allowed values: "trigger-immediate", "trigger", "triggered","untriggered"

* `relation_maint_rs_pol_scheduler` - (Optional) Relation to class trigSchedP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_maint_rs_pol_notif` - (Optional) Relation to class maintUserNotif. Cardinality - N_TO_ONE. Type - String.
                
* `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Maintenance Policy.

## Importing ##

An existing Maintenance Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_maintenance_policy.example <Dn>
```