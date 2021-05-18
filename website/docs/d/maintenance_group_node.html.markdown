---
layout: "aci"
page_title: "ACI: aci_maintenance_group_node"
sidebar_current: "docs-aci-data-source-maintenance_group_node"
description: |-
  Data source for ACI Maintenance Group Node
---

# aci_maintenance_group_node #
Data source for ACI Maintenance Group Node

## Example Usage ##

```hcl
data "aci_maintenance_group_node" "example" {
  pod_maintenance_group_dn = "${aci_pod_maintenance_group.example.id}"
  name                     = "example"
}
```


## Argument Reference ##

* `pod_maintenance_group_dn` - (Required) Distinguished name of parent POD maintenance group object.
* `name` - (Required) Name of maintenance group node object.



## Attribute Reference

* `id` - Attribute id set to the dn of the maintenance group node object.
* `annotation` - (Optional) Annotation for maintenance group node object.
* `from_` - (Optional) From for maintenance group node object.
* `name_alias` - (Optional) Name alias for maintenance group node object.
* `to_` - (Optional) To for maintenance group node object.
