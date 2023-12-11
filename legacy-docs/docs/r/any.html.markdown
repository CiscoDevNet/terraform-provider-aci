---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_any"
sidebar_current: "docs-aci-resource-aci_any"
description: |-
  Manages ACI Any
---

# aci_any

Manages ACI Any

## Example Usage

```hcl
resource "aci_any" "example_vzany" {
  vrf_dn       = aci_vrf.example.id
  description  = "vzAny Description"
  annotation   = "tag_any"
  match_t      = "AtleastOne"
  name_alias   = "alias_any"
  pref_gr_memb = "disabled"
}
```

## Argument Reference

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `annotation` - (Optional) Annotation of the Any object.
* `description` - (Optional) Description of the Any object.
* `match_t` - (Optional) Represents the provider label match criteria. Allowed values are "All", "None", "AtmostOne" and "AtleastOne". Default value is "AtleastOne".
* `name_alias` - (Optional) Name alias of the Any object.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPgs can be divided in a the context can be divided into two subgroups. Allowed values are "disabled" and "enabled". Default is "disabled".
* `relation_vz_rs_any_to_cons` - (Optional) Relation to Consumed Contracts (vzBrCP class)
* `relation_vz_rs_any_to_cons_if` - (Optional) Relation to Consumed Contract Interfaces (vzCPIf class)
* `relation_vz_rs_any_to_prov` - (Optional) Relation to Provided Contracts (vzBrCP class).

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Any.

## Importing

An existing Any object can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_any.example <Dn>
```
