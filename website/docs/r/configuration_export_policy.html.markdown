---
layout: "aci"
page_title: "ACI: aci_configuration_export_policy"
sidebar_current: "docs-aci-resource-configuration_export_policy"
description: |-
  Manages ACI Configuration Export Policy
---

# aci_configuration_export_policy #
Manages ACI Configuration Export Policy

## Example Usage ##

```hcl
resource "aci_configuration_export_policy" "example" {


  name  = "example"
  admin_st  = "example"
  annotation  = "example"
  format  = "example"
  include_secure_fields  = "example"
  max_snapshot_count  = "example"
  name_alias  = "example"
  snapshot  = "example"
  target_dn  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object configuration_export_policy.
* `admin_st` - (Optional) admin state of the export policy
Allowed values: "untriggered", "triggered"
* `annotation` - (Optional) annotation for object configuration_export_policy.
* `format` - (Optional) export data format
Allowed values: "xml", "json"
* `include_secure_fields` - (Optional) include_secure_fields for object configuration_export_policy.
Allowed values: "no", "yes"
* `max_snapshot_count` - (Optional) max_snapshot_count for object configuration_export_policy.
* `name_alias` - (Optional) name_alias for object configuration_export_policy.
* `snapshot` - (Optional) snapshot for object configuration_export_policy.
Allowed values: "no", "yes"
* `target_dn` - (Optional) target export object

* `relation_config_rs_export_destination` - (Optional) Relation to class fileRemotePath. Cardinality - ONE_TO_ONE. Type - String.
                
* `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.
                
* `relation_config_rs_remote_path` - (Optional) Relation to class fileRemotePath. Cardinality - N_TO_ONE. Type - String.
                
* `relation_config_rs_export_scheduler` - (Optional) Relation to class trigSchedP. Cardinality - ONE_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Configuration Export Policy.

## Importing ##

An existing Configuration Export Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_configuration_export_policy.example <Dn>
```