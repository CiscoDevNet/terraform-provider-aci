---
layout: "aci"
page_title: "ACI: aci_ospf_interface_policy"
sidebar_current: "docs-aci-resource-ospf_interface_policy"
description: |-
  Manages ACI OSPF Interface Policy
---

# aci_ospf_interface_policy #
Manages ACI OSPF Interface Policy

## Example Usage ##

```hcl
	resource "aci_ospf_interface_policy" "fooospf_interface_policy" {
		tenant_dn    = "${aci_tenant.dev_tenant.id}"
		description  = "from terraform"
		name         = "demo_ospfpol"
		annotation   = "tag_ospf"
		cost         = "unspecified"
		ctrl         = "unspecified"
		dead_intvl   = "40"
		hello_intvl  = "10"
		name_alias   = "alias_ospf"
		nw_t         = "unspecified"
		pfx_suppress = "inherit"
		prio         = "1"
		rexmit_intvl = "5"
		xmit_delay   = "1"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object ospf interface policy.
* `annotation` - (Optional) annotation for object ospf interface policy.
* `cost` - (Optional) The OSPF cost for the interface. The cost (also called metric) of an interface in OSPF is an indication of the overhead required to send packets across a certain interface. The cost of an interface is inversely proportional to the bandwidth of that interface. A higher bandwidth indicates a lower cost. There is more overhead (higher cost) and time delays involved in crossing a 56k serial line than crossing a 10M ethernet line. The formula used to calculate the cost is: cost= 10000 0000/bandwidth in bps For example, it will cost 10 EXP8/10 EXP7 = 10 to cross a 10M Ethernet line and will cost 10 EXP8/1544000 = 64 to cross a T1 line. By default, the cost of an interface is calculated based on the bandwidth; you can force the cost of an interface with the ip ospf cost value interface sub-configuration mode command. Allowed value range is "0" - "65535". Default is "unspecified(0)".
* `ctrl` - (Optional) Interface policy controls. Allowed values are "unspecified", "passive", "mtu-ignore", "advert-subnet" and "bfd". Default is "unspecified". 
* `dead_intvl` - (Optional) The interval between hello packets from a neighbor before the router declares the neighbor as down. This value must be the same for all networking devices on a specific network. Specifying a smaller dead interval (seconds) will give faster detection of a neighbor being down and improve convergence, but might cause more routing instability. Allowed value range is "1" - "65535". Default value is "40".
* `hello_intvl` - (Optional) The interval between hello packets that OSPF sends on the interface. Note that the smaller the hello interval, the faster topological changes will be detected, but more routing traffic will ensue. This value must be the same for all routers and access servers on a specific network. Allowed value range is "1" - "65535". Default value is "10".
* `name_alias` - (Optional) name_alias for object ospf_interface_policy.
* `nw_t` - (Optional) The OSPF interface policy network type. OSPF supports point-to-point and broadcast. Allowed values are "unspecified", "p2p" and "bcast". Default value is "unspecified".
* `pfx_suppress` - (Optional) pfx-suppression for object ospf_interface_policy. Allowed values are "inherit", "enable" and "disable". Default value is "inherit".
* `prio` - (Optional) The OSPF interface priority used to determine the designated router (DR) on a specific network. The router with the highest OSPF priority on a segment will become the DR for that segment. The same process is repeated for the backup designated router (BDR). In the case of a tie, the router with the highest RID will win. The default for the interface OSPF priority is one. Remember that the DR and BDR concepts are per multiaccess segment. Allowed value range is "0" - "255". Default value is "1".
* `rexmit_intvl` - (Optional) The interval between LSA retransmissions. The retransmit interval occurs while the router is waiting for an acknowledgement from the neighbor router that it received the LSA. If no acknowlegment is received at the end of the interval, then the LSA is resent. Allowed value range is "1" - "65535". Default value is "5".
* `xmit_delay` - (Optional) The delay time needed to send an LSA update packet. OSPF increments the LSA age time by the transmit delay amount before transmitting the LSA update. You should take into account the transmission and propagation delays for the interface when you set this value. Allowed value range is "1" - "450". Default is "1".


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the OSPF Interface Policy.

## Importing ##

An existing OSPF Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_ospf_interface_policy.example <Dn>
```