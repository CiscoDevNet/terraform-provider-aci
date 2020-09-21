---
layout: "aci"
page_title: "ACI: aci_fabric_path_ep"
sidebar_current: "docs-aci-data-source-fabric_path_ep"
description: |-
  Data source for ACI Fabric Path End point
---

# aci_fabric_path_ep #
Data source for ACI Fabric Path End point

## Example Usage ##

```hcl

data "aci_fabric_path_ep" "example" {
  pod_id  = "1"
  node_id = "101"
  name    = "eth1/7"
}

```

## Argument Reference ##
* `pod_id` - (Required) pod ID for Object fabric path endpoint.
* `node_id` - (Required) node ID for Object fabric path endpoint.
* `name` - (Required) name of Object fabric path endpoint.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Fabric Path End-point.
