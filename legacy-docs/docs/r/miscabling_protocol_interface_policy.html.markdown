---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_miscabling_protocol_interface_policy"
sidebar_current: "docs-aci-resource-aci_miscabling_protocol_interface_policy"
description: |-
  Manages ACI Mis-cabling Protocol Interface Policy
---

# aci_miscabling_protocol_interface_policy #

Manages ACI Mis-cabling Protocol Interface Policy

## Example Usage ##

```hcl
resource "aci_miscabling_protocol_interface_policy" "example" {
  description = "example description"
  name        = "demo_mcpol"
  admin_st    = "enabled"
  annotation  = "tag_mcpol"
  name_alias  = "alias_mcpol"
}
```

## Argument Reference ##

* `name` - (Required) Name of Object miscabling protocol interface policy.
* `admin_st` - (Optional) Administrative state of the object or policy. Allowed values are "enabled" and "disabled". Default is "enabled".
* `description` - (Optional) Description for object miscabling protocol interface policy.
* `annotation` - (Optional) Annotation for object miscabling protocol interface policy.
* `name_alias` - (Optional) Name alias for object miscabling protocol interface policy.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Mis-cabling Protocol Interface Policy.

## Importing ##

An existing Mis-cabling Protocol Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_miscabling_protocol_interface_policy.example <Dn>
```
