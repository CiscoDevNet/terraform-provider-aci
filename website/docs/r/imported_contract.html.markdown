---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_imported_contract"
sidebar_current: "docs-aci-resource-imported_contract"
description: |-
  Manages ACI Imported Contract
---

# aci_imported_contract

Manages ACI Imported Contract

## Example Usage

```hcl
resource "aci_imported_contract" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "example"
  annotation  = "tag_imported_contract"
  description = "from terraform"
  name_alias  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of the imported contract.
- `annotation` - (Optional) Specifies the annotation of the imported contract.
- `description` - (Optional) Specifies the description of the imported contract.
- `name_alias` - (Optional) Specifies the alias-name of the imported contract.

- `relation_vz_rs_if` - (Optional) Relation to class vzACtrct. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Imported Contract.

## Importing

An existing Imported Contract can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_imported_contract.example <Dn>
```
