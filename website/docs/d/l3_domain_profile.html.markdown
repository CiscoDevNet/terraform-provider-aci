---
layout: "aci"
page_title: "ACI: aci_l3_domain_profile"
sidebar_current: "docs-aci-data-source-l3_domain_profile"
description: |-
  Data source for ACI L3 Domain Profile
---

# aci_l3_domain_profile #
Data source for ACI L3 Domain Profile

## Example Usage ##

```hcl
data "aci_l3_domain_profile" "example" {
  name  = "l3_domain_profile_1"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object l3_domain_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L3 Domain Profile.
* `annotation` - (Optional) Annotation for object l3_domain_profile.
* `name_alias` - (Optional) Name alias for object l3_domain_profile.
