---
subcategory: "Node Management"
layout: "aci"
page_title: "ACI: aci_managed_node_connectivity_group"
sidebar_current: "docs-aci-data-source-managed_node_connectivity_group"
description: |-
  Data source for ACI Managed Node Connectivity Group
---

# aci_managed_node_connectivity_group #

Data source for ACI Managed Node Connectivity Group

## API Information ##
* `Class` - mgmtGrp
* `Distinguished Name` - uni/infra/funcprof/grp-{name}

## GUI Information ##
* `Location` - Tenants -> mgmt -> Managed Node Connectivity Groups -> Create Managed Node Connectivity Group

## Example Usage ##

```hcl
data "aci_managed_node_connectivity_group" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) name of object Managed Node Connectivity Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Managed Node Connectivity Group.
* `annotation` - (Optional) Annotation of object Managed Node Connectivity Group.

