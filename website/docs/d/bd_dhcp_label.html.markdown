---
layout: "aci"
page_title: "ACI: aci_bd_dhcp_label"
sidebar_current: "docs-aci-data-source-bd_dhcp_label"
description: |-
  Data source for ACI BD DHCP Label
---

# aci_bd_dhcp_label

Data source for ACI BD DHCP Label

## Example Usage

```hcl
data "aci_bd_dhcp_label" "example" {

  bridge_domain_dn  = "${aci_bridge_domain.example.id}"
  name  = "example"
}
```

## Argument Reference

- `bridge_domain_dn` - (Required) Distinguished name of parent Bridge Domain object.
- `name` - (Required) Name of Object bd_dhcp_label.

## Attribute Reference

- `id` - Attribute id set to the Dn of the BD DHCP Label.
- `annotation` - (Optional) Annotation for object bd_dhcp_label.
- `name_alias` - (Optional) Name alias for object bd_dhcp_label.
- `owner` - (Optional) Owner of the target relay servers.
- `tag` - (Optional) Label color.
