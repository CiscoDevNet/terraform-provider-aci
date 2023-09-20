---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_any"
sidebar_current: "docs-aci-data-source-any"
description: |-
  Data source for ACI Any
---

# aci_any #

Data source for ACI Any

## Example Usage ##

```hcl
data "aci_any" "dev_any" {
  vrf_dn = aci_vrf.dev_vrf.id
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Any object.
* `annotation` - (Optional) Annotation of the Any object.
* `description` - (Optional) Description of the Any object.
* `match_t` - (Optional) Represents the provider label match criteria.
* `name_alias` - (Optional) Name alias of the Any object.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPgs can be divided in a the context can be divided into two subgroups.
* `relation_vz_rs_any_to_cons` - (Optional) Relation to Consumed Contracts (vzBrCP class)
* `relation_vz_rs_any_to_cons_if` - (Optional) Relation to Consumed Contract Interfaces (vzCPIf class)
* `relation_vz_rs_any_to_prov` - (Optional) Relation to Provided Contracts (vzBrCP class).