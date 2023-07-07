---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_igmp_interface_policy"
sidebar_current: "docs-aci-resource-igmp_interface_policy"
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
  tenant_dn          = aci_tenant.example.id
  name               = "exampleii"
  grp_timeout        = "260"
  last_mbr_cnt       = "2"
  last_mbr_resp_time = "1"
  querier_timeout    = "255"
  query_intvl        = "125"
  robust_fac         = "2"
  rsp_intvl          = "10"
  start_query_cnt    = "2"
  start_query_intvl  = "31"
  ver                = "v2"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the IGMP Interface Policy object.
* `annotation` - (Optional) Annotation of the IGMP Interface Policy object.
* `name_alias` - (Optional) Name Alias of the IGMP Interface Policy object.
* `grp_timeout` - (Optional) Group timeout. Allowed range is "3-65535" and the default value is "260".
* `if_ctrl` - (Optional) Interface Control. Allowed values are "allow-v3-asm", "fast-leave", "rep-ll".  Type: List.
* `last_mbr_cnt` - (Optional) Last member query count. Allowed range is "1-5" and the default value is "2".
* `last_mbr_resp_time` - (Optional) Last member response time. Allowed range is "1-25" and the default value is "1".
* `querier_timeout` - (Optional) Querier timeout. Allowed range is "1-65535" and the the default value is "255".
* `query_intvl` - (Optional) Query interval. Allowed range is "1-18000" and the default value is "125".
* `robust_fac` - (Optional) Robustness factor. Allowed range is "1-7" and the default value is "2".
* `rsp_intvl` - (Optional) Query response interval. Allowed range is "1-25" and the default value is "10".
* `start_query_cnt` - (Optional) Startup query count. Allowed range is "1-10" and the default value is "2".
* `start_query_intvl` - (Optional) Startup query interval. Allowed range is "1-18000" and the default value is "31".
* `ver` - (Optional) Interface version. Allowed values are "v2", "v3" and the default value is "v2". Type: String.


## Importing ##

An existing IGMP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_igmp_interface_policy.example <Dn>
```