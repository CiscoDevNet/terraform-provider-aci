---
layout: "aci"
page_title: "ACI: aci_miscabling_protocol_interface_policy"
sidebar_current: "docs-aci-resource-miscabling_protocol_interface_policy"
description: |-
  Manages ACI Mis-cabling Protocol Interface Policy
---

# aci_miscabling_protocol_interface_policy #
Manages ACI Mis-cabling Protocol Interface Policy

## Example Usage ##

```hcl
	resource "aci_miscabling_protocol_interface_policy" "foomiscabling_protocol_interface_policy" {
		description = "%s"
		name        = "demo_mcpol"
		admin_st    = "%s"
		annotation  = "tag_mcpol"
		name_alias  = "alias_mcpol"
	}
```
## Argument Reference ##
* `name` - (Required) name of Object miscabling_protocol_interface_policy.
* `admin_st` - (Optional) administrative state of the object or policy. Allowed values are "enabled" and "disabled". Default is "enabled".
* `annotation` - (Optional) annotation for object miscabling_protocol_interface_policy.
* `name_alias` - (Optional) name_alias for object miscabling_protocol_interface_policy.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Mis-cabling Protocol Interface Policy.

## Importing ##

An existing Mis-cabling Protocol Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_miscabling_protocol_interface_policy.example <Dn>
```