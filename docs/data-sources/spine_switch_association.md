---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_switch_association"
sidebar_current: "docs-aci-data-source-aci_spine_switch_association"
description: |-
  Data source for ACI Spine Switch Association
---

# aci_spine_switch_association #
Data source for ACI Spine Switch Association

## Example Usage ##

```hcl

data "aci_spine_switch_association" "example" {
  spine_profile_dn                = aci_spine_profile.example.id
  name                            = "spine_switch_association_1"
  spine_switch_association_type   = "range"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent Spine Profile object.
* `name` - (Required) Name of Object Spine Switch Association.
* `spine_switch_association_type` - (Required) Spine association type of Object Spine Switch Association.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Switch Association.
* `description` - (Optional) Description for object Spine Switch Association.
* `annotation` - (Optional) Annotation for object Spine Switch Association.
* `name_alias` - (Optional) Name alias for object Spine Switch Association.
