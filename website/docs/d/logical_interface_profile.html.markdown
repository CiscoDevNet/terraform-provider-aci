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

## API Information

- `Class` - l3extLIfP
- `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/lnodep-{logical_node_profile}/lifp-{logical_interface_profile}

## GUI Information

- `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage

```hcl
data "aci_logical_interface_profile" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  name  = "example"
}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of the parent Logical Node Profile object.
- `name` - (Required) Name of the object logical interface profile.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Interface Profile.
- `annotation` - (Optional) Annotation of the object logical interface profile.
- `description` - (Optional) Description of the object logical interface profile.
- `name_alias` - (Optional) Name alias of the object logical interface profile.
- `prio` - (Optional) QoS priority class id.
- `tag` - (Optional) Specifies the color of a policy label.

