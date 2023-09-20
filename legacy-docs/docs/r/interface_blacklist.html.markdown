---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_interface_blacklist"
sidebar_current: "docs-aci-resource-interface_blacklist"
description: |-
  Manages ACI Out of Service Fabric Path
---

# aci_interface_blacklist #

Manages ACI interface blacklist which is the equivalent of disabling an interface.

## API Information ##

* `Class` - fabricRsOosPath
* `Distinguished Name` - uni/fabric/outofsvc/rsoosPath-[{tDn}]

## GUI Information ##

* `Location` - Fabric > Inventory > Pod X > Leaf YYY > Interfaces > Physical Interfaces > eth1/Z

## Example Usage ##

```hcl
resource "aci_interface_blacklist" "example" {
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
* `annotation` - (Optional) Annotation of the blacklist object (fabricRsOosPath).


## Importing ##

An existing blacklist can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_interface_blacklist.example <Dn>
```