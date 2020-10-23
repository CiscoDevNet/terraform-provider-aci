---
layout: "aci"
page_title: "ACI: aci_configuration_import_policy"
sidebar_current: "docs-aci-resource-configuration_import_policy"
description: |-
  Manages ACI Configuration Import Policy
---

# aci_configuration_import_policy #
Manages ACI Configuration Import Policy

## Example Usage ##

```hcl
resource "aci_configuration_import_policy" "example" {


  name  = "example"
  admin_st  = "example"
  annotation  = "example"
  fail_on_decrypt_errors  = "example"
  file_name  = "example"
  import_mode  = "example"
  import_type  = "example"
  name_alias  = "example"
  snapshot  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object configuration_import_policy.
* `admin_st` - (Optional) admin state of the import
Allowed values: "untriggered", "triggered"
* `annotation` - (Optional) annotation for object configuration_import_policy.
* `fail_on_decrypt_errors` - (Optional) fail_on_decrypt_errors for object configuration_import_policy.
Allowed values: "no", "yes"
* `file_name` - (Optional) import file name
* `import_mode` - (Optional) data import mode.
Allowed values: "atomic", "best-effort"
* `import_type` - (Optional) data import type.
Allowed values: "merge", "replace"
* `name_alias` - (Optional) name_alias for object configuration_import_policy.
* `snapshot` - (Optional) snapshot for object configuration_import_policy.
Allowed values: "no", "yes"

* `relation_config_rs_import_source` - (Optional) Relation to class fileRemotePath. Cardinality - ONE_TO_ONE. Type - String.
                
* `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.
                
* `relation_config_rs_remote_path` - (Optional) Relation to class fileRemotePath. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Configuration Import Policy.

## Importing ##

An existing Configuration Import Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_configuration_import_policy.example <Dn>
```