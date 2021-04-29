---
layout: "aci"
page_title: "ACI: aci_l3out_vpc_member"
sidebar_current: "docs-aci-resource-l3out_vpc_member"
description: |-
  Manages ACI L3out VPC Member
---

# aci_l3out_vpc_member

Manages ACI L3out VPC Member

## Example Usage

```hcl
resource "aci_l3out_vpc_member" "example" {

  leaf_port_dn  = "${aci_l3out_path_attachment.example.id}"
  side  = "A"
  addr  = "10.0.0.1"
  annotation  = "example"
  ipv6_dad = "enabled"
  ll_addr  = "::"
  name_alias  = "example"

}
```

## Argument Reference

- `leaf_port_dn` - (Required) Distinguished name of parent leaf port object.
- `side` - (Required) side of Object l3out vpc member.  
Allowed values: "A" and "B". Default value: "A".
- `addr` - (Optional) peer address.

- `annotation` - (Optional) annotation for object l3out vpc member.

- `ipv6_dad` - (Optional) ipv6_dad for object l3out vpc member.
  Allowed values: "disabled", "enabled". Default value: "enabled"
- `ll_addr` - (Optional) override of system generated ipv6 link-local address.

- `name_alias` - (Optional) name_alias for object l3out vpc member.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out VPC Member.

## Importing

An existing L3out VPC Member can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_vpc_member.example <Dn>
```
