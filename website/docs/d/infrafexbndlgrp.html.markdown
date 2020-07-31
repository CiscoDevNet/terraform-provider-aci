---
layout: "aci"
page_title: "ACI: aci_fex_bundle_group"
sidebar_current: "docs-aci-data-source-fex_bundle_group"
description: |-
  Data source for ACI Fex Bundle Group
---

# aci_fex_bundle_group #
Data source for ACI Fex Bundle Group

## Example Usage ##

```hcl

data "aci_fex_bundle_group" "example" {
  fex_profile_dn  = "${aci_fex_profile.example.id}"
  name            = "example"
}

```


## Argument Reference ##
* `fex_profile_dn` - (Required) Distinguished name of parent FEXProfile object.
* `name` - (Required) name of Object fex_bundle_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Fex Bundle Group.
* `annotation` - (Optional) annotation for object fex_bundle_group.
* `name_alias` - (Optional) name_alias for object fex_bundle_group.
