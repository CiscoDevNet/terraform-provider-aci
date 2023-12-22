---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_ospf_timers"
sidebar_current: "docs-aci-data-source-aci_ospf_timers"
description: |-
  Data source for ACI OSPF Timers
---

# aci_ospf_timers

Data source for ACI OSPF Timers

## Example Usage

```hcl
data "aci_ospf_timers" "check" {
  tenant_dn = aci_tenant.tenentcheck.id
  name      = "one"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of OSPF timers object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the OSPF timers.
- `annotation` - Annotation for OSPF timers object.
- `description` - Description for object OSPF timers object.
- `bw_ref` - Ospf policy bandwidth for OSPF timers object.
- `ctrl` - Control state for OSPF timers object.
- `dist` - Preferred administrative distance for OSPF timers object.
- `gr_ctrl` - Graceful restart enabled or helper only for OSPF timers object.
- `lsa_arrival_intvl` - Minimum interval between the arrivals of LSAs for OSPF timers object.
- `lsa_gp_pacing_intvl` - LSA group pacing interval for OSPF timers object.
- `lsa_hold_intvl` - Throttle hold interval between LSAs for OSPF timers object.
- `lsa_max_intvl` - Throttle max interval between LSAs for OSPF timers object.
- `lsa_start_intvl` - Throttle start-wait interval between LSAs for OSPF timers object.
- `max_ecmp` - Maximum ECMP for OSPF timers object.
- `max_lsa_action` - Action to take when maximum LSA limit is reached for OSPF timers object.
- `max_lsa_num` - Maximum number of LSAs that are not self-generated for OSPF timers object.
- `max_lsa_reset_intvl` - Time until the sleep count is reset to zero for OSPF timers object.
- `max_lsa_sleep_cnt` - Number of times ospf can be placed in sleep state for OSPF timers object.
- `max_lsa_sleep_intvl` - Maximum LSA threshold for OSPF timers object.
- `max_lsa_thresh` - Maximum LSA threshold for OSPF timers object.
- `name_alias` - Name alias for OSPF timers object.
- `spf_hold_intvl` - Minimum hold time between spf calculations for OSPF timers object.
- `spf_init_intvl` - Initial delay interval for the spf schedule for OSPF timers object.
- `spf_max_intvl` - Maximum interval between SPF calculations for OSPF timers object.
