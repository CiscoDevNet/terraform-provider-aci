---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_tenant"
sidebar_current: "docs-aci-data-source-aci_tenant"
description: |-
  Data source for ACI Tenant
---

# aci_tenant

Data source for ACI Tenant

## Example Usage

```hcl
data "aci_tenant" "example" {
  name  = "dev_ten"
}
```

## Argument Reference

- `name` - (Required) Name of Object tenant.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Tenant.
- `annotation` - (Optional) Annotation for object tenant.
- `name_alias` - (Optional) Name alias for object tenant.
- `description` - (Optional) Description for object tenant.
