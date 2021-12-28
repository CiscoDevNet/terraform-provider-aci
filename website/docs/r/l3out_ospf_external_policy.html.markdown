---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_ospf_external_policy"
sidebar_current: "docs-aci-resource-l3out_ospf_external_policy"
description: |-
  Manages ACI L3-out OSPF External Policy
---

# aci_l3out_ospf_external_policy

Manages ACI L3-out OSPF External Policy

## Example Usage

```hcl
resource "aci_l3out_ospf_external_policy" "example" {
  l3_outside_dn  = aci_l3_outside.example.id
  annotation     = "example"
  description    = "from terraform"
  area_cost      = "1"
  area_ctrl      = ["redistribute", "summary"]
  area_id        = "0.0.0.1"
  area_type      = "nssa"
  multipod_internal = "no"
  name_alias     = "example"
}
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of the parent l3 outside object.
- `annotation` - (Optional) Annotation for object L3-out OSPF External Policy.
- `description` - (Optional) Description for object L3-out OSPF External Policy.
- `area_cost` - (Optional) The OSPF Area cost. Default value: "1".
- `area_ctrl` - (Optional) The controls of redistribution and summary LSA generation into NSSA and Stub areas.  
  Allowed values: "redistribute", "summary", "suppress-fa", "unspecified"  Default value: ["redistribute","summary"].
- `area_id` - (Optional) The OSPF Area ID.
- `area_type` - (Optional) The area type.  
  Allowed values: "nssa", "regular", "stub". Default value: "nssa".
- `multipod_internal` - (Optional) Start OSPF in WAN instance instead of the default. Value "yes" can be set only under infra tenant. Allowed values: "no", "yes". Default value: "no".
- `name_alias` - (Optional) Name alias for object L3-out OSPF External Policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3-out OSPF External Policy.

## Importing

An existing L3-out OSPF External Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_ospf_external_policy.example <Dn>
```
