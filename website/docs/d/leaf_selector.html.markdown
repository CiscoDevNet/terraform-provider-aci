---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_selector"
sidebar_current: "docs-aci-data-source-leaf_selector"
description: |-
  Data source for ACI Leaf Selector
---

# aci_leaf_selector

Data source for ACI Leaf Selector

## Example Usage

```hcl
data "aci_leaf_selector" "example" {
  leaf_profile_dn          = aci_leaf_profile.example.id
  name                     = "example"
  switch_association_type  = "range"
}
```

## Argument Reference

- `leaf_profile_dn` - (Required) Distinguished name of parent Leaf Profile object.
- `name` - (Required) Name of Object switch association.
- `switch_association_type` - (Required) The leaf selector type.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Switch Association.
- `annotation` - (Optional) Annotation for object switch association.
- `description` - (Optional) Description for object switch association.
- `name_alias` - (Optional) Name alias for object switch association.
