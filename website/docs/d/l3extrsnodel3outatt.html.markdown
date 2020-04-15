---
layout: "aci"
page_title: "ACI: aci_fabric_node"
sidebar_current: "docs-aci-data-source-fabric_node"
description: |-
  Data source for ACI Fabric Node
---

# aci_logical_node_to_fabric_node #
Data source for ACI Fabric Node

## Example Usage ##

```hcl
data "aci_logical_node_to_fabric_node" "example" {

  logical_node_profile_dn  = "${aci_logical_node_profile.example.id}"

  tDn  = "example"
}
```
## Argument Reference ##
* `logical_node_profile_dn` - (Required) Distinguished name of parent LogicalNodeProfile object.
* `tDn` - (Required) tDn of Object fabric_node.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Fabric Node.
* `annotation` - (Optional) annotation for object fabric_node.
* `config_issues` - (Optional) configuration issues
* `rtr_id` - (Optional) router identifier
* `rtr_id_loop_back` - (Optional) 
