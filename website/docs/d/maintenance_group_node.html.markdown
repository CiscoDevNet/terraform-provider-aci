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
  pod_maintenance_group_dn = aci_pod_maintenance_group.example.id
  description              = "from terraform"
  name                     = "example"
}
```


## Argument Reference ##

* `pod_maintenance_group_dn` - (Required) Distinguished name of parent POD Maintenance Group Object.
* `name` - (Required) Name of Maintenance Group Node Object.



## Attribute Reference

* `id` - Attribute id set to the dn of the Maintenance Group Node Object.
* `description` - (Optional) Description for Maintenance Group Node Object.
* `annotation` - (Optional) Annotation for Maintenance Group Node Object.
* `from_` - (Optional) From for Maintenance Group Node Object.
* `name_alias` - (Optional) Name alias for Maintenance Group Node Object.
* `to_` - (Optional) To for Maintenance Group Node Object.
