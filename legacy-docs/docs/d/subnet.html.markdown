---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_subnet"
sidebar_current: "docs-aci-data-source-aci_subnet"
description: |-
  Data source for ACI Subnet
---

# aci_subnet

Data source for ACI Subnet

## API Information
Class - fvSubnet
- Distinguished Name - uni/tn-{tenant_name}/BD-{bd_name}/subnet-[{subnet_ip}]
- Distinguished Name - uni/tn-{tenant_name}/ap-{ap_name}/epg-{epg_name}/subnet-[{subnet_ip}]

## GUI Information
- Location - Tenant > Networking > Bridge Domains > Subnets
- Location - Tenant > Application Profiles > Application EPGs > Subnets

## Example Usage

```hcl
data "aci_subnet" "dev_subnet" {
  parent_dn         = aci_bridge_domain.example.id
  ip                = "10.0.3.28/27"
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of parent bridge domain object.
- `ip` - (Required) The IP address and mask of the default gateway.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Subnet.
- `annotation` - (Read-Only) Annotation for object subnet.
- `description` - (Read-Only) Description for object subnet.
- `ctrl` - (Read-Only) The list of subnet control state. The control can be specific protocols applied to the subnet such as IGMP Snooping.
- `ip_data_plane_learning` - (Read-Only) Flag to enable/disable ip-data-plane learning for the Subnet object. Allowed values are "enabled" and "disabled" and default value is "enabled". Type: String.
- `name_alias` - (Read-Only) Name alias for object subnet.
- `preferred` - (Read-Only) Indicates if the subnet is preferred (primary) over the available alternatives. Only one preferred subnet is allowed.
- `scope` - (Read-Only) The List of network visibility of the subnet.
- `virtual` - (Read-Only) Treated as virtual IP address. Used in case of BD extended to multiple sites.
- `relation_fv_rs_bd_subnet_to_out` - (Read-Only) Relation to class l3extOut. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_nd_pfx_pol` - (Read-Only) Relation to class ndPfxPol. Cardinality - N_TO_ONE. Type - String.
- `relation_fv_rs_bd_subnet_to_profile` - (Read-Only) Relation to class rtctrlProfile. Cardinality - N_TO_ONE. Type - String.
- `next_hop_addr` - (Read-Only) EP Reachability of the Application EPGs Subnet object. Type - String.
- `msnlb` - (Read-Only) A block representing MSNLB of the Application EPGs Subnet object. Type - Block.
   - `mode` - Mode of the MSNLB object.
   - `group` - The IGMP mode group IP address of the MSNLB object.
   - `mac` - MAC address of the unicast and static multicast mode of the MSNLB object.
- `anycast_mac` - Anycast MAC of the Application EPGs Subnet object. Type - String.