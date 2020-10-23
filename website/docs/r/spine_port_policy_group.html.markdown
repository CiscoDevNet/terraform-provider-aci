---
layout: "aci"
page_title: "ACI: aci_spine_port_policy_group"
sidebar_current: "docs-aci-resource-spine_port_policy_group"
description: |-
  Manages ACI Spine Port Policy Group
---

# aci_spine_port_policy_group #
Manages ACI Spine Port Policy Group

## Example Usage ##

```hcl

resource "aci_spine_port_policy_group" "example" {
  name        = "example"
  annotation  = "example"
  name_alias  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object aci_spine_port_policy_group.
* `annotation` - (Optional) annotation for object aci_spine_port_policy_group.
* `name_alias` - (Optional) name_alias for object aci_spine_port_policy_group.

* `relation_infra_rs_h_if_pol` - (Optional) Relation to class fabricHIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_cdp_if_pol` - (Optional) Relation to class cdpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_copp_if_pol` - (Optional) Relation to class coppIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_att_ent_p` - (Optional) Relation to class infraAttEntityP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_macsec_if_pol` - (Optional) Relation to class macsecIfPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Access Port Policy Group.

## Importing ##

An existing Spine Access Port Policy Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_port_policy_group.example <Dn>
```