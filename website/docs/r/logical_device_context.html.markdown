---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_logical_device_context"
sidebar_current: "docs-aci-resource-logical_device_context"
description: |-
  Manages ACI Logical Device Context
---

# aci_logical_device_context

Manages ACI Logical Device Context

## Example Usage

```hcl

resource "aci_logical_device_context" "example" {
  tenant_dn         = aci_tenant.tenentcheck.id
  ctrct_name_or_lbl = "default"
  graph_name_or_lbl = "any"
  node_name_or_lbl  = "any"
  annotation        = "example"
  description       = "from terraform"
  context           = "ctx1"
  name_alias        = "example"
  relation_vns_rs_l_dev_ctx_to_l_dev = "uni/tn-test_acc_tenant/lDevVip-LoadBalancer01"
}

```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `ctrct_name_or_lbl` - (Required) Ctrct name or label of Object logical device context.
- `graph_name_or_lbl` - (Required) Graph name or label of Object logical device context.
- `node_name_or_lbl` - (Required) Node name or label of Object logical device context.
- `relation_vns_rs_l_dev_ctx_to_l_dev` - (Required) Relation to either a service device cluster (vnsALDev) or to a proxy object for a logical device cluster (vnsLDevIf). Cardinality - N_TO_ONE. Type - String.
- `annotation` - (Optional) Annotation for object logical device context.
- `description` - (Optional) Description for object logical device context.
- `context` - (Optional) Context for object logical device context.
- `name_alias` - (Optional) Name alias for object logical device context.
- `relation_vns_rs_l_dev_ctx_to_rtr_cfg` - (Optional) Relation to class vnsRtrCfg. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Logical Device Context.

## Importing

An existing Logical Device Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_logical_device_context.example <Dn>
```
