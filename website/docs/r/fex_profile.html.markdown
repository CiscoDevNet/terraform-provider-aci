---
layout: "aci"
page_title: "ACI: aci_fex_profile"
sidebar_current: "docs-aci-resource-fex_profile"
description: |-
  Manages ACI FEX Profile
---

# aci_fex_profile

Manages ACI FEX Profile

## Example Usage

```hcl

resource "aci_fex_profile" "example" {
  name        = "fex_prof"
  annotation  = "example"
  name_alias  = "example"
  description = "from terraform"
}

```

## Argument Reference

- `name` - (Required) The FEX profile name.
- `annotation` - (Optional) Specifies the annotation of the FEX profile name.
- `description` - (Optional) Specifies the description of the FEX profile name.
- `name_alias` - (Optional) Specifies the alias name of the FEX profile name.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the FEX Profile.

## Importing

An existing FEX Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_fex_profile.example <Dn>
```
