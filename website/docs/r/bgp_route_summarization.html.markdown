---
layout: "aci"
page_title: "ACI: aci_bgp_route_summarization"
sidebar_current: "docs-aci-resource-bgp_route_summarization"
description: |-
  Manages ACI BGP Route Summarization
---

# aci_bgp_route_summarization

Manages ACI BGP Route Summarization

## Example Usage

```hcl
resource "aci_bgp_route_summarization" "example" {

  tenant_dn   = aci_tenant.example.id
  name        = "example"
  annotation  = "example"
  description = "from terraform"
  attrmap     = "example"
  ctrl        = "as-set"
  name_alias  = "example"

}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object BGP route summarization.
- `annotation` - (Optional) Annotation for object BGP route summarization.
- `description` - (Optional) Description for object BGP route summarization.
- `attrmap` - (Optional) Summary attribute map.
- `ctrl` - (Optional) The control state.
  Allowed values: "as-set", "none". Default value: "none".
- `name_alias` - (Optional) Name alias for object BGP route summarization.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Route Summarization.

## Importing

An existing BGP Route Summarization can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_route_summarization.example <Dn>
```
