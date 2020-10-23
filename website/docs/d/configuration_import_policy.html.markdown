---
layout: "aci"
page_title: "ACI: aci_configuration_import_policy"
sidebar_current: "docs-aci-data-source-configuration_import_policy"
description: |-
  Data source for ACI Configuration Import Policy
---

# aci_configuration_import_policy #
Data source for ACI Configuration Import Policy

## Example Usage ##

```hcl
data "aci_configuration_import_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object configuration_import_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Configuration Import Policy.
* `admin_st` - (Optional) admin state of the import
* `annotation` - (Optional) annotation for object configuration_import_policy.
* `fail_on_decrypt_errors` - (Optional) fail_on_decrypt_errors for object configuration_import_policy.
* `file_name` - (Optional) import file name
* `import_mode` - (Optional) data import mode
* `import_type` - (Optional) data import type
* `name_alias` - (Optional) name_alias for object configuration_import_policy.
* `snapshot` - (Optional) snapshot for object configuration_import_policy.
