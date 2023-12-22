---
subcategory: "Node Management"
layout: "aci"
page_title: "ACI: aci_mgmt_zone"
sidebar_current: "docs-aci-resource-aci_in_b_managed_nodes_zone"
description: |-
  Manages ACI Management Zone
---

# aci_mgmt_zone

Manages ACI Management Zone

## API Information

- `Class` - mgmtInBZone and mgmtOoBZone
- `Distinguished Name` - uni/infra/funcprof/grp-{name}/inbzone and uni/infra/funcprof/grp-{name}/oobzone

## GUI Information

- `Location` - Tenants -> mgmt -> Managed Node Connectivity Groups -> Create Managed Node Connectivity Group -> Policy

## Example Usage

```hcl
resource "aci_mgmt_zone" "example" {
  managed_node_connectivity_group_dn  = aci_managed_node_connectivity_group.example.id
  type = "in_band"
  name = "inb_zone"
  name_alias = "zone_tag"
  annotation = "orchestrator:terraform"
  description = "from terraform"

  relation_mgmt_rs_in_b = aci_node_mgmt_epg.example.id // type = "in_band"

  relation_mgmt_rs_inb_epg = aci_node_mgmt_epg.example.id // type = "in_band"
}

resource "aci_mgmt_zone" "example2" {
  managed_node_connectivity_group_dn = aci_managed_node_connectivity_group.example.id
  type = "out_of_band"
  name = "oob_zone"
  name_alias = "zone_tag"
  annotation = "orchestrator:terraform"
  description = "from terraform"

  relation_mgmt_rs_oo_b = aci_node_mgmt_epg.example.id // type = "out_of_band"

  relation_mgmt_rs_oob_epg = aci_node_mgmt_epg.example.id // type = "out_of_band"
}
```

## Argument Reference

- `managed_node_connectivity_group_dn` - (Required) Distinguished name of parent ManagedNodeConnectivityGroup object.
- `name` - (Required) Name of the Management Zone.
- `type` - (Required) Type of the Management Zone. Allowed values: "in_band" and "out_of_band".
- `annotation` - (Optional) Annotation of object Management Zone.
- `name_alias` - (Optional) Name Alias of object Management Zone.
- `description` - (Optional) Description of object Management Zone.

- `relation_mgmt_rs_addr_inst` - (Optional) Represents the relation to a Relation to IP Address Pool (class fvnsAddrInst). A source relation to the IP address namespace/IP address range. Type: String.

### `type = "in_band"`

- `relation_mgmt_rs_in_b` - (Optional) Represents the relation to a In-Band Management EPg (class mgmtInB). Relationship to an in-band EPG Type: String.

- `relation_mgmt_rs_inb_epg` - (Optional) Represents the relation to a In-Band Management EPg (class mgmtInB). A source relation to the in-band management endpoint group. Type: String.

### `type = "out_of_band"`

- `relation_mgmt_rs_oo_b` - (Optional) Represents the relation to a Out-Of-Band Management EPg (class mgmtOoB). Relationship to an out-of-band EPG Type: String.

- `relation_mgmt_rs_oob_epg` - (Optional) Represents the relation to a Out-Of-Band Management EPg (class mgmtOoB). A source relation to an out-of-band management endpoint group. Type: String.

## Importing

An existing InBManagedNodesZone can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_mgmt_zone.example <Dn>
```
