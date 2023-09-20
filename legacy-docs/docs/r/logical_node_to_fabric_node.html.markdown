---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_node_to_fabric_node"
sidebar_current: "docs-aci-resource-logical_node_to_fabric_node"
description: |-
  Manages ACI Logical Node to Fabric Node
---

# aci_logical_node_to_fabric_node #
Manages ACI Logical Node to Fabric Node

## Example Usage ##

```hcl
resource "aci_logical_node_to_fabric_node" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  tdn               = "topology/pod-1/node-201"
  annotation        = "annotation"
  config_issues     = "none"
  rtr_id            = "10.0.1.1"
  rtr_id_loop_back  = "no"
}
```
## Argument Reference ##
* `logical_node_profile_dn` - (Required) Distinguished name of parent LogicalNodeProfile object.
* `tdn` - (Required) Tdn of Object Fabric Node.
* `annotation` - (Optional) Annotation for object Fabric Node.
* `config_issues` - (Optional) Configuration issues. Allowed values: "anchor-node-mismatch", "bd-profile-missmatch", "loopback-ip-missing", "missing-mpls-infra-l3out", "missing-rs-export-route-profile", "node-path-misconfig", "node-vlif-misconfig", "none", "routerid-not-changable-with-mcast", "subnet-mismatch". Default value: "none"
* `rtr_id` - (Optional) Router identifier
* `rtr_id_loop_back` - (Optional) Allowed values: "yes", "no". Default value: "yes"



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fabric Node.

## Importing ##

An existing Fabric Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_logical_node_to_fabric_node.example <Dn>
```