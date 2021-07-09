---
layout: "aci"
page_title: "ACI: aci_action_rule_profile"
sidebar_current: "docs-aci-resource-action_rule_profile"
description: |-
  Manages ACI Action Rule Profile
---

# aci_action_rule_profile #

Manages ACI Action Rule Profile

## Example Usage ##

```hcl
resource "aci_action_rule_profile" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name       = "example"
  annotation = "example"
  name_alias = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object action_rule_profile.
* `annotation` - (Optional) annotation for object action_rule_profile.
* `name_alias` - (Optional) name_alias for object action_rule_profile.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Action Rule Profile.

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
