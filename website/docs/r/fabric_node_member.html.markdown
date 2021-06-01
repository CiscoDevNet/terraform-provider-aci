---
layout: "aci"
page_title: "ACI: aci_fabric_node_member"
sidebar_current: "docs-aci-resource-fabric_node_member"
description: |-
  Manages ACI Fabric Node Member
---

# aci_fabric_node_member

Manages ACI Fabric Node Member

## Example Usage

```hcl
resource "aci_fabric_node_member" "example" {
  name        = "test"
  serial      = "1"
  annotation  = "example"
  description = "from terraform"
  ext_pool_id = "0"
  fabric_id   = "1"
  name_alias  = "example"
  node_id     = "201"
  node_type   = "unspecified"
  pod_id      = "1"
  role        = "unspecified"
}
```

## Argument Reference

- `serial` - (Required) Serial Number for the new Fabric Node Member.
- `name` - (Required) Name of Fabric Node member.
- `annotation` - (Optional) Specifies the annotation of a Fabric Node member.
- `description` - (Optional) Specifies the description of a Fabric Node member.
- `ext_pool_id` - (Optional) external pool ID for object Fabric Node member. Default value: "0".
- `fabric_id` - (Optional) Fabric ID for the new Fabric Node Member. Default value: "1".
- `name_alias` - (Optional) Name alias for object Fabric Node member.
- `node_id` - (Optional) Node ID Number for the new Fabric Node Member. Default value: "0".
- `node_type` - (Optional) Node type for object Fabric Node member.
  Allowed values: "unspecified", "remote-leaf-wan". Default value: "unspecified".
- `pod_id` - (Optional) The pod id of the new Fabric Node Member. Default value: "1".
- `role` - (Optional) Role for the new Fabric Node Member. 
  Allowed values: "unspecified", "leaf", "spine". Default value: "unspecified".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Fabric Node Member.

## Importing

An existing Fabric Node Member can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_fabric_node_member.example <Dn>
```
