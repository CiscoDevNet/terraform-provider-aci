---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bgp_timers"
sidebar_current: "docs-aci-resource-bgp_timers"
description: |-
  Manages ACI BGP Timers
---

# aci_bgp_timers

Manages ACI BGP Timers

## Example Usage

```hcl
resource "aci_bgp_timers" "example" {
  tenant_dn    = aci_tenant.tenentcheck.id
  description  = "from terraform"
  name         = "one"
  annotation   = "example"
  gr_ctrl      = "helper"
  hold_intvl   = "189"
  ka_intvl     = "65"
  max_as_limit = "70"
  name_alias   = "aliasing"
  stale_intvl  = "15"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of bgp timers object.
- `annotation` - (Optional) Annotation for bgp timers object.
- `description` - (Optional) Description for bgp timers object.
- `gr_ctrl` - (Optional) Graceful restart enabled or helper only for bgp timers object. Default value is "helper".
- `hold_intvl` - (Optional) Time period before declaring neighbor down for bgp timers object. Default value is "180".
- `ka_intvl` - (Optional) Interval time between keepalive messages for bgp timers object. Default value is "60".
- `max_as_limit` - (Optional) Maximum AS limit for bgp timers object. Range of allowed values is "0" to "2000". Default value is "0".
- `name_alias` - (Optional) Name alias for bgp timers object. Default value is "default".
- `stale_intvl` - (Optional) Stale interval for routes advertised by peer for bgp timers object. Range of allowed values is "1" to "3600". Default value is "300".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Timers.

## Importing

An existing BGP Timers can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_timers.example <Dn>
```
