---
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
data "aci_miscabling_protocol_interface_policy" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object miscabling_protocol_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Mis-cabling Protocol Interface Policy.
* `admin_st` - (Optional) administrative state of the object or policy.
* `annotation` - (Optional) annotation for object miscabling_protocol_interface_policy.
* `name_alias` - (Optional) name_alias for object miscabling_protocol_interface_policy.
