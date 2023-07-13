---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_pim_interface_policy"
sidebar_current: "docs-aci-data-source-pim_interface_policy"
description: |-
  Data source for ACI PIM Interface Policy
---

# aci_pim_interface_policy #

Data source for ACI PIM Interface Policy

## API Information ##

* `Class` - pimIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/pimifpol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> PIM

## Example Usage ##

```hcl
data "aci_pim_interface_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the PIM Interface Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the PIM Interface Policy.
* `annotation` - (Read-Only) Annotation of the PIM Interface Policy object.
* `name_alias` - (Read-Only) Name Alias of the PIM Interface Policy object.
* `auth_type` - (Read-Only) Authentication type.
* `auth_key` - (Read-Only) Secure authentication key.
* `control_state` - (Read-Only) Interface controls.
* `designated_router_delay` - (Read-Only) Designated router delay.
* `designated_router_priority` - (Read-Only) Designated router priority.
* `hello_interval` - (Read-Only) Hello traffic policy.
* `join_prune_interval` - (Read-Only) Join Prune Traffic Policy.
* `inbound_join_prune_filter_policy` - (Read-Only) Inbound join prune filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol).
* `outbound_join_prune_filter_policy` - (Read-Only) Outbound join prune filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol).
* `neighbor_filter_policy` - (Read-Only) Neighbor filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol).
