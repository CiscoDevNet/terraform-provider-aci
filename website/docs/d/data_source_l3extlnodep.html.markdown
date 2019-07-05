---
layout: "aci"
page_title: "ACI: aci_logical_node_profile"
sidebar_current: "docs-aci-data-source-logical_node_profile"
description: |-
  Data source for ACI Logical Node Profile
---

# aci_logical_node_profile #
Data source for ACI Logical Node Profile

## Example Usage ##

```hcl
data "aci_logical_node_profile" "example" {

  l3_outside_dn  = "${aci_l3_outside.example.id}"

  name  = "example"
}
```
## Argument Reference ##
* `l3_outside_dn` - (Required) Distinguished name of parent L3Outside object.
* `name` - (Required) name of Object logical_node_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Logical Node Profile.
* `annotation` - (Optional) annotation for object logical_node_profile.
* `config_issues` - (Optional) configuration issues
* `name_alias` - (Optional) name_alias for object logical_node_profile.
* `tag` - (Optional) label color
* `target_dscp` - (Optional) target dscp
