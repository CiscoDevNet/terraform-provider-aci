---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_cdp_interface_policy"
sidebar_current: "docs-aci-resource-aci_cdp_interface_policy"
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
  description = "From Terraform"
}
```

## Argument Reference ##

* `name` - (Required) Name of Object cdp interface policy.
* `admin_st` - (Optional) Administrative state.  Allowed values: "enabled", "disabled".  Default value is "enabled".
* `annotation` - (Optional) Annotation for object cdp interface policy.
* `name_alias` - (Optional) Name alias for object cdp interface policy.
* `description` - (Optional) Description for object cdp interface policy.
## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the CDP Interface Policy.

## Importing ##

An existing CDP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_cdp_interface_policy.example <Dn>
```
