---
layout: "aci"
page_title: "ACI: aci_spanning_tree_interface_policy"
sidebar_current: "docs-aci-data-source-spanning_tree_interface_policy"
description: |-
  Data source for ACI Spanning Tree Interface Policy
API Information:
 - Class: "stpIfPol"
 - Distinguished Named: "uni/infra/ifPol"
GUI Location:
 - Fabric > Access Policies > Policies > Interface > Spanning Tree Interface
---

# aci_spanning_tree_interface_policy #
Data source for ACI Spanning Tree Interface Policy

## Example Usage ##

```hcl
data "aci_spanning_tree_interface_policy" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object spanning_tree_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Spanning Tree Interface Policy.
* `annotation` - (Optional) annotation for object spanning_tree_interface_policy.
* `ctrl` - (Optional) stp interface control
* `name_alias` - (Optional) name_alias for object spanning_tree_interface_policy.
