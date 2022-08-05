---
subcategory: "L3Out"
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
  l3out_hsrp_interface_group_dn = aci_l3out_hsrp_interface_group.example.id
  ip = "example"
}
```

## Argument Reference

- `l3out_hsrp_interface_group_dn` - (Required) Distinguished name of parent HSRP group profile object.
- `ip` - (Required) IP of Object L3out HSRP Secondary VIP.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out HSRP Secondary VIP.
- `annotation` - (Optional) Annotation for object L3out HSRP Secondary VIP.
- `description` - (Optional) Description for object L3out HSRP Secondary VIP.
- `config_issues` - (Optional) Configuration Issues.
- `name_alias` - (Optional) Name alias for object L3out HSRP Secondary VIP.
