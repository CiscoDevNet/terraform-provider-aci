---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_hsrp_interface_policy"
sidebar_current: "docs-aci-data-source-hsrp_interface_policy"
description: |-
  Data source for ACI HSRP Interface Policy
---

# aci_hsrp_interface_policy

Data source for ACI HSRP Interface Policy

## Example Usage

```hcl
data "aci_hsrp_interface_policy" "check" {
  tenant_dn = aci_tenant.tenentcheck.id
  name      = "one"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of HSRP interface policy object.

## Attribute Reference

- `id` - Attribute id set to the Dn of HSRP interface policy object.
- `annotation` - Annotation for HSRP interface policy object.
- `ctrl` - Control state for HSRP interface policy object.
- `delay` - Administrative port delay for HSRP interface policy object.
- `name_alias` - Name alias for HSRP interface policy object.
- `reload_delay` - Reload delay for HSRP interface policy object.
- `description` - Description for object HSRP interface policy object.
