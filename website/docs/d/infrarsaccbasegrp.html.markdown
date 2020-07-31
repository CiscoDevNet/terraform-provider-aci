---
layout: "aci"
page_title: "ACI: aci_access_group"
sidebar_current: "docs-aci-data-source-access_group"
description: |-
  Data source for ACI Access Group
---

# aci_access_group #
Data source for ACI Access Group

## Example Usage ##

```hcl

data "aci_access_group" "example" {
  access_port_selector_dn  = "${aci_access_port_selector.example.id}"
}

```


## Argument Reference ##
* `access_port_selector_dn` - (Required) Distinguished name of parent AccessPortSelector object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Access Access Group.
* `annotation` - (Optional) annotation for object access_access_group.
* `fex_id` - (Optional) interface policy group fex id
* `tdn` - (Optional) interface policy group's target rn
