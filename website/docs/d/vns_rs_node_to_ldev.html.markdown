---
layout: "aci"
page_title: "ACI: aci_relationfroma_abs_nodetoan_l_dev"
sidebar_current: "docs-aci-data-source-relationfroma_abs_nodetoan_l_dev"
description: |-
  Data source for ACI Relation from a AbsNode to an LDev
---

# aci_relationfroma_abs_nodetoan_l_dev #
Data source for ACI Relation from a AbsNode to an LDev

## Example Usage ##

```hcl
data "aci_relationfroma_abs_nodetoan_l_dev" "example" {

  function_node_dn  = "${aci_function_node.example.id}"
}
```
## Argument Reference ##
* `function_node_dn` - (Required) Distinguished name of parent FunctionNode object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Relation from a AbsNode to an LDev.
* `annotation` - (Optional) annotation for object relationfroma_abs_nodetoan_l_dev.
* `t_dn` - (Optional) distinguished name of the target
