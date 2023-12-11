---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_selector"
sidebar_current: "docs-aci-resource-aci_leaf_selector"
description: |-
  Manages ACI Leaf Selector
---

# aci_leaf_selector

Manages ACI Leaf Selector

## Example Usage

```hcl
resource "aci_leaf_selector" "example" {
  leaf_profile_dn         = aci_leaf_profile.example.id
  name                    = "example_leaf_selector"
  switch_association_type = "range"
  annotation              = "orchestrator:terraform"
  description             = "from terraform"
  name_alias              = "tag_leaf_selector"
}
```

## Argument Reference

- `leaf_profile_dn` - (Required) Distinguished name of parent Leaf Profile object.
- `name` - (Required) Name of Object switch association.
- `switch_association_type` - (Required) The leaf selector type.
  Allowed values: "ALL", "range", "ALL_IN_POD".
- `annotation` - (Optional) Annotation for object switch association.
- `description` - (Optional) Description for object switch association.
- `name_alias` - (Optional) Name alias for object switch association.

- `relation_infra_rs_acc_node_p_grp` - (Optional) Relation to class infraAccNodePGrp. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Switch Association.

## Importing

An existing Switch Association can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_leaf_selector.example <Dn>
```
