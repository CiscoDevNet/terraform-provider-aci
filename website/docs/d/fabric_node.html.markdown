---
layout: "aci"
page_title: "ACI: aci_fabric_node"
sidebar_current: "docs-aci-data-source-fabric_node"
description: |-
  Data source for ACI Fabric Node
---

# aci_fabric_node #
Data source for ACI Fabric Node

## Example Usage ##

```hcl
data "aci_fabric_node" "example" {
  fabric_pod_dn  = aci_fabric_pod.example.id
  fabric_node_id  = "example"
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
* `fabric_node_id` - (Optional) The identifier of object Fabric Node.
* `last_state_mod_ts` - (Optional) Last state mode time stamp for object Fabric Node.
* `name_alias` - (Optional) Name alias for object Fabric Node.
* `description` - (Optional) Description for object Fabric Node.
