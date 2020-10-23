---
layout: "aci"
page_title: "ACI: aci_bridge_domain"
sidebar_current: "docs-aci-resource-bridge_domain"
description: |-
  Manages ACI Bridge Domain
---

# aci_bridge_domain #
Manages ACI Bridge Domain

## Example Usage ##

```hcl
	resource "aci_bridge_domain" "foobridge_domain" {
		tenant_dn                   = "${aci_tenant.tenant_for_bd.id}"
		description                 = "sample bridge domain"
		name                        = "demo_bd"
		optimize_wan_bandwidth      = "no"
		annotation                  = "tag_bd"
		arp_flood                   = "no"
		ep_clear                    = "no"
		ep_move_detect_mode         = "garp"
		host_based_routing          = "no"
		intersite_bum_traffic_allow = "yes"
		intersite_l2_stretch        = "yes"
		ip_learning                 = "yes"
		ipv6_mcast_allow            = "no"
		limit_ip_learn_to_subnets   = "yes"
		mac                         = "00:22:BD:F8:19:FF"
		mcast_allow                 = "yes"
		multi_dst_pkt_act           = "bd-flood"
		name_alias                  = "alias_bd"
		bridge_domain_type          = "regular"
		unicast_route               = "no"
		unk_mac_ucast_act           = "flood"
		unk_mcast_act               = "flood"
		vmac                        = "not-applicable"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object bridge_domain.
* `optimize_wan_bandwidth` - (Optional) Flag to enable OptimizeWanBandwidth between sites. Allowed values are "yes" and "no". Default is "no".
* `annotation` - (Optional) annotation for object bridge_domain.
* `arp_flood` - (Optional) A property to specify whether ARP flooding is enabled. If flooding is disabled, unicast routing will be performed on the target IP address. Allowed values are "yes" and "no". Default is "no".
* `ep_clear` - (Optional) Represents the parameter used by the node (i.e. Leaf) to clear all EPs in all leaves for this BD. Allowed values are "yes" and "no". Default is "no".
* `ep_move_detect_mode` - (Optional) The End Point move detection option uses the Gratuitous Address Resolution Protocol (GARP). A gratuitous ARP is an ARP broadcast-type of packet that is used to verify that no other device on the network has the same IP address as the sending device.
Allowed value: "garp"
* `host_based_routing` - (Optional) enables advertising host routes out of l3outs of this BD. Allowed values are "yes" and "no". Default is "no".
* `intersite_bum_traffic_allow` - (Optional)  Control whether BUM traffic is allowed between sites
.Allowed values are "yes" and "no". Default is "no".
* `intersite_l2_stretch` - (Optional) Flag to enable l2Stretch between sites. Allowed values are "yes" and "no". Default is "no".
* `ip_learning` - (Optional) Endpoint Dataplane Learning.Allowed values are "yes" and "no". Default is "yes".
* `ipv6_mcast_allow` - (Optional) Flag to indicate multicast IpV6 is allowed or not.Allowed values are "yes" and "no". Default is "no".
* `limit_ip_learn_to_subnets` - (Optional) Limits IP address learning to the bridge domain subnets only. Every BD can have multiple subnets associated with it. By default, all IPs are learned. Allowed values are "yes" and "no". Default is "yes".
* `ll_addr` - (Optional) override of system generated ipv6 link-local address.
* `mac` - (Optional) The MAC address of the bridge domain (BD) or switched virtual interface (SVI). Every BD by default takes the fabric-wide default MAC address. You can override that address with a different one. By default the BD will take a 00:22:BD:F8:19:FF mac address.
* `mcast_allow` - (Optional) Flag to indicate if multicast is enabled for IpV4 addresses. Allowed values are "yes" and "no". Default is "no".
* `multi_dst_pkt_act` - (Optional) The multiple destination forwarding method for L2 Multicast, Broadcast, and Link Layer traffic types. Allowed values are "bd-flood", "encap-flood" and "drop". Default is "bd-flood".
* `name_alias` - (Optional) name_alias for object bridge_domain.
* `bridge_domain_type` - (Optional) The specific type of the object or component. Allowed values are "regular" and "fc". Default is "regular".
* `unicast_route` - (Optional) The forwarding method based on predefined forwarding criteria (IP or MAC address). Allowed values are "yes" and "no". Default is "yes".
* `unk_mac_ucast_act` - (Optional) The forwarding method for unknown layer 2 destinations. Allowed values are "flood" and "proxy". Default is "proxy".
* `unk_mcast_act` - (Optional) The parameter used by the node (i.e. a leaf) for forwarding data for an unknown multicast destination. Allowed values are "flood" and "opt-flood". Default is "flood".
* `v6unk_mcast_act` - (Optional) v6unk_mcast_act for object bridge_domain.
* `vmac` - (Optional) Virtual MAC address of the BD/SVI. This is used when the BD is extended to multiple sites using l2 Outside. Only allowed values is "not-applicable".

* `relation_fv_rs_bd_to_profile` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_mldsn` - (Optional) Relation to class mldSnoopPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_abd_pol_mon_pol` - (Optional) Relation to class monEPGPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_to_nd_p` - (Optional) Relation to class ndIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_flood_to` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_fv_rs_bd_to_fhs` - (Optional) Relation to class fhsBDPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_to_relay_p` - (Optional) Relation to class dhcpRelayP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_ctx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_to_netflow_monitor_pol` - (Optional) Relation to class netflowMonitorPol. Cardinality - N_TO_M. Type - Set of Map.
                
* `relation_fv_rs_igmpsn` - (Optional) Relation to class igmpSnoopPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_to_ep_ret` - (Optional) Relation to class fvEpRetPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_fv_rs_bd_to_out` - (Optional) Relation to class l3extOut. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Bridge Domain.

## Importing ##

An existing Bridge Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bridge_domain.example <Dn>
```