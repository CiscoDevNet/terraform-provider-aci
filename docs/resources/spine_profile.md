---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_profile"
sidebar_current: "docs-aci-resource-aci_spine_profile"
description: |-
  Manages ACI Spine Profile
---

# aci_spine_profile #
Manages ACI Spine Profile

## Example Usage ##

```hcl

resource "aci_spine_profile" "example" {
  name        = "spine_profile_1"
  description = "from terraform"
  annotation  = "spine_profile_tag"
  name_alias  = "check"
}

```


## Argument Reference ##
* `name` - (Required) Name of Object Spine Profile.
* `description` - (Optional) Description for object Spine Profile.
* `annotation` - (Optional) Annotation for object Spine Profile.
* `name_alias` - (Optional) Name alias for object Spine Profile.

- `spine_selector` - (Optional) Spine Selector block to attach with the Spine Profile.
- `spine_selector.name` - (Required) Name of the Spine Selector.
- `spine_selector.switch_association_type` - (Required) Type of switch association.
  Allowed values: "ALL", "range", "ALL_IN_POD"
- `spine_selector.node_block` - (Optional) Node block to attach with Spine Selector.
- `spine_selector.node_block.name` - (Required) Name of the node block.
- `spine_selector.node_block.from_` - (Optional) Start of Node Block range. Range from 1 to 16000. Default value is "1".
- `spine_selector.node_block.to_` - (Optional) End of Node Block range. Range from 1 to 16000. Default value is "1".

* `relation_infra_rs_sp_acc_port_p` - (Optional) Relation to class infraSpAccPortP. Cardinality - N_TO_M. Type - Set of String.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Profile.

## Importing ##

An existing Spine Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_spine_profile.example <Dn>
```