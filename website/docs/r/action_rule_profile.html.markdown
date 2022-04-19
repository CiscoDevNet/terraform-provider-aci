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
* `Distinguished Name` - uni/tn-{tenant_name}/attr-{rule_name}

## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules

## Example Usage ##

```hcl
resource "aci_action_rule_profile" "example" {
  tenant_dn     = aci_tenant.example.id
  name          = "Rule-1"
  annotation    = "orchestrator:terraform"
  set_route_tag = 100
  set_preference = 100
  set_weight = 100
  set_metric = 100
  set_metric_type = "ospf-type1"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Action Rule Profile object.
* `annotation` - (Optional) Annotation of the Action Rule Profile object.
* `description` - (Optional) Description of the Action Rule Profile object.
* `name_alias` - (Optional) Name alias of the Action Rule Profile object.
* `set_route_tag` - (Optional) Set Route Tag of the Action Rule Profile object.
* `set_preference` - (Optional) Set Preference of the Action Rule Profile object.
* `set_weight` - (Optional) Set Weight of the Action Rule Profile object.
* `set_metric` - (Optional) Set Metric of the Action Rule Profile object.
* `set_metric_type` - (Optional) Set Metric Type of the Action Rule Profile object. Allowed values are "ospf-type1", "ospf-type2".

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
