---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_system"
sidebar_current: "docs-aci-data-source-aci_system"
description: |-
  Data source for ACI System
---

# aci_system

Data source for ACI System

## Example Usage

```hcl
data "aci_system" "example" {
  pod_id = "1"
  system_id = "1"
}
```

## Argument Reference

- `system_id` - (Required) Unique system ID.
- `pod_id` - (Required) POD Idenitfier.

## Attribute Reference

- `id` - Attribute id set to the Dn of the system.
- `address` - (Optional) The IP address of the system.
- `etep_addr` - (Optional) The External TEP IP address of this PoD.
- `name_alias` - (Optional) Name alias of the system.
- `node_type`- (Optional) Role of this system.
- `remote_network_id`- (Optional) Remote Network ID.
- `remote_node` - (Optional) Remote system.
- `rldirect_mode` - (Optional) Remote Leaf Direct Mode.
- `role` - (Optional) The system role type.
- `server_type` - (Optional) Type of server.
- `bootstrap_state` - (Optional) Bootstrap state of this system.
- `child_action` - (Optional) Delete or ignore. For internal use only.
- `config_issues` - (Optional) Bitmask representation of the configuration issues found during the endpoint group deployment.
- `control_plane_mtu` - (Optional) MTU for control plane (SUP-originated) packets.
- `current_time` - (Optional) The current time on this system.
- `enforce_subnet_check` - (Optional) Enforce subnet check on all VRFs.
- `fabric_domain` - (Optional) Fabric domain of this node.
- `fabric_id` - (Optional) The latest system health score. Use the navigation bar at the top right of the table to select which health level to view.
- `fabric_mac` - (Optional) MAC address of fabric.
- `inb_mgmt_addr` - (Optional) The in-band management IPv4 address.
- `inb_mgmt_addr6` - (Optional) In-band management IPv6 address.
- `inb_mgmt_addr6_mask` - (Optional) In-band management IPv6 address Mask.
- `inb_mgmt_addr_mask` - (Optional) In-band management IP address Mask.
- `inb_mgmt_addr_gateway` - (Optional) In-band management IP Gateway.
- `inb_mgmt_addr_gateway6` - (Optional) In-band management IPv6 Gateway.
- `last_reset_reason`- (Optional) Last reset reason for this system.
- `lc_own` - (Optional) A value that indicates how this object was created. For internal use only.
- `mod_ts` - (Optional) The time when this object was last modified.
- `mode` - (Optional) Specifies if this system is configured in standalone mode or HA pair.
- `mon_pol_dn` - (Optional) The monitoring policy attached to this observable object.
- `name` - (Optional) The system name.
- `oob_mgmt_addr` - (Optional) The out-of-band management IPv4 address.
- `oob_mgmt_addr6` - (Optional) The out-of-band management IPv6 address.
- `oob_mgmt_addr6_mask` - (Optional) The out-of-band management IPv6 address Mask.
- `oob_mgmt_addr_mask` - (Optional) The out-of-band management IP address Mask.
- `oob_mgmt_gateway` - (Optional) Out-of-band management IP Gateway.
- `oob_mgmt_gateway6` - (Optional) Out-of-band management IPv6 Gateway
- `rl_oper_pod_id` - (Optional) Operational POD Idenitfier for RL Pod Redundancy.
- `rl_routable_mode` -  (Optional) Is Remote-Leaf Routable.
- `serial` - (Optional) Serial Number of the system.
- `state` - (Optional) Operational state of this system.
- `system_uptime` - (Optional) The time (in seconds) since the system was booted.
- `tep_pool` - (Optional) Tep-Pool for this system
- `unicast_xr_ep_learn_disable` - (Optional) Disable xrLeanrs.
- `version` - (Optional) The version of the compatibility catalog.
- `virtual_mode` - (Optional) Virtual mode of system.
