---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_miscabling_protocol_interface_policy"
sidebar_current: "docs-aci-data-source-miscabling_protocol_interface_policy"
description: |-
  Data source for ACI Mis-cabling Protocol Interface Policy
---

# aci_miscabling_protocol_interface_policy #
Data source for ACI Mis-cabling Protocol Interface Policy

## Example Usage ##

```hcl
data "aci_miscabling_protocol_interface_policy" "dev_miscable_pol" {
  name  = "foo_miscable_pol"
}
```
## Argument Reference ##
* `name` - (Required) name of Object miscabling_protocol_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Mis-cabling Protocol Interface Policy.
* `admin_st` - (Optional) Administrative state of the object or policy.
* `description` - (Optional) Description for object miscabling protocol interface policy.
* `annotation` - (Optional) Annotation for object miscabling protocol interface policy.
* `name_alias` - (Optional) Name alias for object miscabling protocol interface policy.
