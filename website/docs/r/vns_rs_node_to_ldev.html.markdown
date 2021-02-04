---
layout: "aci"
page_title: "ACI: aci_relationfroma_abs_nodetoan_l_dev"
sidebar_current: "docs-aci-resource-relationfroma_abs_nodetoan_l_dev"
description: |-
  Manages ACI Relation from a AbsNode to an LDev
---

# aci_relationfroma_abs_nodetoan_l_dev #
Manages ACI Relation from a AbsNode to an LDev

## Example Usage ##

```hcl
resource "aci_relationfroma_abs_nodetoan_l_dev" "example" {

  function_node_dn  = "${aci_function_node.example.id}"
  annotation  = "example"
  t_dn  = "example"
}
```
## Argument Reference ##
* `function_node_dn` - (Required) Distinguished name of parent FunctionNode object.
* `annotation` - (Optional) annotation for object relationfroma_abs_nodetoan_l_dev.
* `t_dn` - (Optional) distinguished name of the target



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Relation from a AbsNode to an LDev.

## Importing ##

An existing Relation from a AbsNode to an LDev can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_relationfroma_abs_nodetoan_l_dev.example <Dn>
```