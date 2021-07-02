---
layout: "aci"
page_title: "ACI: aci_bgp_timers"
sidebar_current: "docs-aci-data-source-bgp_timers"
description: |-
  Data source for ACI BGP Timers
---

# aci_bgp_timers

Data source for ACI BGP Timers

## Example Usage

```hcl
data "aci_bgp_timers" "check" {
  tenant_dn = aci_tenant.tenentcheck.id
  name      = "one"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of bgp timers object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the BGP Timers Policy.
- `annotation` - Annotation for bgp timers object.
- `description` - Description for bgp timers object.
- `gr_ctrl` - Graceful restart enabled or helper only for bgp timers object.
- `hold_intvl` - Time period before declaring neighbor down for bgp timers object.
- `ka_intvl` - Interval time between keepalive messages for bgp timers object.
- `max_as_limit` - Maximum AS limit for bgp timers object.
- `name_alias` - Name alias for bgp timers object.
- `stale_intvl` - Stale interval for routes advertised by peer for bgp timers object.
