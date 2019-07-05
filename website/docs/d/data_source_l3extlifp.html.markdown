---
layout: "aci"
page_title: "ACI: aci_logical_interface_profile"
sidebar_current: "docs-aci-data-source-logical_interface_profile"
description: |-
  Data source for ACI Logical Interface Profile
---

# aci_logical_interface_profile #
Data source for ACI Logical Interface Profile

## Example Usage ##

```hcl
data "aci_logical_interface_profile" "example" {

  logical_node_profile_dn  = "${aci_logical_node_profile.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `logical_node_profile_dn` - (Required) Distinguished name of parent LogicalNodeProfile object.
* `name` - (Required) name of Object logical_interface_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Logical Interface Profile.
* `annotation` - (Optional) annotation for object logical_interface_profile.
* `name_alias` - (Optional) name_alias for object logical_interface_profile.
* `prio` - (Optional) qos priority class id
* `tag` - (Optional) label color
