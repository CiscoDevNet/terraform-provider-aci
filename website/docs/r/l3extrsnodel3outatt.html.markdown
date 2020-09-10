---
layout: "aci"
page_title: "ACI: aci_fabric_node"
sidebar_current: "docs-aci-resource-fabric_node"
description: |-
  Manages ACI Fabric Node
---

# aci_logical_node_to_fabric_node #
Manages ACI Fabric Node

## Example Usage ##

```hcl
resource "aci_logical_node_to_fabric_node" "example" {

  logical_node_profile_dn  = "${aci_logical_node_profile.example.id}"

  tDn  = "example"
  annotation  = "example"
  config_issues  = "example"
  rtr_id  = "example"
  rtr_id_loop_back  = "example"
}
```
## Argument Reference ##
* `logical_node_profile_dn` - (Required) Distinguished name of parent LogicalNodeProfile object.
* `tDn` - (Required) tDn of Object fabric_node.
* `annotation` - (Optional) annotation for object fabric_node.
* `config_issues` - (Optional) configuration issues.
Allowed values: "none", "node-path-misconfig","routerid-not-changable-with-mcast", "loopback-ip-missing"
* `rtr_id` - (Optional) router identifier
* `rtr_id_loop_back` - (Optional) Allowed values: "yes", "no"



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fabric Node.

## Importing ##

An existing Fabric Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_node.example <Dn>
```