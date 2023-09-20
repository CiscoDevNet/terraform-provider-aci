---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment_secondary_ip"
sidebar_current: "docs-aci-resource-l3out_path_attachment_secondary_ip"
description: |-
  Manages ACI L3-out Path Attachment Secondary IP
---

# aci_l3out_path_attachment_secondary_ip

Manages ACI L3-out Path Attachment Secondary IP

## API Information ##

* `Class` - l3extIp
* `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/lnodep-{lnodep}/lifp-{lifp}/rspathL3OutAtt-[leaf_path_dn]/addr-[policy_ip_addr]

## GUI Information ##

* `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles -> SVI -> IPv4 Secondary / IPv6 Additional Addresses

## Example Usage

```hcl
resource "aci_l3out_path_attachment_secondary_ip" "example" {
  l3out_path_attachment_dn = aci_l3out_path_attachment.example.id
  addr                     = "10.0.0.1/24"
  annotation               = "example"
  description              = "from terraform"
  ipv6_dad                 = "disabled"
  name_alias               = "example"
  dhcp_relay               = "enabled"
}
```

## Argument Reference

- `l3out_path_attachment_dn` - (Required) Distinguished name of the parent L3-out Path Attachment object.
- `addr` - (Required) The peer IP address of the L3-out Path Attachment Secondary IP object.
- `description` - (Optional) Description of the L3-out Path Attachment Secondary IP object.
- `annotation` - (Optional) Annotation of the L3-out Path Attachment Secondary IP object.
- `ipv6_dad` - (Optional) IPv6 DAD of the L3-out Path Attachment Secondary IP object. Allowed values are "enabled" and "disabled". Default value is "enabled".
- `name_alias` - (Optional) Name alias of the L3-out Path Attachment Secondary IP object.
- `dhcp_relay` - (Optional) DHCP relay gateway of the L3-out Path Attachment Secondary IP object. Allowed values are "enabled" and "disabled". Default value is "disabled".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the L3 out Path Attachment Secondary IP.

## Importing

An existing L3-out Path Attachment Secondary IP can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_path_attachment_secondary_ip.example <Dn>
```
