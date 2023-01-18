---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3_outside"
sidebar_current: "docs-aci-resource-l3_outside"
description: |-
  Manages ACI L3 Outside
---

# aci_l3_outside

Manages ACI L3 Outside

## API Information ##

* `Class` - l3extOut
* `Distinguished Name` - uni/tn-{tenant_name}/out-{l3_outside_name}

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs

## Example Usage

```hcl
resource "aci_l3_outside" "foo_l3_outside" {
  tenant_dn      = aci_tenant.terraform_tenant.id
  name           = "foo_l3_outside"
  enforce_rtctrl = ["export", "import"]
  target_dscp    = "unspecified"
  mpls_enabled   = "yes"

  // Relation to Route Control for Dampening
  relation_l3ext_rs_dampening_pol {
    tn_rtctrl_profile_dn = data.aci_route_control_profile.shared_route_control_profile.id
    af                   = "ipv6-ucast"
  }

  relation_l3ext_rs_dampening_pol {
    tn_rtctrl_profile_dn = data.aci_route_control_profile.shared_route_control_profile.id
    af                   = "ipv4-ucast"
  }

  // Target VRF object should belong to the parent tenant or be a shared object.
  relation_l3ext_rs_ectx = data.aci_vrf.default_vrf.id

  // Relation to Route Profile for Interleak - Interleak Policy object should belong to the parent tenant or be a shared object.
  relation_l3ext_rs_interleak_pol = data.aci_route_control_profile.shared_route_control_profile.id

  // Relation to L3 Domain
  relation_l3ext_rs_l3_dom_att = aci_l3_domain_profile.l3_domain_profile.id

  // Relation to Route Profile for Redistribution
  relation_l3extrs_redistribute_pol {
    target_dn = data.aci_route_control_profile.shared_route_control_profile.id
    source    = "static"
  }

  relation_l3extrs_redistribute_pol {
    target_dn = data.aci_route_control_profile.shared_route_control_profile.id
    source    = "direct"
  }
}
```

## Argument Reference

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L3 Outside object.
* `description`- (Optional) Description of the L3 Outside object.
* `annotation` - (Optional) Annotation of the L3 Outside object.
* `enforce_rtctrl` - (Optional) Enforce route control type. Allowed values are "import" and "export". Default is "export". Type - String.
* `name_alias` - (Optional) Name alias of the L3 Outside object.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the L3 Outside object. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".
* `mpls_enabled` - (Optional) Indiscate whether MPLS is enabled or not. Allowed values are "no", "yes". Default value is "no".
* `relation_l3ext_rs_dampening_pol` - (Optional) Relation to Route Control Profile for Dampening Policies (class rtctrlProfile). Can't configure multiple Dampening Policies for the same address-family. Cardinality - N_TO_M. Type - Block.
  * tn_rtctrl_profile_name - (Deprecated) Name of the Route Control Profile for Dampening Policies.
  * tn_rtctrl_profile_dn - (Optional) Distinguished name of the Route Control Profile for Dampening Policies.
  * af - (Optional) Address Family of the Dampening Policies. Allowed values are "ipv4-ucast" and "ipv6-ucast". Default is "ipv4-ucast".
* `relation_l3ext_rs_ectx` - (Optional) Relation to VRF (class fvCtx). Target VRF object should belong to the parent tenant or be a shared object. Cardinality - N_TO_ONE. Type - String.
* `relation_l3ext_rs_interleak_pol` - (Optional) Relation to Route Profile for Interleak (class rtctrlProfile). Interleak Policy object should belong to the parent tenant or be a shared object. Cardinality - N_TO_ONE. Type - String.
* `relation_l3ext_rs_l3_dom_att` - (Optional) Relation to a L3 Domain (class extnwDomP). Cardinality - N_TO_ONE. Type - String.
* `relation_l3extrs_redistribute_pol` - (Optional) A block representing the relation to a Route Profile for Redistribution (class rtctrlProfile). Type: Block.
  * `source` - (Optional) Route Map Source for the Route Profile for Redistribution. Allowed values are "attached-host", "direct", "static". Default value is "static".
  * `target_dn` - (Required) Distinguished name of the Route Control Profile for the Route Profile for Redistribution.
## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3 Outside.

## Importing

An existing L3 Outside can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3_outside.example <Dn>
```
