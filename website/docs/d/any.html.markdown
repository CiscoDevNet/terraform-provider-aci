---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_any"
sidebar_current: "docs-aci-data-source-any"
description: |-
  Data source for ACI Any
---

# aci_any

Data source for ACI Any

## Example Usage

```hcl
data "aci_any" "dev_any" {
  vrf_dn  = aci_vrf.dev_vrf.id
}
```

## Argument Reference

- `vrf_dn` - (Required) Distinguished name of parent VRF object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Any.
- `annotation` - (Optional) Annotation for object any.
- `description` - (Optional) Description for object any.
- `match_t` - (Optional) Represents the provider label match criteria.
- `name_alias` - (Optional) Name alias for object any.
- `pref_gr_memb` - (Optional) Represents parameter used to determine if EPgs can be divided in a the context can be divided into two subgroups.
