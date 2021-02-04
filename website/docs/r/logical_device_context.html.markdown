---
layout: "aci"
page_title: "ACI: aci_logical_device_context"
sidebar_current: "docs-aci-resource-logical_device_context"
description: |-
  Manages ACI Logical Device Context
---

# aci_logical_device_context #
Manages ACI Logical Device Context

## Example Usage ##

```hcl

resource "aci_logical_device_context" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  ctrct_name_or_lbl  = "example"
  graph_name_or_lbl  = "example"
  node_name_or_lbl  = "example"
  annotation  = "example"
  context  = "example"
  name_alias  = "example"
}

```


## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `ctrct_name_or_lbl` - (Required) ctrct name or label of Object logical_device_context.
* `graph_name_or_lbl` - (Required) graph name or label of Object logical_device_context.
* `node_name_or_lbl` - (Required) node name or label of Object logical_device_context.
* `annotation` - (Optional) annotation for object logical_device_context.
* `context` - (Optional) context for object logical_device_context.
* `name_alias` - (Optional) name_alias for object logical_device_context.


* `relation_vns_rs_l_dev_ctx_to_l_dev` - (Optional) Relation to class vnsALDevIf. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vns_rs_l_dev_ctx_to_rtr_cfg` - (Optional) Relation to class vnsRtrCfg. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Logical Device Context.

## Importing ##

An existing Logical Device Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_logical_device_context.example <Dn>
```