---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_taboo_contract"
sidebar_current: "docs-aci-resource-taboo_contract"
description: |-
  Manages ACI Taboo Contract
---

# aci_taboo_contract

Manages ACI Taboo Contract

## Example Usage

```hcl
resource "aci_taboo_contract" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example_contract"
  description = "from terraform"
  annotation  = "orchestrator:terraform"
  name_alias  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object Taboo Contract.
- `description` - (Optional) Description for object Taboo Contract.
- `annotation` - (Optional) Annotation for object Taboo Contract.
- `name_alias` - (Optional) Name alias for object Taboo Contract.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Taboo Contract.

## Importing

An existing Taboo Contract can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_taboo_contract.example <Dn>
```
