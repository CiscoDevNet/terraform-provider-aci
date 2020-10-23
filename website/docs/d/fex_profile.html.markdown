---
layout: "aci"
page_title: "ACI: aci_fex_profile"
sidebar_current: "docs-aci-data-source-fex_profile"
description: |-
  Data source for ACI FEX Profile
---

# aci_fex_profile #
Data source for ACI FEX Profile

## Example Usage ##

```hcl

data "aci_fex_profile" "example" {
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object fex_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the FEX Profile.
* `annotation` - (Optional) annotation for object fex_profile.
* `name_alias` - (Optional) name_alias for object fex_profile.
