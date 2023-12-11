---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_access_bundle_policy_group"
sidebar_current: "docs-aci-data-source-aci_leaf_access_bundle_policy_group"
description: |-
  Data source for ACI leaf access bundle policy group
---

# aci_leaf_access_bundle_policy_group

Data source for ACI leaf access bundle policy group

## Example Usage

```hcl

data "aci_leaf_access_bundle_policy_group" "dev_pol_grp" {
  name  = "foo_pol_grp"
}

```

## Argument Reference

- `name` - (Required) The bundled ports group name. This name can be up to 64 alphanumeric characters. Note that you cannot change this name after the object has been saved.

## Attribute Reference

- `id` - Attribute id set to the Dn of the ACI leaf access bundle policy group.
- `annotation` - (Optional) Annotation for object leaf access bundle policy group.
- `description` - (Optional) Specifies a description of the policy definition.
- `lag_t` - (Optional) The bundled ports group link aggregation type: port channel vs virtual port channel.
- `name_alias` - (Optional) Name alias for object leaf access bundle policy group.
