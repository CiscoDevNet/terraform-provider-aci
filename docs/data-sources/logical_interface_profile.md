---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_interface_profile"
sidebar_current: "docs-aci-data-source-aci_logical_interface_profile"
description: |-
  Data source for ACI Logical Interface Profile
---

# aci_logical_interface_profile

Data source for ACI Logical Interface Profile

## API Information

- `Class` - l3extLIfP
- `Distinguished Name` - uni/tn-{tenant_name}/out-{l3out}/lnodep-{logical_node_profile}/lifp-{logical_interface_profile}

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
- `annotation` - (Read-Only) Annotation of the logical interface profile object.
- `description` - (Read-Only) Description of the logical interface profile object.
- `name_alias` - (Read-Only) Name alias of the logical interface profile object.
- `prio` - (Read-Only) QoS priority class id.
- `tag` - (Read-Only) Specifies the color of a policy label.

- `relation_l3ext_rs_pim_ip_if_pol` - (Read-Only) Represents the relation to the PIM Interface Policy (class pimIfPol).
- `relation_l3ext_rs_pim_ipv6_if_pol` - (Read-Only) Represents the relation to the PIM IPv6 Interface Policy (class pimIfPol).
- `relation_l3ext_rs_igmp_if_pol` - (Read-Only) Represents the relation to the IGMP Interface Policy (class igmpIfPol).
- `relation_l3ext_rs_l_if_p_to_netflow_monitor_pol` - (Read-Only) Relation to the Netflow Monitor Policy (class netflowMonitorPol). Cardinality - N_TO_M. Type - Block.
  - `tn_netflow_monitor_pol_name` - (Deprecated) Distinguished name of the target Netflow Monitor Policy.
	- `tn_netflow_monitor_pol_dn` -  (Read-Only) Distinguished name of the target Netflow Monitor Policy.
	- `flt_type` - (Read-Only) Netflow IP filter type.
- `relation_l3ext_rs_egress_qos_dpp_pol` - (Read-Only) Relation to the Egress Data Plane Policing Policy (class qosDppPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_ingress_qos_dpp_pol` - (Read-Only) Relation to the Ingress Data Plane Policing Policy (class qosDppPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_l_if_p_cust_qos_pol` - (Read-Only) Relation to the Custom QoS Policy (class qosCustomPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_nd_if_pol` - (Read-Only) Relation to the IPv6 ND policy (class ndIfPol). Cardinality - N_TO_ONE. Type - String.
