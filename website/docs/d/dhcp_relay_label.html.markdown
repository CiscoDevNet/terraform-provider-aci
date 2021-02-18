---
layout: "aci"
page_title: "ACI: aci_dhcp_relay_label"
sidebar_current: "docs-aci-data-source-dhcp_relay_label"
description: |-
  Data source for ACI DHCP Relay Label
---

# aci_dhcp_relay_label

Data source for ACI DHCP Relay Label

## Example Usage

```hcl
data "aci_dhcp_relay_label" "example" {

  bridge_domain_dn  = "${aci_bridge_domain.example.id}"
  name  = "example"
}
```

## Argument Reference

- `bridge_domain_dn` - (Required) Distinguished name of parent BridgeDomain object.
- `name` - (Required) Name of Object dhcp_relay_label.

## Attribute Reference

- `id` - Attribute id set to the Dn of the DHCP Relay Label.
- `annotation` - (Optional) Annotation for object dhcp_relay_label.
- `name_alias` - (Optional) name_alias for object dhcp_relay_label.
- `owner` - (Optional) Owner of the target relay servers.
- `tag` - (Optional) Label color.
