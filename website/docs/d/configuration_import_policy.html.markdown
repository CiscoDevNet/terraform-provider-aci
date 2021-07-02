---
layout: "aci"
page_title: "ACI: aci_configuration_import_policy"
sidebar_current: "docs-aci-data-source-configuration_import_policy"
description: |-
  Data source for ACI Configuration Import Policy
---

# aci_configuration_import_policy

Data source for ACI Configuration Import Policy

## Example Usage

```hcl
data "aci_configuration_import_policy" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) The name of the import policy. For ease of reference, include details such as: the full or partial name of the file to be imported, the type/mode of import, and the remote location where the file is stored. The name cannot be changed after the policy has been created.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Configuration Import Policy.
- `admin_st` - (Optional) The administrative state of the executable policies. A policy can be triggered at any time by setting the admin_st to triggered. The value on APIC will reset back to untriggered once trigger is done.
- `annotation` - (Optional) Specifies the annotation of a policy component.
- `description` - (Optional) Specifies the description of a policy component.
- `fail_on_decrypt_errors` - (Optional) Determines if the APIC should fail the rollback if unable to decrypt secured data.
- `file_name` - (Optional) The name of the file to be imported from the remote location listed below.
- `import_mode` - (Optional) The import mode. The configuration data is imported per shard with each shard holding certain part of the system configuration objects.
- `import_type` - (Optional) The import type specifies whether the existing fabric configuration will be merged or replaced with the backup configuration being imported.
- `name_alias` - (Optional) Name alias for object configuration import policy.
- `snapshot` - (Optional) Snapshot for object configuration import policy.
