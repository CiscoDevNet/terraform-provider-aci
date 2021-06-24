---
layout: "aci"
page_title: "ACI: aci_l3out_ospf_external_policy"
sidebar_current: "docs-aci-data-source-l3out_ospf_external_policy"
description: |-
  Data source for ACI L3-out OSPF External Policy
---

# aci_l3out_ospf_external_policy

Data source for ACI L3-out OSPF External Policy

## Example Usage

```hcl
data "aci_l3out_ospf_external_policy" "example" {
  l3_outside_dn  = aci_l3_outside.example.id
}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of the parent l3 outside object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3-out OSPF External Policy.
- `annotation` - (Optional) Annotation for object L3-out OSPF External Policy.
- `description` - (Optional) Description for object L3-out OSPF External Policy.
- `area_cost` - (Optional) The OSPF Area cost.
- `area_ctrl` - (Optional) The controls of redistribution and summary LSA generation into NSSA and Stub areas.
- `area_id` - (Optional) The OSPF Area ID.
- `area_type` - (Optional) The area type.
- `multipod_internal` - (Optional) Start OSPF in WAN instance instead of the default.
- `name_alias` - (Optional) Name alias for object L3-out OSPF External Policy.
