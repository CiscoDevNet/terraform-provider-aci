---
layout: "aci"
page_title: "ACI: aci_ospf_interface_policy"
sidebar_current: "docs-aci-data-source-ospf_interface_policy"
description: |-
  Data source for ACI OSPF Interface Policy
---

# aci_ospf_interface_policy

Data source for ACI OSPF Interface Policy

## Example Usage

```hcl
data "aci_ospf_interface_policy" "dev_ospf_pol" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "foo_ospf_pol"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object ospf interface policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the OSPF Interface Policy.
- `annotation` - (Optional) Annotation for object ospf interface policy.
- `cost` - (Optional) The OSPF cost for the interface. The cost (also called metric) of an interface in OSPF is an indication of the overhead required to send packets across a certain interface. The cost of an interface is inversely proportional to the bandwidth of that interface. A higher bandwidth indicates a lower cost. There is more overhead (higher cost) and time delays involved in crossing a 56k serial line than crossing a 10M ethernet line. The formula used to calculate the cost is: cost= 10000 0000/bandwidth in bps For example, it will cost 10 EXP8/10 EXP7 = 10 to cross a 10M Ethernet line and will cost 10 EXP8/1544000 = 64 to cross a T1 line. By default, the cost of an interface is calculated based on the bandwidth; you can force the cost of an interface with the ip ospf cost value interface sub-configuration mode command.
- `ctrl` - (Optional) Interface policy controls
- `dead_intvl` - (Optional) The interval between hello packets from a neighbor before the router declares the neighbor as down. This value must be the same for all networking devices on a specific network. Specifying a smaller dead interval (seconds) will give faster detection of a neighbor being down and improve convergence, but might cause more routing instability.
- `hello_intvl` - (Optional) The interval between hello packets that OSPF sends on the interface. Note that the smaller the hello interval, the faster topological changes will be detected, but more routing traffic will ensue. This value must be the same for all routers and access servers on a specific network.
- `name_alias` - (Optional) Name alias for object ospf interface policy.
- `nw_t` - (Optional) The OSPF interface policy network type. OSPF supports point-to-point and broadcast.
- `pfx_suppress` - (Optional) pfx-suppression for object ospf interface policy.
- `prio` - (Optional) The OSPF interface priority used to determine the designated router (DR) on a specific network. The router with the highest OSPF priority on a segment will become the DR for that segment. The same process is repeated for the backup designated router (BDR). In the case of a tie, the router with the highest RID will win. The default for the interface OSPF priority is one. Remember that the DR and BDR concepts are per multiaccess segment.
- `rexmit_intvl` - (Optional) The interval between LSA retransmissions. The retransmit interval occurs while the router is waiting for an acknowledgement from the neighbor router that it received the LSA. If no acknowlegment is received at the end of the interval, then the LSA is resent.
- `xmit_delay` - (Optional) The delay time needed to send an LSA update packet. OSPF increments the LSA age time by the transmit delay amount before transmitting the LSA update. You should take into account the transmission and propagation delays for the interface when you set this value.
