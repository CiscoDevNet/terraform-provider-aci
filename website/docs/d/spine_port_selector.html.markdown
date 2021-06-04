---
layout: "aci"
page_title: "ACI: aci_spine_port_selector"
sidebar_current: "docs-aci-data-source-spine_port_selector"
description: |-
  Data source for ACI Spine Port Selector
---

# aci_spine_port_selector #
Data source for ACI Spine Port Selector

## Example Usage ##

```hcl

data "aci_spine_port_selector" "example" {
  spine_profile_dn  = "${aci_spine_profile.example.id}"
  tdn               = "example"
}

```


## Argument Reference ##
* `spine_profile_dn` - (Required) Distinguished name of parent SpineProfile object.
* `tdn` - (Required) tDn of Object Interface profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the port selector.
* `annotation` - (Optional) Annotation for object port selector.

