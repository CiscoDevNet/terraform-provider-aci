---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_igmp_interface_policy"
sidebar_current: "docs-aci-data-source-igmp_interface_policy"
description: |-
  Data source for ACI IGMP Interface Policy
---

# aci_igmp_interface_policy #

Data source for ACI IGMP Interface Policy

## API Information ##

* `Class` - igmpIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/igmpIfPol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> IGMP Interface

## Example Usage ##

```hcl
data "aci_igmp_interface_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the IGMP Interface Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the IGMP Interface Policy.
* `annotation` - (Read-Only) Annotation of the IGMP Interface Policy object.
* `name_alias` - (Read-Only) Name Alias of the IGMP Interface Policy object.
* `group_timeout` - (Read-Only) Group timeout.
* `control` - (Read-Only) Interface Control.
* `last_member_count` - (Read-Only) Last member query count.
* `last_member_response_time` - (Read-Only) Last member response time.
* `querier_timeout` - (Read-Only) Querier timeout.
* `query_interval` - (Read-Only) Query interval.
* `robustness_variable` - (Read-Only) Robustness factor.
* `response_interval` - (Read-Only) Query response interval.
* `startup_query_count` - (Read-Only) Startup query count.
* `startuo_query_interval` - (Read-Only) Startup query interval.
* `version` - (Read-Only) Interface version.
* `maximum_mulitcast_entries` - (Read-Only) Maximum Multicast Entries. Type: String.
* `reserved_mulitcast_entries` - (Read-Only) Reserved Multicast Entries. Type: String.
* `state_limit_route_map` - (Read-Only) State limit route map which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol).
* `report_policy_route_map` - (Read-Only) Report policy route map which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol).
* `static_report_route_map` - (Read-Only) Static report policy route map which represents the relationship to a PIM Route Map Filter (class rtdmcARtMapPol).
