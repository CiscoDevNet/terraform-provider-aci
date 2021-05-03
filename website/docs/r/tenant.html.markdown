---
layout: "aci"
page_title: "ACI: aci_tenant"
sidebar_current: "docs-aci-resource-tenant"
description: |-
  Manages ACI Tenant
---

# aci_tenant

Manages ACI Tenant

## Example Usage

```hcl
resource "aci_tenant" "example" {
  name        = "demo_tenant"
  description = "from terraform"
  annotation  = "tag"
  name_alias  = "tenant"
}
```

## Argument Reference

- `name` - (Required) Name of Object tenant.
- `annotation` - (Optional) Annotation for object tenant.
- `name_alias` - (Optional) Name alias for object tenant.
- `description` - (Optional) Description for object tenant.

- `relation_fv_rs_tn_deny_rule` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_tenant_mon_pol` - (Optional) Relation to class monEPGPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Tenant.

## Importing

An existing Tenant can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_tenant.example <Dn>
```
