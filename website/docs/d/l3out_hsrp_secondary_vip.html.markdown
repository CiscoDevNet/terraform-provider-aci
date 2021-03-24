---
layout: "aci"
page_title: "ACI: aci_l3out_hsrp_secondary_vip"
sidebar_current: "docs-aci-data-source-l3out_hsrp_secondary_vip"
description: |-
  Data source for ACI L3out HSRP Secondary VIP
---

# aci_l3out_hsrp_secondary_vip

Data source for ACI L3out HSRP Secondary VIP

## Example Usage

```hcl
data "aci_l3out_hsrp_secondary_vip" "example" {
  hsrp_group_profile_dn = "uni/tn-example/out-example/lnodep-example/lifp-example/hsrpIfP/hsrpGroupP-example"
  ip = "example"
}
```

## Argument Reference

- `hsrp_group_profile_dn` - (Required) Distinguished name of parent HSRP group profile object.
- `ip` - (Required) IP of Object L3out HSRP Secondary VIP.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out HSRP Secondary VIP.
- `annotation` - (Optional) Annotation for object L3out HSRP Secondary VIP.
- `description` - (Optional) Description for object L3out HSRP Secondary VIP.
- `config_issues` - (Optional) Configuration Issues.
- `ip` - (Optional) IP address.
- `name_alias` - (Optional) Name alias for object L3out HSRP Secondary VIP.
