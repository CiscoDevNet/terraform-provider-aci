---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_taboo_contract"
sidebar_current: "docs-aci-data-source-taboo_contract"
description: |-
  Data source for ACI Taboo Contract
---

# aci_taboo_contract

Data source for ACI Taboo Contract

## Example Usage

```hcl
data "aci_taboo_contract" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object Taboo Contract.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Taboo Contract.
- `description` - (Optional) Description for object Taboo Contract.
- `annotation` - (Optional) Annotation for object Taboo Contract.
- `name_alias` - (Optional) Name alias for object Taboo Contract.
