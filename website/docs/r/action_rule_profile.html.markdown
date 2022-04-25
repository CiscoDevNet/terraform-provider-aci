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
  tenant_dn       = aci_tenant.example.id
  description     = "From Terraform"
  name            = "Rule-1"
  annotation      = "orchestrator:terraform"
  name_alias      = "example"
  set_route_tag   = 100
  set_preference  = 100
  set_weight      = 100
  set_metric      = 100
  set_metric_type = "ospf-type1"
  set_next_hop    = "1.1.1.1"
  set_communities = {
    community = "no-advertise"
    criteria  = "replace"
  }
  next_hop_propagation = "yes"
  multipath            = "yes"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Action Rule Profile object.
* `annotation` - (Optional) Annotation of the Action Rule Profile object.
* `description` - (Optional) Description of the Action Rule Profile object.
* `name_alias` - (Optional) Name alias of the Action Rule Profile object.
* `set_route_tag` - (Optional) Set Route Tag of the Action Rule Profile object. Type: Integer.
* `set_preference` - (Optional) Set Preference of the Action Rule Profile object. Type: Integer.
* `set_weight` - (Optional) Set Weight of the Action Rule Profile object. Type: Integer.
* `set_metric` - (Optional) Set Metric of the Action Rule Profile object. Type: Integer.
* `set_metric_type` - (Optional) Set Metric Type of the Action Rule Profile object. Allowed values are `ospf-type1`, `ospf-type2`.
* `set_next_hop` - (Optional) Set Next Hop of the Action Rule Profile object.
* `set_communities` - (Optional and Map of String) Map of the key-value pairs which represents the attributes of Set Communities object. The expected map attributes are ```community``` and ```criteria```.
* `next_hop_propagation` - (Optional) Next Hop Propagation of the Action Rule Profile object.
* `multipath` - (Optional) Multipath of the Action Rule Profile object.

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
