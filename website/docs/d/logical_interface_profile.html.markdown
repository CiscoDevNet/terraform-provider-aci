---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_interface_profile"
sidebar_current: "docs-aci-data-source-logical_interface_profile"
description: |-
  Data source for ACI Logical Interface Profile
---

# aci_logical_interface_profile

Data source for ACI Logical Interface Profile

## API Information

- `Class` - l3extLIfP
- `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/lnodep-{logical_node_profile}/lifp-{logical_interface_profile}

## GUI Information

- `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage

```hcl
data "aci_logical_interface_profile" "example" {
  logical_node_profile_dn  = aci_logical_node_profile.example.id
  name  = "example"
}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of the parent Logical Node Profile object.
- `name` - (Required) Name of the logical interface profile object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Logical Interface Profile.
- `annotation` - (Optional) Annotation of the logical interface profile object.
- `description` - (Optional) Description of the logical interface profile object.
- `name_alias` - (Optional) Name alias of the logical interface profile object.
- `prio` - (Optional) QoS priority class id.
- `tag` - (Optional) Specifies the color of a policy label.

- `relation_l3ext_rs_l_if_p_to_netflow_monitor_pol` - (Optional) Relation to the Netflow Monitor Policy (class netflowMonitorPol). Cardinality - N_TO_M. Type - Block.
  - `tn_netflow_monitor_pol_name` - (Deprecated) Distinguished name of the target Netflow Monitor Policy.
	- `tn_netflow_monitor_pol_dn` -  (Optional) Distinguished name of the target Netflow Monitor Policy.
	- `flt_type` - (Optional) Netflow IP filter type.
- `relation_l3ext_rs_egress_qos_dpp_pol` - (Optional) Relation to the class qosDppPol. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_ingress_qos_dpp_pol` - (Optional) Relation to the class qosDppPol. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_l_if_p_cust_qos_pol` - (Optional) Relation to the class qosCustomPol. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_arp_if_pol` - (Optional) Relation to the class arpIfPol. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_nd_if_pol` - (Optional) Relation to the class ndIfPol. Cardinality - N_TO_ONE. Type - String.
