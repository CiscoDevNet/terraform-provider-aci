---
layout: "aci"
page_title: "ACI: aci_fabric_node"
sidebar_current: "docs-aci-data-source-fabric_node"
description: |-
  Data source for ACI Fabric Node
---

# aci_fabric_node #
Data source for ACI Fabric Node

## API Information ##
* `Class` - fabricNode
* `Distinguished Named` - topology/pod-{id}/node-{id}

## GUI Information ##
* `Location` - Go to help & tool > object store browser > enter fabricNode in Class or DN or URL field then click Run Query

## Example Usage ##

```hcl
data "aci_fabric_node" "example" {
  fabric_pod_dn  = "topology/pod-1"
  fabric_node_id  = "101"
}
```

## Argument Reference ##
* `fabric_pod_dn` - (Required) Distinguished name of parent Fabric Pod object.
* `fabric_node_id` - (Required) fabric_node_id of object Fabric Node.

## Attribute Reference
* `id` - Attribute id set to the Dn of the Fabric Node.
* `ad_st` - (Optional) The administrative state of object Fabric Node.
* `annotation` - (Optional) Annotation for object Fabric Node.
* `apic_type` - (Optional) The APIC type for object Fabric Node.
* `fabric_st` - (Optional) Fabric state for object Fabric Node.
* `address` - (Optional) IP address of object Fabric Node.
* `name` - (Optional) Name of object Fabric Node.
* `node_type` - (Optional) Fabric Node type of object Fabric Node.
* `role` - (Optional) Fabric Node role of object Fabric Node.
* `name_alias` - (Optional) Name alias for object Fabric Node.

