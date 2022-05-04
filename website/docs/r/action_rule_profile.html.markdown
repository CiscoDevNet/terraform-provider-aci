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
  name            = "example"
  annotation      = "example"
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
  next_hop_propagation    = "yes"
  multipath               = "yes"
  saspath_prepend_last_as = 10
  saspath_prepend_asn = {
    order = 20
    asn   = 30
  }
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Action Rule Profile object.
* `annotation` - (Optional) Annotation of the Action Rule Profile object.
* `description` - (Optional) Description of the Action Rule Profile object.
* `name_alias` - (Optional) Name alias of the Action Rule Profile object.
* `set_route_tag` - (Optional) Set Route Tag of the Action Rule Profile object. Can not be configured along with `multipath`.
* `set_preference` - (Optional) Set Preference of the Action Rule Profile object.
* `set_weight` - (Optional) Set Weight of the Action Rule Profile object.
* `set_metric` - (Optional) Set Metric of the Action Rule Profile object.
* `set_metric_type` - (Optional) Set Metric Type of the Action Rule Profile object. Allowed values are `ospf-type1`, `ospf-type2`.
* `set_next_hop` - (Optional) Set Next Hop of the Action Rule Profile object.
* `set_communities` - (Optional) A block representing the attributes of Set Communities object. Type: Block.
  * `criteria` - (Optional) Criteria of the Set Communities object. Allowed values are `append` or `replace`.
  * `community` - (Optional) Community of the Set Communities object. Allowed input formats are `regular:as2-nn2:4:15`, `extended:as4-nn2:5:16`, `no-export` and `no-advertise`.
* `next_hop_propagation` - (Optional) Next Hop Propagation of the Action Rule Profile object. Allowed values are `yes` or `no`.
* `multipath` - (Optional) Multipath of the Action Rule Profile object. Allowed values are `yes` or `no`. Can not be configured along with `set_route_tag`.
* `saspath_prepend_last_as` - (Optional) Set As Path - Prepend Last-AS of the Action Rule Profile object. The value must be between 1 to 10.
* `saspath_prepend_asn` - (Optional) A block representing the attributes of Set As Path - Prepend AS of the Action Rule Profile object. Type: Block.
  * `asn` - (Optional) ASN of the Set As Path - Prepend AS object.
  * `order` - (Optional) Order of the Set As Path - Prepend AS object. Order must be between 0 to 31.

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
