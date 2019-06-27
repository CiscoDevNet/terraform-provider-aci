---
layout: "aci"
page_title: "ACI: aci_attachable_access_entity_profile"
sidebar_current: "docs-aci-data-source-attachable_access_entity_profile"
description: |-
  Data source for ACI Attachable Access Entity Profile
---

# aci_attachable_access_entity_profile #
Data source for ACI Attachable Access Entity Profile

## Example Usage ##

```hcl
data "aci_attachable_access_entity_profile" "dev_ent_prof" {
  name  = "foo_ent_prof"
}
```
## Argument Reference ##
* `name` - (Required) name of Object attachable_access_entity_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Attachable Access Entity Profile.
* `annotation` - (Optional) annotation for object attachable_access_entity_profile.
* `name_alias` - (Optional) name_alias for object attachable_access_entity_profile.
