---
layout: "aci"
page_title: "ACI: aci_spanning_tree_interface_policy"
sidebar_current: "docs-aci-resource-spanning_tree_interface_policy"
description: |-
  Manages ACI Spanning Tree Interface Policy
API Information:
 - Class: "stpIfPol"
 - Distinguished Named: "uni/infra/ifPol"
GUI Location:
 - Fabric > Access Policies > Policies > Interface > Spanning Tree Interface
---

# aci_spanning_tree_interface_policy #
Manages ACI Spanning Tree Interface Policy

## Example Usage ##

```hcl
resource "aci_spanning_tree_interface_policy" "example" {
  name  = "demo_stpifpol"
  annotation  = "tag_stp"
  ctrl  = ["%s"]
  name_alias  = "alias_stpifpol"
}
```
## Argument Reference ##
* `name` - (Required) name of Object spanning_tree_interface_policy.
* `annotation` - (Optional) annotation for object spanning_tree_interface_policy.
* `ctrl` - (Optional) stp interface control. Allowed values are "bpdu-filter", "bpdu-guard".
* `name_alias` - (Optional) name_alias for object spanning_tree_interface_policy.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spanning Tree Interface Policy.

## Importing ##

An existing Spanning Tree Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spanning_tree_interface_policy.example <Dn>
```