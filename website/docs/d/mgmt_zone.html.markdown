---
layout: "aci"
page_title: "ACI: aci_in_b_managed_nodes_zone"
sidebar_current: "docs-aci-data-source-in_b_managed_nodes_zone"
description: |-
  Data source for ACI Management Zone
---

# aci_in_b_managed_nodes_zone

Data source for ACI Management Zone

## API Information

- `Class` - mgmtInBZone
- `Distinguished Named` - uni/infra/funcprof/grp-{name}/inbzone

## GUI Information

- `Location` -

## Example Usage

```hcl
data "aci_in_b_managed_nodes_zone" "example" {
  managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.example.id
  type = "in_band"
  name = "inb_zone"
}

data "aci_in_b_managed_nodes_zone" "example" {
  managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.example.id
  type = "out_of_band"
  name = "oob_zone"
}

```

## Argument Reference

- `managed_node_connectivity_group_dn` - (Required) Distinguished name of parent Managed Node Connectivity Group object.
- `type` - (Required) Type of the Management Zone. Allowed values: "in_band" and "out_of_band".
- `name` - (Required) Name of the Management Zone.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Management Zone.
- `annotation` - (Optional) Annotation of object Management Zone.
- `name_alias` - (Optional) Name Alias of object Management Zone.
- `description` - (Optional) Description of object Management Zone.
