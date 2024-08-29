---
subcategory: "Fabric Inventory"
layout: "aci"
page_title: "ACI: aci_system"
sidebar_current: "docs-aci-data-source-aci_system"
description: |-
  Data source for ACI System
---

# aci_system #

Data source for ACI System

## API Information ##

* Class: [topSystem](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/topSystem/overview)

* Supported in ACI versions: 1.0(1e) and later.

* Distinguished Name Formats:
  - `topology/pod-{id}/node-{id}/sys`

## GUI Information ##

* Location: `Generic`

## Example Usage ##

```hcl

data "aci_system" "example" {
  pod_id    = "1"
  system_id = "1"
}

```

## Schema ##

### Required ###

* `system_id` (id) - (string) The identifier of the system object.
* `pod_id` (podId) - (string) The pod identifier.

### Read-Only ###

* `id` - (string) The distinguished name (DN) of the system object.
* `address` (address) - (string) The IP address of the system.
* `bootstrap_state` (bootstrapState) - (string) Bootstrap state of this system.
* `cluster_time_diff` (clusterTimeDiff) - (string) Difference in cluster time from local time for this system.
* `control_plane_mtu` (controlPlaneMTU) - (string) MTU for control plane (SUP-originated) packets.
* `current_time` (currentTime) - (string) The current time on this system.
* `enforce_subnet_check` (enforceSubnetCheck) - (string) Enforce subnet check on all VRFs.
* `external_tep_address` (etepAddr) - (string) The external TEP IP address of this pod.
* `fabric_domain` (fabricDomain) - (string) Fabric domain of this node.
* `fabric_id` (fabricId) - (string) The latest system health score.
* `fabric_mac` (fabricMAC) - (string) The MAC address of the fabric.
* `inband_management_address` (inbMgmtAddr) - (string) The In-band management IPv4 address.
* `inband_management_address_ipv6` (inbMgmtAddr6) - (string) The In-band management IPv6 address.
* `inband_management_address_mask_ipv6` (inbMgmtAddr6Mask) - (string) The In-band management IPv6 address subnet mask.
* `inband_management_address_mask` (inbMgmtAddrMask) - (string) The In-band management IPv4 address subnet mask.
* `inband_management_gateway` (inbMgmtGateway) - (string) The In-band management IPv4 gateway address.
* `inband_management_gateway_ipv6` (inbMgmtGateway6) - (string) The In-band management IPv6 gateway address.
* `last_reboot_time` (lastRebootTime) - (string) The last reboot time for this system.
* `last_reset_reason` (lastResetReason) - (string) The last reset reason for this system.
* `mod_ts` (modTs) - (string) The time when this object was last modified.
* `mode` (mode) - (string) Specifies if this system is configured in standalone mode or HA pair.
* `monitoring_policy_dn` (monPolDn) - (string) The monitoring policy attached to this observable object.
* `name` (name) - (string) The name of the system object.
* `name_alias` (nameAlias) - (string) The name alias of the system object.
* `node_type` (nodeType) - (string) The role of this system.
* `out_of_band_management_address` (oobMgmtAddr) - (string) The Out-of-band management IPv4 address.
* `out_of_band_management_address_ipv6` (oobMgmtAddr6) - (string) The Out-of-band management IPv6 address.
* `out_of_band_management_address_mask_ipv6` (oobMgmtAddr6Mask) - (string) The Out-of-band management IPv6 address subnet mask.
* `out_of_band_management_address_mask` (oobMgmtAddrMask) - (string) The Out-of-band management IPv4 address subnet mask.
* `out_of_band_management_gateway` (oobMgmtGateway) - (string) The Out-of-band management IPv4 gateway address.
* `out_of_band_management_gateway_ipv6` (oobMgmtGateway6) - (string) The Out-of-band management IPv6 gateway address.
* `remote_network_id` (remoteNetworkId) - (string) The remote network ID.
* `remote_node` (remoteNode) - (string) The remote system.
* `remote_leaf_auto_mode` (rlAutoMode) - (string) The remote leaf auto mode.
* `remote_leaf_group_id` (rlGroupId) - (string) The remote leaf site identifier.
* `remote_leaf_operational_pod_id` (rlOperPodId) - (string) The operational pod identifier for RL pod redundancy.
* `remote_leaf_routable_mode` (rlRoutableMode) - (string) Indicates whether the remote leaf is using a routable TEP IP address.
* `remote_leaf_direct_mode` (rldirectMode) - (string) The remote leaf direct mode.
* `role` (role) - (string) The system role type.
* `serial` (serial) - (string) The serial number of the system.
* `server_type` (serverType) - (string) DHCP server type.
* `site_id` (siteId) - (string) The site identifier.
* `state` (state) - (string) The operational state of this system.
* `system_uptime` (systemUpTime) - (string) The time (in seconds) since the system was booted.
* `tep_pool` (tepPool) - (string) The pool of TEP IP addresses allocated for this system.
* `unicast_xr_endpoint_learn_disable` (unicastXrEpLearnDisable) - (string) Indicates whether the learning of unicast endpoints for external routing is disabled.
* `version` (version) - (string) The version for this system.
* `virtual_mode` (virtualMode) - (string) The virtual mode of this system.
