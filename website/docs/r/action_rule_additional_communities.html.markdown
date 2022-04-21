---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_action_rule_additional_communities"
sidebar_current: "docs-aci-action-rule-additional-communities"
description: |-
  Manages ACI Action Rule Profile Set Additional Communities
---

# aci_action_rule_additional_communities #

  Manages ACI Action Rule Profile Set Additional Communities

## API Information ##

* `Class` - rtctrlSetAddComm
* `Distinguished Name` - uni/tn-{tenant_name}/attr-{rule_name}/saddcomm-{community}


## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules -> Rule -> Additional Communities


## Example Usage ##

```hcl
resource "aci_action_rule_additional_communities" "example" {
  action_rule_profile_dn  = aci_action_rule_profile.example.id
  community  = "no-advertise"
  annotation = "orchestrator:terraform"
  set_criteria = "append"
}
```

## Argument Reference ##

* `action_rule_profile_dn` - (Required) Distinguished name of the parent action rule profile object.
* `community` - (Required) The community value of the set action rule profile additional communities object.
* `description` - (Optional) The description of the set action rule profile additional communities object.
* `annotation` - (Optional) Annotation of the action rule profile additional communities object.
* `name_alias` - (Optional) Name Alias of the additional communities object.
* `set_criteria` - (Optional) The criteria for setting the (extended) community attribute for a BGP route update. Allowed values are "append", "none", "replace", and default value is "append". Type: String.


## Importing ##

An existing Action Rule Profile Additional Communities can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_action_rule_additional_communities.example <Dn>
```