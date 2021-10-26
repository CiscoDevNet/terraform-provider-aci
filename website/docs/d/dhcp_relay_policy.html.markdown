---
layout: "aci"
page_title: "ACI: aci_dhcp_relay_policy"
sidebar_current: "docs-aci-data-source-dhcp_relay_policy"
description: |-
  Data source for ACI DHCP Relay Policy
---

# aci_dhcp_relay_policy

Data source for ACI DHCP Relay Policy.

## Example Usage

```hcl
data "aci_dhcp_relay_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object DHCP Relay Policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the DHCP Relay Policy.
- `annotation` - (Optional) Annotation for object DHCP Relay Policy.
- `description` - (Optional) Description for object DHCP Relay Policy.
- `mode` - (Optional) DHCP relay policy mode.
- `name_alias` - (Optional) Name alias for object DHCP Relay Policy.
- `owner` - (Optional) Owner of the target relay servers.
