---
layout: "aci"
page_title: "ACI: aci_leaf_selector"
sidebar_current: "docs-aci-data-source-leaf_selector"
description: |-
  Data source for ACI Leaf Selector
---

# aci_leaf_selector #
Data source for ACI Leaf Selector

## Example Usage ##

```hcl
data "aci_leaf_selector" "example" {

  leaf_profile_dn  = "${aci_leaf_profile.example.id}"

  name  = "example"

  switch_association_type  = "example"
}
```
## Argument Reference ##
* `leaf_profile_dn` - (Required) Distinguished name of parent LeafProfile object.
* `name` - (Required) name of Object switch_association.
* `switch_association_type` - (Required) switch_association_type of Object switch_association.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Switch Association.
* `annotation` - (Optional) annotation for object switch_association.
* `name_alias` - (Optional) name_alias for object switch_association.
* `switch_association_type` - (Optional) leaf selector type
