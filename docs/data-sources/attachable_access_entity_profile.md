---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_attachable_access_entity_profile"
sidebar_current: "docs-aci-data-source-aci_attachable_access_entity_profile"
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
* `name` - (Required) Name of Object attachable access entity profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the attachable access entity profile.
* `annotation` - (Optional) Annotation for object attachable access entity profile.
* `name_alias` - (Optional) Name alias for object attachable access entity profile.
* `description` - (Optional) Description for object attachable access entity profile.