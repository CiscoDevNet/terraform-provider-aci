---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_maintenance_group_node"
sidebar_current: "docs-aci-resource-aci_maintenance_group_node"
description: |-
  Manages ACI Maintenance Group Node
---

# aci_maintenance_group_node #
Manages ACI Maintenance Group Node

## Example Usage ##

```hcl
resource "aci_maintenance_group_node" "example" {
  pod_maintenance_group_dn = aci_pod_maintenance_group.example.id
  description              = "from terraform"
  name                     = "First"
  annotation               = "example"
  from_                    = "1"
  name_alias               = "aliasing"
  to_                      = "5"
}
```


## Argument Reference ##

* `pod_maintenance_group_dn` - (Required) Distinguished name of parent POD Maintenance Group Object.
* `name` - (Required) Name of Maintenance Group Node Object.
* `description` - (Optional) Description for Maintenance Group Node Object.
* `annotation` - (Optional) Annotation for Maintenance Group Node Object.
* `from_` - (Optional) From value for Maintenance Group Node Object. Range : 1 - 16000. DefaultValue : "1"
* `name_alias` - (Optional) Name alias for Maintenance Group Node Object.
* `to_` - (Optional) To value for Maintenance Group Node Object. Range : 1 - 16000. DefaultValue : "1"



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Maintenance Group Node.

## Importing ##

An existing Maintenance Group Node can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_maintenance_group_node.example <Dn>
```