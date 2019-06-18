---
layout: "aci"
page_title: "ACI: aci_leaf_profile"
sidebar_current: "docs-aci-resource-leaf_profile"
description: |-
  Manages ACI Leaf Profile
---

# aci_leaf_profile #
Manages ACI Leaf Profile

## Example Usage ##

```hcl
resource "aci_leaf_profile" "example" {
  name        = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object leaf_profile.
* `annotation` - (Optional) annotation for object leaf_profile.
* `name_alias` - (Optional) name_alias for object leaf_profile.

* `relation_infra_rs_acc_card_p` - (Optional) Relation to class infraAccCardP. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_acc_port_p` - (Optional) Relation to class infraAccPortP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Profile.

## Importing ##

An existing Leaf Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_profile.example <Dn>
```