---
layout: "aci"
page_title: "ACI: aci_configuration_export_policy"
sidebar_current: "docs-aci-resource-configuration_export_policy"
description: |-
  Manages ACI Configuration Export Policy
---

# aci_configuration_export_policy

Manages ACI Configuration Export Policy

## Example Usage

```hcl
resource "aci_configuration_export_policy" "example" {
  name                  = "example"
  description           = "from terraform"
  admin_st              = "untriggered"
  annotation            = "example"
  format                = "json"
  include_secure_fields = "yes"
  max_snapshot_count    = "10"
  name_alias            = "example"
  snapshot              = "yes"
  target_dn             = "uni/tn-test"
}
```

## Argument Reference

- `name` - (Required) Name of Object configuration export policy.
- `admin_st` - (Optional) Admin state of the export policy. A policy can be triggered at any time by setting the admin_st to triggered. The value on APIC will reset back to untriggered once trigger is done. 
  Allowed values: "untriggered", "triggered". Default value is "untriggered".
- `annotation` - (Optional) Annotation for object configuration export policy.
- `description` - (Optional) Description for object configuration export policy.
- `format` - (Optional) Export data format.
  Allowed values: "xml", "json". Default value is "json".
- `include_secure_fields` - (Optional) Include_secure_fields for object configuration export policy.
  Allowed values: "no", "yes".Default value is "yes".
- `max_snapshot_count` - (Optional) Max snapshot count for object configuration export policy.
  Allowed Values are betwwen 0 to 10. Default value is "global-limit" (0 is consider as a global limit).
- `name_alias` - (Optional) Name alias for object configuration export policy.
- `snapshot` - (Optional) Snapshot for object configuration export policy.
  Allowed values: "no", "yes"Default value is "no".
- `target_dn` - (Optional) Target export object. The distinguished name of the object to be exported.

- `relation_config_rs_export_destination` - (Optional) Relation to class fileRemotePath. Cardinality - ONE_TO_ONE. Type - String.
- `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.
- `relation_config_rs_remote_path` - (Optional) Relation to class fileRemotePath. Cardinality - N_TO_ONE. Type - String.
- `relation_config_rs_export_scheduler` - (Optional) Relation to class trigSchedP. Cardinality - ONE_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Configuration Export Policy.

## Importing

An existing Configuration Export Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_configuration_export_policy.example <Dn>
```
