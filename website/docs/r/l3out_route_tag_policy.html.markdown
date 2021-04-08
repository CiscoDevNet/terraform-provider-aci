---
layout: "aci"
page_title: "ACI: aci_l3out_route_tag_policy"
sidebar_current: "docs-aci-resource-l3out_route_tag_policy"
description: |-
  Manages ACI L3out Route Tag Policy
---

# aci_l3out_route_tag_policy

Manages ACI L3out Route Tag Policy

## Example Usage

```hcl
resource "aci_l3out_route_tag_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  tag  = "1"

}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of Object L3out route tag policy.
- `annotation` - (Optional) Annotation for object L3out route tag policy.

- `name_alias` - (Optional) Name alias for object L3out route tag policy.

- `tag` - (Optional) Tagged number. Default value: "4294967295".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3out Route Tag Policy.

## Importing

An existing L3out Route Tag Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3out_route_tag_policy.example <Dn>
```
