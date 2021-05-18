---
layout: "aci"
page_title: "ACI: aci_logical_interface_context"
sidebar_current: "docs-aci-resource-logical_interface_context"
description: |-
  Manages ACI Logical Interface Context.
---

# aci_logical_interface_context

Manages ACI Logical Interface Context

## Example Usage

```hcl
resource "aci_logical_interface_context" "example" {

  logical_device_context_dn  = "${aci_logical_device_context.example.id}"
  annotation  = "example"
  conn_name_or_lbl  = "example"
  l3_dest  = "no"
  name_alias  = "example"
  permit_log  = "no"
}
```

## Argument Reference

- `logical_device_context_dn` - (Required) Distinguished name of parent LogicalDeviceContext object.
- `conn_name_or_lbl` - (Required) The connector name or label for the logical interface context.
- `annotation` - (Optional) Annotation for object logical_interface_context.
- `l3_dest` - (Optional) Is this LIF a L3 Destination.
  Allowed values: "yes", "no". Default is "yes".
- `name_alias` - (Optional) Name alias for object logical_interface_context.
- `permit_log` - (Optional) Permit logging action for object logical_interface_context.
  Allowed values: "yes", "no". Default is "no".

- `relation_vns_rs_l_if_ctx_to_cust_qos_pol` - (Optional) Relation to class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_svc_e_pg_pol` - (Optional) Relation to class vnsSvcEPgPol. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_svc_redirect_pol` - (Optional) Relation to class vnsSvcRedirectPol. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_l_if` - (Optional) Relation to class vnsALDevLIf. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_out_def` - (Optional) Relation to class l3extOutDef. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_inst_p` - (Optional) Relation to class fvEPg. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
- `relation_vns_rs_l_if_ctx_to_out` - (Optional) Relation to class l3extOut. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Logical Interface Context.

## Importing

An existing Logical Interface Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_logical_interface_context.example <Dn>
```
