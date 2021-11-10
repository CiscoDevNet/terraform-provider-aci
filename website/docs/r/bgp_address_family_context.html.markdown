---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bgp_address_family_context"
sidebar_current: "docs-aci-resource-bgp_address_family_context"
description: |-
  Manages ACI BGP Address Family Context
---

# aci_bgp_address_family_context

Manages ACI BGP Address Family Context

## Example Usage

```hcl
resource "aci_bgp_address_family_context" "example" {
  tenant_dn     = aci_tenant.tenentcheck.id
  name          = "one"
  description   = "from terraform"
  annotation    = "example"
  ctrl          = "host-rt-leak"
  e_dist        = "25"
  i_dist        = "198"
  local_dist    = "100"
  max_ecmp      = "18"
  max_ecmp_ibgp = "25"
  name_alias    = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of BGP address family context object.
- `description` - (Optional) Description for BGP address family context object.
- `annotation` - (Optional) Annotation for BGP address family context object.
- `ctrl` - (Optional) Control state for BGP address family context object. Allowed value is "host-rt-leak".
- `e_dist` - (Optional) Administrative distance of EBGP routes for BGP address family context object. Range of allowed values is "1" to "255". Default value is "20".
- `i_dist` - (Optional) Administrative distance of IBGP routes for BGP address family context object. Range of allowed values is "1" to "255". Default value is "200".
- `local_dist` - (Optional) Administrative distance of local routes for BGP address family context object. Range of allowed values is "1" to "255". Default value is "220".
- `max_ecmp` - (Optional) Maximum number of equal-cost paths for BGP address family context object.Range of allowed values is "1" to "64". Default value is "16".
- `max_ecmp_ibgp` - (Optional) Maximum ECMP IBGP for BGP address family context object. Range of allowed values is "1" to "64". Default value is "16".
- `name_alias` - (Optional) Name alias for BGP address family context object.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Address Family Context.

## Importing

An existing BGP Address Family Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_address_family_context.example <Dn>
```
