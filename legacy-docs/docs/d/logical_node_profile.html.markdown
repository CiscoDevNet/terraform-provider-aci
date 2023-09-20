---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_node_profile"
sidebar_current: "docs-aci-data-source-logical_node_profile"
description: |-
  Data source for ACI Logical Node Profile
---

# aci_logical_node_profile

Data source for ACI Logical Node Profile

## Example Usage

```hcl
data "aci_logical_node_profile" "example" {
  l3_outside_dn  = aci_l3_outside.example.id
  name  = "example"
}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of parent l3-outside object.
- `name` - (Required) Name of Object logical node profile.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Node Profile.
- `annotation` - (Optional) Annotation for object logical node profile.
- `description` - (Optional) Description for object logical node profile.
- `config_issues` - (Optional) Bitmask representation of the configuration issues found during the endpoint group deployment.
- `name_alias` - (Optional) Name alias for object logical node profile.
- `tag` - (Optional) Specifies the color of a policy label.
- `target_dscp` - (Optional) Node level DSCP value.
