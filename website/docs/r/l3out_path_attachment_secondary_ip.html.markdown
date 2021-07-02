---
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment_secondary_ip"
sidebar_current: "docs-aci-resource-l3out_path_attachment_secondary_ip"
description: |-
  Manages ACI L3-out Path Attachment Secondary IP
---

# aci_l3out_path_attachment_secondary_ip

Manages ACI L3-out Path Attachment Secondary IP

## Example Usage

```hcl
resource "aci_l3out_path_attachment_secondary_ip" "example" {
  l3out_path_attachment_dn  = aci_l3out_path_attachment.example.id
  addr  = "10.0.0.1/24"
  annotation  = "example"
  description = "from terraform"
  ipv6_dad = "disabled"
  name_alias  = "example"
}
```

## Argument Reference

- `l3out_path_attachment_dn` - (Required) Distinguished name of parent L3 out path attachment object.
- `addr` - (Required) The peer IP address.
- `description` - (Optional) Description for object L3 out path attachment secondary IP.
- `annotation` - (Optional) Annotation for object L3 out path attachment secondary IP.
- `ipv6_dad` - (Optional) IPv6 DAD for object L3 out path attachment secondary IP.  
  Allowed values: "disabled", "enabled". Default value is "enabled".
- `name_alias` - (Optional) Name alias for object L3 out path attachment secondary IP.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the L3 out Path Attachment Secondary IP.

## Importing

An existing L3-out Path Attachment Secondary IP can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_path_attachment_secondary_ip.example <Dn>
```
