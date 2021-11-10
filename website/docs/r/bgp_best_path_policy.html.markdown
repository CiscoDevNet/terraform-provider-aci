---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bgp_best_path_policy"
sidebar_current: "docs-aci-resource-bgp_best_path_policy"
description: |-
  Manages ACI BGP Best Path Policy
---

# aci_bgp_best_path_policy

Manages ACI BGP Best Path Policy

## Example Usage

```hcl
resource "aci_bgp_best_path_policy" "foobgp_best_path_policy" {
    tenant_dn   = aci_tenant.example.id
    name        = "example"
    annotation  = "example"
    description = "from terraform"
    ctrl        = "asPathMultipathRelax"
    name_alias  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object BGP Best Path Policy.
- `annotation` - (Optional) Annotation for object BGP Best Path Policy.
- `description` - (Optional) Description for object BGP Best Path Policy.
- `ctrl` - (Optional) The control state.  
  Allowed values: "asPathMultipathRelax", "0". Default Value: "0".
- `name_alias` - (Optional) Name alias for object BGP Best Path Policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Best Path Policy.

## Importing

An existing BGP Best Path Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_best_path_policy.example <Dn>
```
