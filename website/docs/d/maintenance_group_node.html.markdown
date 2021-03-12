---
layout: "aci"
page_title: "ACI: aci_maintenance_group_node"
sidebar_current: "docs-aci-data-source-maintenance_group_node"
description: |-
  Data source for ACI Node Block
---

# aci_maintenance_group_node #
Data source for ACI Node Block

## Example Usage ##

```hcl
data "aci_maintenance_group_node" "example" {
  pod_maintenance_group_dn = "${aci_pod_maintenance_group.example.id}"
  name                     = "example"
}
```


## Argument Reference ##

* `pod_maintenance_group_dn` - (Required) distinguished name of parent POD maintenance group object.
* `name` - (Required) name of maintenance group node object.



## Attribute Reference

* `id` - Attribute id set to the dn of the maintenance group node object.
* `annotation` - (Optional) annotation for maintenance group node object.
* `from_` - (Optional) from for maintenance group node object.
* `name_alias` - (Optional) name alias for maintenance group node object.
* `to_` - (Optional) to for maintenance group node object.
