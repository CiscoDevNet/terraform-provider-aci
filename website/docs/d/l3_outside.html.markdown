---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3_outside"
sidebar_current: "docs-aci-data-source-l3_outside"
description: |-
  Data source for ACI L3 Outside
---

# aci_l3_outside #
Data source for ACI L3 Outside

## API Information ##

* `Class` - l3extOut
* `Distinguished Name` - uni/tn-{tenant_name}/out-{l3_outside_name}

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs

## Example Usage ##

```hcl
data "aci_l3_outside" "foo_l3_outside" {
  tenant_dn      = aci_tenant.terraform_tenant.id
  name           = "foo_l3_outside"
}
```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L3 Outside object.

## Attribute Reference

* `id` - Attribute id set to the Dn of the L3 Outside.
* `description`- (Optional) Description of the L3 Outside object.
* `annotation` - (Optional) Annotation of the L3 Outside object.
* `enforce_rtctrl` - (Optional) Enforce route control type of the L3 Outside object. 
* `name_alias` - (Optional) Name alias of the L3 Outside object.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the L3 Outside object.
* `relation_l3ext_rs_dampening_pol` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_M. Type - Block.
  * tn_rtctrl_profile_name - (Deprecated) Name of the Route Control Profile for Dampening Policies.
  * tn_rtctrl_profile_dn - (Optional) Distinguished name of the Route Control Profile for Dampening Policies.
  * af - (Optional) Address Family of the Dampening Policies.
* `relation_l3ext_rs_ectx` - (Optional) Relation to class fvCtx.
* `relation_l3ext_rs_out_to_bd_public_subnet_holder` - (Optional) Relation to class fvBDPublicSubnetHolder.
* `relation_l3ext_rs_interleak_pol` - (Optional) Relation to class rtctrlProfile.
* `relation_l3ext_rs_l3_dom_att` - (Optional) Relation to class extnwDomP.