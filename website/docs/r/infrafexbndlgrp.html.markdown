---
layout: "aci"
page_title: "ACI: aci_fex_bundle_group"
sidebar_current: "docs-aci-resource-fex_bundle_group"
description: |-
  Manages ACI Fex Bundle Group
---

# aci_fex_bundle_group #
Manages ACI Fex Bundle Group

## Example Usage ##

```hcl

resource "aci_fex_bundle_group" "example" {
  fex_profile_dn  = "${aci_fex_profile.example.id}"
  name            = "example"
  annotation      = "example"
  name_alias      = "example"
}

```


## Argument Reference ##
* `fex_profile_dn` - (Required) Distinguished name of parent FEXProfile object.
* `name` - (Required) name of Object fex_bundle_group.
* `annotation` - (Optional) annotation for object fex_bundle_group.
* `name_alias` - (Optional) name_alias for object fex_bundle_group.

* `relation_infra_rs_mon_fex_infra_pol` - (Optional) Relation to class monInfraPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_fex_bndl_grp_to_aggr_if` - (Optional) Relation to class pcAggrIf. Cardinality - ONE_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fex Bundle Group.

## Importing ##

An existing Fex Bundle Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fex_bundle_group.example <Dn>
```