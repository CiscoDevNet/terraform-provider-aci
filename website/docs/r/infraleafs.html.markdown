---
layout: "aci"
page_title: "ACI: aci_switch_association"
sidebar_current: "docs-aci-resource-switch_association"
description: |-
  Manages ACI Switch Association
---

# aci_switch_association #
Manages ACI Switch Association

## Example Usage ##

```hcl
resource "aci_switch_association" "example" {

  leaf_profile_dn  = "${aci_leaf_profile.example.id}"

  name  = "example"

  switch_association_type  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `leaf_profile_dn` - (Required) Distinguished name of parent LeafProfile object.
* `name` - (Required) name of Object switch_association.
* `switch_association_type` - (Required) switch_association_type of Object switch_association.
* `annotation` - (Optional) annotation for object switch_association.
* `name_alias` - (Optional) name_alias for object switch_association.
* `switch_association_type` - (Optional) leaf selector type

* `relation_infra_rs_acc_node_p_grp` - (Optional) Relation to class infraAccNodePGrp. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Switch Association.

## Importing ##

An existing Switch Association can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_switch_association.example <Dn>
```