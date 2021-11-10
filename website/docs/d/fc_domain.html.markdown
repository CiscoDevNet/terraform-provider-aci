---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_fc_domain"
sidebar_current: "docs-aci-data-source-fc_domain"
description: |-
  Data source for ACI FC Domain
---

# aci_fc_domain

Data source for ACI FC Domain

## Example Usage

```hcl
data "aci_fc_domain" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) Name of Object fibre channel domain.

## Attribute Reference

- `id` - Attribute id set to the Dn of the FC Domain.
- `annotation` - (Optional) Annotation for object fibre channel domain.
- `name_alias` - (Optional) Name alias for object fibre channel domain.
