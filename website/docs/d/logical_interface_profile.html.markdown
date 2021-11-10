---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_interface_profile"
sidebar_current: "docs-aci-data-source-logical_interface_profile"
description: |-
  Data source for ACI Logical Interface Profile
---

# aci_logical_interface_profile

Data source for ACI Logical Interface Profile

## Example Usage

```hcl
data "aci_logical_interface_profile" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  name  = "example"
}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of parent Logical Node Profile object.
- `name` - (Required) Name of Object logical interface profile.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Interface Profile.
- `annotation` - (Optional) Annotation for object logical interface profile.
- `description` - (Optional) Description for object logical interface profile.
- `name_alias` - (Optional) Name alias for object logical interface profile.
- `prio` - (Optional) QoS priority class id.
- `tag` - (Optional) Specifies the color of a policy label.

