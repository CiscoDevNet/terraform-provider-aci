---
layout: "aci"
page_title: "ACI: aci_maintenance_group_node"
sidebar_current: "docs-aci-resource-maintenance_group_node"
description: |-
  Manages ACI Maintenance Group Node
---

# aci_maintenance_group_node #
Manages ACI Maintenance Group Node

## Example Usage ##

```hcl
resource "aci_maintenance_group_node" "example" {
  pod_maintenance_group_dn = "${aci_pod_maintenance_group.example.id}"
  name                     = "First"
  annotation               = "example"
  from_                    = "1"
  name_alias               = "aliasing"
  to_                      = "5"
}
```


## Argument Reference ##

* `pod_maintenance_group_dn` - (Required) Distinguished name of parent POD maintenance group object.
* `name` - (Required) Name of maintenance group node object.
* `annotation` - (Optional) Annotation for maintenance group node object.
* `from_` - (Optional) From value for maintenance group node object.
* `name_alias` - (Optional) Name alias for maintenance group node object.
* `to_` - (Optional) To value for maintenance group node object.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Maintenance Group Node.

## Importing ##

An existing Maintenance Group Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_maintenance_group_node.example <Dn>
```