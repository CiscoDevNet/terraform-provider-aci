---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_action_rule_profile_additional_communities"
sidebar_current: "docs-aci-action-rule-profile-additional-communities"
description: |-
  Data source for ACI Action Rule Profile Set Additional Communities
---

# aci_action_rule_profile_additional_communities #

Data source for ACI Action Rule Profile Set Additional Communities


## API Information ##

* `Class` - rtctrlSetAddComm
* `Distinguished Name` - uni/tn-{tenant_name}/attr-{rule_name}/saddcomm-{community}

## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules -> Rule -> Additional Communities


## Example Usage ##

```hcl
data "aci_action_rule_profile_additional_communities" "example" {
  action_rule_profile_dn  = aci_action_rule_profile.example.id
  community  = "example"
}
```

## Argument Reference ##

* `action_rule_profile_dn` - (Required) Distinguished name of the parent action rule profile object.
* `community` - (Required) The community value of the set action rule profile additional communities object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the additional communities object.
* `annotation` - (Optional) Annotation of the additional communities object.
* `name_alias` - (Optional) Name Alias of the additional communities object.
* `set_criteria` - (Optional) The criteria for setting the (extended) community attribute for a BGP route update.
* `type` - (Optional) The type of the set action rule profile additional communities object.
