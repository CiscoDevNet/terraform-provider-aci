---
layout: "aci"
page_title: "ACI: aci_fex_bundle_group"
sidebar_current: "docs-aci-data-source-fex_bundle_group"
description: |-
  Data source for ACI Fex Bundle Group
---

# aci_fex_bundle_group

Data source for ACI Fex Bundle Group

## Example Usage

```hcl

data "aci_fex_bundle_group" "example" {
  fex_profile_dn  = aci_fex_profile.example.id
  name            = "example"
}

```

## Argument Reference

- `fex_profile_dn` - (Required) Distinguished name of parent FEX Profile object.
- `name` - (Required) Name of Object FEX bundle group.

## Attribute Reference

- `id` - Attribute id set to the Dn of the FEX Bundle Group.
- `annotation` - (Optional) Specifies the annotation of the FEX bundle group.
- `description` - (Optional) Specifies the description of the FEX bundle group.
- `name_alias` - (Optional) Specifies the alias name of the FEX bundle group.
