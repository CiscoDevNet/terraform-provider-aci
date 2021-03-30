---
layout: "aci"
page_title: "ACI: aci_l3out_path_attachment_secondary_ip"
sidebar_current: "docs-aci-data-source-l3out_path_attachment_secondary_ip"
description: |-
  Data source for ACI L3-out Path Attachment Secondary IP
---

# aci_l3out_path_attachment_secondary_ip

Data source for ACI L3-out Path Attachment Secondary IP

## Example Usage

```hcl
data "aci_l3out_path_attachment_secondary_ip" "example" {
  l3out_path_attachment_dn  = "${aci_l3out_path_attachment.example.id}"
  addr  = "10.0.0.1/24"
}
```

## Argument Reference

- `l3out_path_attachment_dn` - (Required) Distinguished name of parent L3-out Path Attachment object.
- `addr` - (Required) The peer IP address.

## Attribute Reference

- `id` - Attribute id set to the Dn of the L3-out path attachment secondary IP.
- `annotation` - (Optional) Annotation for object L3-out path attachment secondary IP.
- `ipv6_dad` - (Optional) ipv6_dad for object L3-out path attachment secondary IP.
- `name_alias` - (Optional) name_alias for object L3-out path attachment secondary IP.
