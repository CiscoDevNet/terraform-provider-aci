---
layout: "aci"
page_title: "ACI: aci_fabric_node_member"
sidebar_current: "docs-aci-resource-fabric_node_member"
description: |-
  Manages ACI Fabric Node Member
---

# aci_fabric_node_member #
Manages ACI Fabric Node Member

## Example Usage ##

```hcl
resource "aci_fabric_node_member" "example" {

  name = "test"
  serial  = "example"
  annotation  = "example"
  ext_pool_id  = "example"
  fabric_id  = "example"
  name_alias  = "example"
  node_id  = "example"
  node_type  = "example"
  pod_id  = "example"
  role  = "example"
}
```
## Argument Reference ##
* `serial` - (Required) serial of Object fabric_node_member.
* `name` - (Required) Name of Fabric Node member.
* `annotation` - (Optional) annotation for object fabric_node_member.
* `ext_pool_id` - (Optional) ext_pool_id for object fabric_node_member.
* `fabric_id` - (Optional) place holder for a value
* `name_alias` - (Optional) name_alias for object fabric_node_member.
* `node_id` - (Optional) node id
* `node_type` - (Optional) node_type for object fabric_node_member.
* `pod_id` - (Optional) pod id
* `role` - (Optional) system role type
* `serial` - (Optional) serial number



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fabric Node Member.

## Importing ##

An existing Fabric Node Member can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_node_member.example <Dn>
```