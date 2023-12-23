---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_port_selector"
sidebar_current: "docs-aci-data-source-aci_spine_port_selector"
description: |-
  Data source for ACI Spine Port Selector
---

# aci_spine_port_selector #
Data source for ACI Spine Port Selector

## Example Usage ##

```hcl

data "aci_spine_port_selector" "example" {
  spine_profile_dn  = aci_spine_profile.example.id
  tdn               = aci_spine_interface_profile.example.id
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent Spine Profile.
* `tdn` - (Required) tDn of the Spine Interface Profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spine Port selector.
* `annotation` - (Optional) Annotation for Spine Port selector.
