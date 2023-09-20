---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_fex_profile"
sidebar_current: "docs-aci-data-source-fex_profile"
description: |-
  Data source for ACI FEX Profile
---

# aci_fex_profile

Data source for ACI FEX Profile

## Example Usage

```hcl

data "aci_fex_profile" "example" {
  name  = "example"
}

```

## Argument Reference

- `name` - (Required) The fex profile name.

## Attribute Reference

- `id` - Attribute id set to the Dn of the FEX Profile.
- `annotation` - (Optional) Specifies the annotation of the policy definition.
- `name_alias` - (Optional) Specifies the description of the policy definition.
- `name_alias` - (Optional) Specifies the alias name of the policy definition.
