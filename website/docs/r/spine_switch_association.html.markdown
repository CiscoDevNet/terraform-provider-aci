---
layout: "aci"
page_title: "ACI: aci_spine_switch_association"
sidebar_current: "docs-aci-resource-spine_switch_association"
description: |-
  Manages ACI Spine Switch Association
---

# aci_spine_switch_association #
Manages ACI Spine Association

## Example Usage ##

```hcl

resource "aci_spine_switch_association" "example" {
  spine_profile_dn                = "${aci_spine_profile.example.id}"
  name                            = "check"
  description                     = "hello"
  spine_switch_association_type   = "range"
  annotation                      = "example"
  name_alias                      = "example"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent Spine Profile object.
* `name` - (Required) name of Object spine Switch association.
* `spine_switch_association_type` - (Required) spine association type of Object spine Switch association.
Allowed values: "ALL", "range", "ALL_IN_POD"
* `annotation` - (Optional) annotation for object spine Switch association.
* `name_alias` - (Optional) name alias for object spine Switch association.

* `relation_infra_rs_spine_acc_node_p_grp` - (Optional) Relation to class infraSpineAccNodePGrp. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Association.

## Importing ##

An existing Spine Association can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_association.example <Dn>
```