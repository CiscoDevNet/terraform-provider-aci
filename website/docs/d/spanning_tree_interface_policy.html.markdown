---
layout: "aci"
page_title: "ACI: aci_spanning_tree_interface_policy"
sidebar_current: "docs-aci-data-source-spanning_tree_interface_policy"
description: |-
  Data source for ACI Spanning Tree Interface Policy
---

# aci_spanning_tree_interface_policy #

Data source for ACI Spanning Tree Interface Policy

## API Information ##

* `Class` - stpIfPol
* `Distinguished Named` - uni/infra/ifPol-{name}

## GUI Information ##

* `Location` - Fabric > Access Policies > Policies > Interface > Spanning Tree Interface

## Example Usage ##

```hcl
data "aci_spanning_tree_interface_policy" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) name of object Spanning Tree Interface Policy.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Spanning Tree Interface Policy.
* `annotation` - (Optional) Annotation of object Spanning Tree Interface Policy.
* `name_alias` - (Optional) Name Alias of object Spanning Tree Interface Policy.
* `description` - (Optional) Description of object Spanning Tree Interface Policy.
* `ctrl` - (Optional) Interface controls.
