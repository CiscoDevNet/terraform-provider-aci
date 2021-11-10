---
subcategory: "Import/Export"
layout: "aci"
page_title: "ACI: aci_configuration_import_policy"
sidebar_current: "docs-aci-resource-configuration_import_policy"
description: |-
  Manages ACI Configuration Import Policy
---

# aci_configuration_import_policy

Manages ACI Configuration Import Policy

## Example Usage

```hcl
resource "aci_configuration_import_policy" "example" {
  name  = "import_pol"
  admin_st  = "untriggered"
  annotation  = "example"
  description = "from terraform"
  fail_on_decrypt_errors  = "yes"
  file_name  = "file.tar.gz"
  import_mode  = "best-effort"
  import_type  = "replace"
  name_alias  = "example"
  snapshot  = "no"
}
```

## Argument Reference

- `name` - (Required) The name of the import policy. For ease of reference, include details such as: the full or partial name of the file to be imported, the type/mode of import, and the remote location where the file is stored. The name cannot be changed after the policy has been created.
- `admin_st` - (Optional) The administrative state of the executable policies. A policy can be triggered at any time by setting the admin_st to triggered. The value on APIC will reset back to untriggered once trigger is done. Allowed values: "untriggered", "triggered". Default value: "untriggered".
- `annotation` - (Optional) Specifies the annotation of a policy component.
- `description` - (Optional) Specifies the description of a policy component.
- `fail_on_decrypt_errors` - (Optional) Determines if the APIC should fail the rollback if unable to decrypt secured data. Allowed values: "no", "yes". Default value: "yes".
- `file_name` - (Optional) The name of the file to be imported from the remote location listed below.
- `import_mode` - (Optional) The import mode. The configuration data is imported per shard with each shard holding certain part of the system configuration objects. Allowed values: "atomic", "best-effort". Default value: "atomic".
- `import_type` - (Optional) The import type specifies whether the existing fabric configuration will be merged or replaced with the backup configuration being imported. Allowed values: "merge", "replace". Default value: "merge".
- `name_alias` - (Optional) Name alias for object configuration import policy.
- `snapshot` - (Optional) Snapshot for object configuration import policy. Allowed values: "no", "yes". Default value: "no".

- `relation_config_rs_import_source` - (Optional) Relation to class fileRemotePath. Cardinality - ONE_TO_ONE. Type - String.
- `relation_trig_rs_triggerable` - (Optional) Relation to class trigTriggerable. Cardinality - ONE_TO_ONE. Type - String.
- `relation_config_rs_remote_path` - (Optional) Relation to class fileRemotePath. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Configuration Import Policy.

## Importing

An existing Configuration Import Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_configuration_import_policy.example <Dn>
```
