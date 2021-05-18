---
layout: "aci"
page_title: "ACI: aci_cdp_interface_policy"
sidebar_current: "docs-aci-resource-cdp_interface_policy"
description: |-
  Manages ACI CDP Interface Policy
---

# aci_cdp_interface_policy #

Manages ACI CDP Interface Policy

## Example Usage ##

```hcl
resource "aci_cdp_interface_policy" "example" {
  name        = "example"
  admin_st    = "enabled"
  annotation  = "tag_cdp"
  name_alias  = "alias_cdp"
}
```

## Argument Reference ##

* `name` - (Required) name of Object cdp_interface_policy.
* `admin_st` - (Optional) administrative state.  Allowed values: "enabled", "disabled".  Default value is "enabled".
* `annotation` - (Optional) annotation for object cdp_interface_policy.
* `name_alias` - (Optional) name_alias for object cdp_interface_policy.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the CDP Interface Policy.

## Importing ##

An existing CDP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_cdp_interface_policy.example <Dn>
```
