---
layout: "aci"
page_title: "ACI: aci_logical_node_to_fabric_node"
sidebar_current: "docs-aci-data-source-logical_node_to_fabric_node"
description: |-
  Data source for ACI Logical Node to Fabric Node
---

# aci_logical_node_to_fabric_node #
Data source for ACI Logical Node to Fabric Node

## Example Usage ##

```hcl
data "aci_logical_node_to_fabric_node" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  tdn  = "topology/pod-1/paths-101/pathep-[eth1/4]"
}
```
## Argument Reference ##
* `logical_node_profile_dn` - (Required) Distinguished name of parent LogicalNodeProfile object.
* `tdn` - (Required) Tdn of Object Fabric Node.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Fabric Node.
* `annotation` - (Optional) Annotation for object Fabric Node.
* `config_issues` - (Optional) Configuration issues
* `rtr_id` - (Optional) Router identifier
* `rtr_id_loop_back` - (Optional) 
