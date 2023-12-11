---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment_secondary_ip"
sidebar_current: "docs-aci-data-source-aci_l3out_path_attachment_secondary_ip"
description: |-
  Data source for ACI L3-out Path Attachment Secondary IP
---

# aci_l3out_path_attachment_secondary_ip

Data source for ACI L3-out Path Attachment Secondary IP

## API Information ##

* `Class` - l3extIp
* `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/lnodep-{lnodep}/lifp-{lifp}/rspathL3OutAtt-[leaf_path_dn]/addr-[policy_ip_addr]

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles -> SVI -> IPv4 Secondary / IPv6 Additional Addresses

## Example Usage

```hcl
data "aci_l3out_path_attachment_secondary_ip" "example" {
  l3out_path_attachment_dn = aci_l3out_path_attachment.example.id
  addr                     = "10.0.0.1/24"
}
```

## Argument Reference

- `l3out_path_attachment_dn` - (Required) Distinguished name of the parent L3-out Path Attachment object.
- `addr` - (Required) The peer IP address of the L3-out Path Attachment Secondary IP object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3-out Path Attachment Secondary IP object.
- `description` - (Optional) Description of the L3-out Path Attachment Secondary IP object.
- `annotation` - (Optional) Annotation of the L3-out Path Attachment Secondary IP object.
- `ipv6_dad` - (Optional) IPv6 DAD of the L3-out Path Attachment Secondary IP object.
- `name_alias` - (Optional) Name alias of the L3-out Path Attachment Secondary IP object.
- `dhcp_relay` - (Optional) DHCP relay gateway of the L3-out Path Attachment Secondary IP object.