---
layout: "aci"
page_title: "ACI: aci_spine_switch_association"
sidebar_current: "docs-aci-data-source-spine_switch_association"
description: |-
  Data source for ACI Spine Switch Association
---

# aci_spine_switch_association #
Data source for ACI Spine Switch Association

## Example Usage ##

```hcl

data "aci_spine_switch_association" "example" {
  spine_profile_dn                = "${aci_spine_profile.example.id}"
  name                            = "check"
  spine_switch_association_type   = "range"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent SpineProfile object.
* `name` - (Required) name of Object Spine Switch association.
* `spine_switch_association_type` - (Required) spine association type of Object Spine Switch association.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Switch Association.
* `annotation` - (Optional) annotation for object Spine Switch association.
* `name_alias` - (Optional) name alias for object Spine Switch association.
