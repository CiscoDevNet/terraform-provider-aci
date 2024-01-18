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
  parent_dn = aci_bridge_domain.example.id
  ip        = "10.0.3.28/27"
}
```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of parent bridge domain object. Type: String
- `ip` - (Required) The IP address and mask of the default gateway. Type: String

## Attribute Reference

- `id` - (Read-Only) Attribute id set to the Dn of the Subnet. Type: String
- `annotation` - (Read-Only) Annotation for object subnet. Type: String
- `description` - (Read-Only) Description for object subnet. Type: String
- `ctrl` - (Read-Only) The list of subnet control state. The control can be specific protocols applied to the subnet such as IGMP Snooping. Type: List
- `ip_data_plane_learning` - (Read-Only) Flag to enable/disable ip-data-plane learning for the Subnet object. Allowed values are "enabled" and "disabled" and default value is "enabled". Type: String.
- `name_alias` - (Read-Only) Name alias for object subnet. Type: String
- `preferred` - (Read-Only) Indicates if the subnet is preferred (primary) over the available alternatives. Only one preferred subnet is allowed. Type: String
- `scope` - (Read-Only) The List of network visibility of the subnet. Type: List
- `virtual` - (Read-Only) Treated as virtual IP address. Used in case of BD extended to multiple sites. Type: String
- `relation_fv_rs_bd_subnet_to_out` - (Read-Only) Relation to class l3extOut. Cardinality - N_TO_M. Type: List
- `relation_fv_rs_nd_pfx_pol` - (Read-Only) Relation to class ndPfxPol. Cardinality - N_TO_ONE. Type: List.
- `relation_fv_rs_bd_subnet_to_profile` - (Read-Only) Relation to class rtctrlProfile. Cardinality - N_TO_ONE. Type: List.
- `next_hop_addr` - (Read-Only) EP Reachability of the Application EPGs Subnet object. Type: String.
- `msnlb` - (Read-Only) A block representing MSNLB of the Application EPGs Subnet object. Type: Block.
   - `mode` - (Read-Only) Mode of the MSNLB object. Type: String
   - `group` - (Read-Only) The IGMP mode group IP address of the MSNLB object. Type: String
   - `mac` - (Read-Only) MAC address of the unicast and static multicast mode of the MSNLB object. Type: String
- `anycast_mac` - (Read-Only) Anycast MAC of the Application EPGs Subnet object. Type - String. Type: String