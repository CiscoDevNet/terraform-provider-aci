---
layout: "aci"
page_title: "ACI: aci_dhcp_option_policy"
sidebar_current: "docs-aci-data-source-dhcp_option_policy"
description: |-
  Data source for ACI DHCP Option Policy.
---

# aci_dhcp_option_policy

Data source for ACI DHCP Option Policy.

## Example Usage

```hcl
data "aci_dhcp_option_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object  DHCP Option Policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the DHCP Option Policy.
- `annotation` - (Optional) Annotation for object  DHCP Option Policy.
- `name_alias` - (Optional) Name alias for object  DHCP Option Policy.
