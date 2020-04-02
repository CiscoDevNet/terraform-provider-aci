---
layout: "aci"
page_title: "ACI: aci_action_rule_profile"
sidebar_current: "docs-aci-data-source-action_rule_profile"
description: |-
  Data source for ACI Action Rule Profile
---

# aci_action_rule_profile #
Data source for ACI Action Rule Profile

## Example Usage ##

```hcl
data "aci_action_rule_profile" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object action_rule_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Action Rule Profile.
* `annotation` - (Optional) annotation for object action_rule_profile.
* `name_alias` - (Optional) name_alias for object action_rule_profile.
