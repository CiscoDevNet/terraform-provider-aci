---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_route_control_context"
sidebar_current: "docs-aci-resource-aci_route_control_context"
description: |-
  Manages ACI Route Control Context
---

# aci_route_control_context #

Manages ACI Route Control Context

## API Information ##

* `Class` - rtctrlCtxP
* `Distinguished Name` - uni/tn-{name}/prof-{name}/ctx-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Route Maps for Route Control

## Example Usage ##

```hcl
resource "aci_route_control_context" "control" {
  route_control_profile_dn  = aci_route_control_profile.bgp.id
  name  = "control"
  action = "permit"
  annotation = "orchestrator:terraform"
  order = "0"
  set_rule = aci_action_rule_profile.set_rule1.id
  relation_rtctrl_rs_ctx_p_to_subj_p = [aci_match_rule.rule.id]
}
```

## Argument Reference ##

* `route_control_profile_dn` - (Required) Distinguished name of parent Route Control Profile object.
* `name` - (Required) Name of object Route Control Context.
* `annotation` - (Optional) Annotation of object Route Control Context.
* `action` - (Optional) Action. The action required when the condition is met. Allowed values are "deny", "permit", and default value is "permit". Type: String.
* `order` - (Optional) Local Order.The order of the policy context. Allowed range is 0-9 and default value is "0".
* `set_rule` - (Optional) Represents the relation to an Attribute Profile (class rtctrlAttrP). Type: String.
* `relation_rtctrl_rs_ctx_p_to_subj_p` - (Optional) Represents the relation to a Subject Profile (class rtctrlSubjP). Type: List.


## Importing ##

An existing Route Control Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_route_control_context.example <Dn>
```