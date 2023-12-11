---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_interface_blacklist"
sidebar_current: "docs-aci-data-source-aci_interface_blacklist"
description: |-
  Data source for ACI Out of Service Fabric Path
---

# aci_outof_service_fabric_path #

Data source for ACI interface blacklist which is the equivalent of a disabled interface.

## API Information ##

* `Class` - fabricRsOosPath
* `Distinguished Name` - uni/fabric/outofsvc/rsoosPath-[{tDn}]

## GUI Information ##

* `Location` - Fabric > Inventory > Pod X > Leaf YYY > Interfaces > Physical Interfaces > ethX/Z

## Example Usage ##

```hcl
data "aci_outof_service_fabric_path" "example" {
  pod_id  = 1
  node_id = 101
  interface = eth1/1
}
```

## Argument Reference ##

* `pod_id` - (Required) The Pod ID of the switch that own the interface that need to be disabled.
* `node_id` - (Required) The Node ID of the switch that own the interface that need to be disabled.
* `fex_id` - (Required) The FEX ID of the FEX that own the interface that need to be disabled.
* `interface` - (Required) The interface name of the interface that need to be disabled.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Out of Service Fabric Path.
* `annotation` - (Optional) Annotation of object Out of Service Fabric Path.
