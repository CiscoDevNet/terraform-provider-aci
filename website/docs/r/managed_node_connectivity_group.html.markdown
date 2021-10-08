---
layout: "aci"
page_title: "ACI: aci_managed_node_connectivity_group"
sidebar_current: "docs-aci-resource-managed_node_connectivity_group"
description: |-
  Manages ACI Managed Node Connectivity Group
---

# aci_managed_node_connectivity_group #

Manages ACI Managed Node Connectivity Group

## API Information ##
* `Class` - mgmtGrp
* `Distinguished Named` - uni/infra/funcprof/grp-{name}

## GUI Information ##
* `Location` - Tenants -> mgmt -> Managed Node Connectivity Groups -> Create Managed Node Connectivity Group


## Example Usage ##

```hcl
resource "aci_managed_node_connectivity_group" "example" {
  name  = "example"
  annotation = "test_annotation"
}
```

## Argument Reference ##
* `name` - (Required) Name of object Managed Node Connectivity Group.
* `annotation` - (Optional) Annotation of object Managed Node Connectivity Group.

## Importing ##

An existing Managed Node Connectivity Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_managed_node_connectivity_group.example <Dn>
```