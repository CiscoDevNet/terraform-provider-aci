---
subcategory: -
layout: "aci"
page_title: "ACI: aci_igmp_interface_policy"
sidebar_current: "docs-aci-data-source-igmp_interface_policy"
description: |-
  Data source for ACI IGMP Interface Policy
---

# aci_igmpinterface_policy #

Data source for ACI IGMP Interface Policy

## API Information ##

* `Class` - igmpIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/igmpIfPol-{name}

## GUI Information ##

* `Location` - 

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
* `grp_timeout` - (Read-Only) Group Timeout.
* `if_ctrl` - (Read-Only) Interface Control.
* `last_mbr_cnt` - (Read-Only) Last member query count.
* `last_mbr_resp_time` - (Read-Only) Last member response time.
* `querier_timeout` - (Read-Only) Querier Timeout. 
* `query_intvl` - (Read-Only) Query interval.
* `robust_fac` - (Read-Only) Robustness factor.
* `rsp_intvl` - (Read-Only) Query response interval.
* `start_query_cnt` - (Read-Only) Startup query count.
* `start_query_intvl` - (Read-Only) Startup query interval.
* `ver` - (Read-Only) Interface version.
