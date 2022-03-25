---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_action_rule_profile"
sidebar_current: "docs-aci-data-source-action_rule_profile"
description: |-
  Data source for ACI Action Rule Profile
---

# aci_action_rule_profile #

Data source for ACI Action Rule Profile

## API Information ##

* `Class` - rtctrlAttrP
* `Distinguished Name` - uni/tn-{name}/attr-{name}

## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules

## Example Usage ##

```hcl
data "aci_action_rule_profile" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "Rule-1"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of object Action Rule Profile.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Action Rule Profile.
* `annotation` - (Optional) Annotation of object Action Rule Profile.
* `name_alias` - (Optional) Name Alias of object Action Rule Profile.
* `description` - (Optional) Description of object Action Rule Profile.
