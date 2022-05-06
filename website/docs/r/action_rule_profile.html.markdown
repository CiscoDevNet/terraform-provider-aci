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
  set_route_tag   = 100 # Can not be configured along with next_hop_propagation and multipath
  set_preference  = 100
  set_weight      = 100
  set_metric      = 100
  set_metric_type = "ospf-type1"
  set_next_hop    = "1.1.1.1"
  set_communities = {
    community = "no-advertise"
    criteria  = "replace"
  }
  set_as_path_prepend_last_as = 10
  set_as_path_prepend_as {
    order = 10
    asn   = 20
  }
  set_as_path_prepend_as {
    order = 20
    asn   = 30
  }
  set_dampening = {
    half_life        = 10 # Half time must be at least 9% of the maximum suppress time
    reuse           = 1
    suppress        = 10  # Suppress limit must be larger than reuse limit
    max_suppress_time = 100 # Max Suppress Time - should not be less than suppress limit
  }
}

resource "aci_action_rule_profile" "example2" {
  tenant_dn       = aci_tenant.example.id
  name            = "example2"
  set_preference  = 100
  set_weight      = 100
  set_metric      = 100
  set_metric_type = "ospf-type1"
  set_next_hop    = "1.1.1.1"
  set_communities = {
    community = "no-advertise"
    criteria  = "replace"
  }
  next_hop_propagation    = "yes" # Can not be configured along with set_route_tag
  multipath               = "yes" # Can not be configured along with set_route_tag
  set_as_path_prepend_last_as = 10
  set_as_path_prepend_as {
    order = 10
    asn   = 20
  }
  set_as_path_prepend_as {
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
* `set_route_tag` - (Optional) Set Route Tag of the Action Rule Profile object. Can not be configured along with `next_hop_propagation` and `multipath`. Type: Integer.
* `set_preference` - (Optional) Set Preference of the Action Rule Profile object. Type: Integer.
* `set_weight` - (Optional) Set Weight of the Action Rule Profile object. Type: Integer.
* `set_metric` - (Optional) Set Metric of the Action Rule Profile object. Type: Integer.
* `set_metric_type` - (Optional) Set Metric Type of the Action Rule Profile object. Allowed values are `ospf-type1`, `ospf-type2`.
* `set_next_hop` - (Optional) Set Next Hop of the Action Rule Profile object.
* `set_communities` - (Optional) A block representing the attributes of Set Communities object. Type: Block.
  * `criteria` - (Optional) Criteria of the Set Communities object. Allowed values are `append` or `replace`. Type: String.
  * `community` - (Optional) Community of the Set Communities object. Allowed input formats are `regular:as2-nn2:4:15`, `extended:as4-nn2:5:16`, `no-export` and `no-advertise`. Type: String.
* `next_hop_propagation` - (Optional) Next Hop Propagation of the Action Rule Profile object. Allowed values are `yes` or `no`. Can not be configured along with `set_route_tag`. Type: String.
* `multipath` - (Optional) Multipath of the Action Rule Profile object. Allowed values are `yes` or `no`. Can not be configured along with `set_route_tag`. Type: String.
* `set_as_path_prepend_last_as` - (Optional) Number of ASN to be prepended to AS Path of the Action Rule Profile object.
* `set_as_path_prepend_as` - (Optional) A block representing ASNs to be configured as Set As Path - Prepend AS of the Action Rule Profile object. Type: Block.
  * `asn` - ASN to be prepended to Set AS Path.
  * `order` - Order in which the ASN should be prepended to Set AS Path.
* `set_dampening` - (Optional) A block representing the attributes of Set Dampening object. Type: Block.
  * `half_life` - Half Life of the Set Dampening object, the maximum value for this field is 60 in minutes.
  * `reuse` - Reuse Limit of the Set Dampening object, the maximum value for this field is 20000.
  * `suppress` - Suppress Limit of the Set Dampening object, the maximum value for this field is 20000.
  * `max_suppress_time` - Max Suppress Time of the Set Dampening object, the maximum value for this field is 255 in minutes.

## Importing ##

An existing Action Rule Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_action_rule_profile.example <Dn>
```
