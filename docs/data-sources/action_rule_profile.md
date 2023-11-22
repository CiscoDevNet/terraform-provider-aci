---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_action_rule_profile"
sidebar_current: "docs-aci-data-source-action_rule_profile"
description: |-
  Data source for the ACI Action Rule Profile
---

# aci_action_rule_profile #

Data source for the ACI Action Rule Profile

## API Information ##

* `Class` - rtctrlAttrP
* `Distinguished Name` - uni/tn-{tenant_name}/attr-{rule_name}

## GUI Information ##

* `Location` - Tenant > Policies > Protocols > Set Rules

## Example Usage ##

```hcl
data "aci_action_rule_profile" "example" {
  tenant_dn     = aci_tenant.example.id
  name          = "Rule-1"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Action Rule Profile object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Action Rule Profile object.
* `annotation` - (Optional) Annotation of the Action Rule Profile object.
* `name_alias` - (Optional) Name Alias of the Action Rule Profile object.
* `description` - (Optional) Description of the Action Rule Profile object.
* `set_route_tag` - (Optional) Set Route Tag of the Action Rule Profile object.
* `set_preference` - (Optional) Set Preference of the Action Rule Profile object.
* `set_weight` - (Optional) Set Weight of the Action Rule Profile object.
* `set_metric` - (Optional) Set Metric of the Action Rule Profile object.
* `set_metric_type` - (Optional) Set Metric Type of the Action Rule Profile object.
* `set_next_hop` - (Optional) Set Next Hop of the Action Rule Profile object.
* `set_communities` - (Optional) A block representing the attributes of Set Communities object. Type: Block.
  * `criteria` - (Optional) Criteria of the Set Communities object.
  * `community` - (Optional) Community of the Set Communities object.
* `next_hop_propagation` - (Optional) Next Hop Propagation of the Action Rule Profile object.
* `multipath` - (Optional) Multipath of the Action Rule Profile object.
* `set_as_path_prepend_last_as` - (Optional) Number of ASN prepended to AS Path of the Action Rule Profile object.
* `set_as_path_prepend_as` - (Optional) A block representing ASNs configured as Set As Path - Prepend AS of the Action Rule Profile object. Type: Block.
  * `asn` - ASN prepended to Set AS Path.
  * `order` - Order in which the ASN is prepended to Set AS Path.
* `set_dampening` - (Optional) A block representing the attributes of Set Dampening object. Type: Block.
  * `half_life` - Half Life of the Set Dampening object.
  * `reuse` - Reuse Limit of the Set Dampening object.
  * `suppress` - Suppress Limit of the Set Dampening object.
  * `max_suppress_time` - Max Suppress Time of the Set Dampening object.
