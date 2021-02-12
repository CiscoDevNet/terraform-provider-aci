---
layout: "aci"
page_title: "ACI: aci_dhcp_option"
sidebar_current: "docs-aci-data-source-dhcp_option"
description: |-
  Data source for ACI DHCP Option
---

# aci_dhcp_option

Data source for ACI DHCP Option.

## Example Usage

```hcl
data "aci_dhcp_option" "example" {

  dhcp_option_policy_dn  = "${aci_dhcp_option_policy.example.id}"
  name  = "example"
}
```

## Argument Reference

- `dhcp_option_policy_dn` - (Required) Distinguished name of parent DHCPOptionPolicy object.
- `name` - (Required) Name of Object dhcp_option.

## Attribute Reference

- `id` - Attribute id set to the Dn of the DHCP Option.
- `annotation` - (Optional) Annotation for object dhcp_option.
- `data` - (Optional) DHCP option data
- `dhcp_option_id` - (Optional) DHCP option id
- `name_alias` - (Optional) name_alias for object dhcp_option.
