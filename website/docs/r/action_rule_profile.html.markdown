---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_action_rule_profile"
sidebar_current: "docs-aci-resource-action_rule_profile"
description: |-
  Manages ACI Action Rule Profile
---

# aci_action_rule_profile #

Manages ACI Action Rule Profile

## API Information ##

* `Class` - rtctrlAttrP
* `Distinguished Name` - uni/tn-{name}/attr-{name}

## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules

## Example Usage ##

```hcl
resource "aci_action_rule_profile" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "Rule-1"
  annotation = "orchestrator:terraform"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the object Action Rule Profile.
* `annotation` - (Optional) Annotation of the object Action Rule Profile.
* `description` - (Optional) Description of the object Action Rule Profile.
* `name_alias` - (Optional) Name alias of the object Action Rule Profile.

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
