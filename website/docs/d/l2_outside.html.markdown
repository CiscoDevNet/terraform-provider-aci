---
layout: "aci"
page_title: "ACI: aci_l2_outside"
sidebar_current: "docs-aci-data-source-l2_outside"
description: |-
  Data source for ACI L2 Outside
---

# aci_l2_outside

Data source for ACI L2 Outside

## Example Usage

```hcl
data "aci_l2_outside" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) name of object l2 outside.

## Attribute Reference

- `id` - Attribute id set to the Dn of the l2 outside.
- `annotation` - (Optional) Annotation for object l2 outside.
- `name_alias` - (Optional) Name alias for object l2 outside.
- `target_dscp` - (Optional) Target dscp.
