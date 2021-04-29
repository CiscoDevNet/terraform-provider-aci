---
layout: "aci"
page_title: "ACI: aci_l3out_vpc_member"
sidebar_current: "docs-aci-data-source-l3out_vpc_member"
description: |-
  Data source for ACI L3out VPC Member
---

# aci_l3out_vpc_member

Data source for ACI L3out VPC Member

## Example Usage

```hcl
data "aci_l3out_vpc_member" "example" {

  leaf_port_dn  = "${aci_l3out_path_attachment.example.id}"
  side  = "A"
}
```

## Argument Reference

- `leaf_port_dn` - (Required) Distinguished name of parent leaf port object.
- `side` - (Required) side of Object l3out vpc member.  
  Allowed values: "A" and "B". Default value: "A".

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3out VPC Member.
- `addr` - (Optional) peer address
- `annotation` - (Optional) annotation for object l3out vpc member.
- `ipv6_dad` - (Optional) ipv6_dad for object l3out vpc member.
- `ll_addr` - (Optional) override of system generated ipv6 link-local address
- `name_alias` - (Optional) name_alias for object l3out vpc member.
- `side` - (Optional) node id
