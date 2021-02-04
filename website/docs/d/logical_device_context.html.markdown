---
layout: "aci"
page_title: "ACI: aci_logical_device_context"
sidebar_current: "docs-aci-data-source-logical_device_context"
description: |-
  Data source for ACI Logical Device Context
---

# aci_logical_device_context #
Data source for ACI Logical Device Context

## Example Usage ##

```hcl

data "aci_logical_device_context" "check" {
  tenant_dn         = "${aci_tenant.tenentcheck.id}"
  ctrct_name_or_lbl = "example"
  graph_name_or_lbl = "example"
  node_name_or_lbl  = "example"
}

```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `ctrct_name_or_lbl` - (Required) ctrct name or label of Object logical_device_context.
* `graph_name_or_lbl` - (Required) graph name or label of Object logical_device_context.
* `node_name_or_lbl` - (Required) node name or label of Object logical_device_context.



## Attribute Reference

* `id` - attribute id set to the Dn of the Logical Device Context.
* `annotation` - annotation for object logical_device_context.
* `context` - context for object logical_device_context.
* `name_alias` - name_alias for object logical_device_context.