---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_igmp_interface_policy"
sidebar_current: "docs-aci-resource-aci_igmp_interface_policy"
description: |-
  Manages ACI IGMP Interface Policy
---

# aci_igmp_interface_policy #

Manages ACI IGMP Interface Policy

## API Information ##

* `Class` - igmpIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/igmpIfPol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> IGMP Interface

## Example Usage ##

```hcl
resource "aci_igmp_interface_policy" "example_igmp" {
  tenant_dn                  = aci_tenant.example.id
  name                       = "example_igmp"
  group_timeout              = "260"
  last_member_count          = "2"
  last_member_response_time  = "1"
  querier_timeout            = "255"
  query_interval             = "125"
  robustness_variable        = "2"
  response_interval          = "10"
  startup_query_count        = "2"
  startup_query_interval     = "31"
  version                    = "v2"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the IGMP Interface Policy object.
* `annotation` - (Optional) Annotation of the IGMP Interface Policy object.
* `name_alias` - (Optional) Name Alias of the IGMP Interface Policy object.
* `group_timeout` - (Optional) Group timeout. Allowed range is "3-65535" and the default value is "260".
* `control` - (Optional) Interface Control. Allowed values are "allow-v3-asm", "fast-leave", "rep-ll".  Type: List.
* `last_member_count` - (Optional) Last member query count. Allowed range is "1-5" and the default value is "2".
* `last_member_response_time` - (Optional) Last member response time. Allowed range is "1-25" and the default value is "1".
* `querier_timeout` - (Optional) Querier timeout. Allowed range is "1-65535" and the the default value is "255".
* `query_interval` - (Optional) Query interval. Allowed range is "1-18000" and the default value is "125".
* `robustness_variable` - (Optional) Robustness factor. Allowed range is "1-7" and the default value is "2".
* `response_interval` - (Optional) Query response interval. Allowed range is "1-25" and the default value is "10".
* `startup_query_count` - (Optional) Startup query count. Allowed range is "1-10" and the default value is "2".
* `startuo_query_interval` - (Optional) Startup query interval. Allowed range is "1-18000" and the default value is "31".
* `version` - (Optional) Interface version. Allowed values are "v2", "v3" and the default value is "v2". Type: String.
* `maximum_mulitcast_entries` - (Optional) Maximum Multicast Entries. Type: String.
* `reserved_mulitcast_entries` - (Optional) Reserved Multicast Entries. Type: String.
* `state_limit_route_map` - (Optional) State limit route map which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.
* `report_policy_route_map` - (Optional) Report policy route map which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.
* `static_report_route_map` - (Optional) Static report policy route map which represents the relationship to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.


## Importing ##

An existing IGMP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_igmp_interface_policy.example <Dn>
```