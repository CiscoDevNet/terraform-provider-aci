---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_dhcp_option"
sidebar_current: "docs-aci-data-source-aci_dhcp_option"
description: |-
  Data source for ACI DHCP Option
---

# aci_dhcp_option

Data source for ACI DHCP Option.

## Example Usage

```hcl
data "aci_dhcp_option" "example" {
  dhcp_option_policy_dn  = aci_dhcp_option_policy.example.id
  name  = "example"
}
```

## Argument Reference

- `dhcp_option_policy_dn` - (Required) Distinguished name of parent DHCP Option Policy object.
- `name` - (Required) Name of Object DHCP Option.

## Attribute Reference

- `id` - Attribute id set to the Dn of the DHCP Option.
- `annotation` - (Optional) Annotation for object DHCP Option.
- `data` - (Optional) DHCP Option data.
- `dhcp_option_id` - (Optional) DHCP Option id.
- `name_alias` - (Optional) Name alias for object DHCP Option.
