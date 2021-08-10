---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_ospf_timers"
sidebar_current: "docs-aci-resource-ospf_timers"
description: |-
  Manages ACI OSPF Timers
---

# aci_ospf_timers

Manages ACI OSPF Timers

## Example Usage

```hcl
resource "aci_ospf_timers" "example" {
  tenant_dn           = aci_tenant.tenentcheck.id
  name                = "one"
  annotation          = "example"
  description         = "from terraform"
  bw_ref              = "30000"
  ctrl                = ["pfx-suppress", "name-lookup"]
  dist                = "200"
  gr_ctrl             = "helper"
  lsa_arrival_intvl   = "2000"
  lsa_gp_pacing_intvl = "50"
  lsa_hold_intvl      = "1000"
  lsa_max_intvl       = "1000"
  lsa_start_intvl     = "5"
  max_ecmp            = "10"
  max_lsa_action      = "restart"
  max_lsa_num         = "20"
  max_lsa_reset_intvl = "5"
  max_lsa_sleep_cnt   = "5"
  max_lsa_sleep_intvl = "10"
  max_lsa_thresh      = "50"
  name_alias          = "example"
  spf_hold_intvl      = "100"
  spf_init_intvl      = "500"
  spf_max_intvl       = "10"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of OSPF timers object.
- `annotation` - (Optional) Annotation for OSPF timers object.
- `description` - (Optional) Description for OSPF timers object.
- `bw_ref` - (Optional) OSPF policy bandwidth for OSPF timers object. Range of allowed values is "1" to "4000000". Default value is "40000".
- `ctrl` - (Optional) Control state for OSPF timers object. It is in the form of a comma-separated string and allowed values are "name-lookup" and "pfx-suppress".
- `dist` - (Optional) Preferred administrative distance for OSPF timers object. Range of allowed values is "1" to "255". Default value is "110".
- `gr_ctrl` - (Optional) Graceful restart enabled or helper only for OSPF timers object. The allowed value is "helper". The default value is "helper". To deselect the option, just pass `gr_ctrl=""`
- `lsa_arrival_intvl` - (Optional) Minimum interval between the arrivals of lsas for OSPF timers object. The range of allowed values is "10" to "600000". The default value is "1000".
- `lsa_gp_pacing_intvl` - (Optional) LSA group pacing interval for OSPF timers object. The range of allowed values is "1" to "1800". The default value is "10".
- `lsa_hold_intvl` - (Optional) Throttle hold interval between LSAs for OSPF timers object. The range of allowed values is "50" to "30000". The default value is "5000".
- `lsa_max_intvl` - (Optional) throttle max interval between LSAs for OSPF timers object. The range of allowed values is "50" to "30000". The default value is "5000".
- `lsa_start_intvl` - (Optional) Throttle start-wait interval between LSAs for OSPF timers object. The range of allowed values is "0" to "5000". The default value is "0".
- `max_ecmp` - (Optional) Maximum ECMP for OSPF timers object. The range of allowed values is "1" to "64". The default value is "8".
- `max_lsa_action` - (Optional) Action to take when maximum LSA limit is reached for OSPF timers object. Allowed values are "reject", "log" and "restart". The default value is "reject".
- `max_lsa_num` - (Optional) Maximum number of LSAs that are not self-generated for OSPF timers object. The range of allowed values is "1" to "4294967295". The default value is "20000".
- `max_lsa_reset_intvl` - (Optional) Time until the sleep count is reset to zero for OSPF timers object. The range of allowed values is "1" to "1440". The default value is "10".
- `max_lsa_sleep_cnt` - (Optional) Number of times OSPF can be placed in a sleep state for OSPF timers object. The range of allowed values is "1" to "4294967295". The default value is "5".
- `max_lsa_sleep_intvl` - (Optional) Maximum LSA threshold for OSPF timers object. The range of allowed values is "1" to "1440". The default value is "5".
- `max_lsa_thresh` - (Optional) Maximum LSA threshold for OSPF timers object. The range of allowed values is "1" to "100". The default value is "75".
- `name_alias` - (Optional) Name alias for OSPF timers object.
- `spf_hold_intvl` - (Optional) Minimum hold time between SPF calculations for OSPF timers object. The range of allowed values is "1" to "600000". The default value is "1000".
- `spf_init_intvl` - (Optional) Initial delay interval for the SPF schedule for OSPF timers object. The range of allowed values is "1" to "600000". The default value is "200".
- `spf_max_intvl` - (Optional) Maximum interval between SPF calculations for OSPF timers object. The range of allowed values is "1" to "600000". The default value is "5000".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the OSPF Timers.

## Importing

An existing OSPF Timers can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_ospf_timers.example <Dn>
```
