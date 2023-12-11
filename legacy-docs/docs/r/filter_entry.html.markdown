---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_filter_entry"
sidebar_current: "docs-aci-resource-aci_filter_entry"
description: |-
  Manages ACI Filter Entry
---

# aci_filter_entry #
Manages ACI Filter Entry

## Example Usage ##

```hcl
	resource "aci_filter_entry" "foofilter_entry" {
		filter_dn     = aci_filter.example.id
		description   = "From Terraform"
		name          = "demo_entry"
		annotation    = "tag_entry"
		apply_to_frag = "no"
		arp_opc       = "unspecified"
		d_from_port   = "unspecified"
		d_to_port     = "unspecified"
		ether_t       = "ipv4"
		icmpv4_t      = "unspecified"
		icmpv6_t      = "unspecified"
		match_dscp    = "CS0"
		name_alias    = "alias_entry"
		prot          = "tcp"
		s_from_port   = "0"
		s_to_port     = "0"
		stateful      = "no"
		tcp_rules     = ["ack","rst"]
	}
```
## Argument Reference ##
* `filter_dn` - (Required) Distinguished name of parent Filter object.
* `name` - (Required) Name of Object filter entry.
* `description` - (Optional) Description for object filter entry.
* `annotation` - (Optional) Annotation for object filter entry.
* `apply_to_frag` - (Optional) Flag to determine whether to apply changes to fragment. Allowed values are "yes" and "no". Default is "no". 
* `arp_opc` - (Optional) Open peripheral codes. Allowed values are "unspecified", "req" and "reply". Default is "unspecified".
* `d_from_port` - (Optional) Destination From Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `d_to_port` - (Optional) Destination To Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `ether_t` - (Optional) Ether type for the entry. Allowed values are "unspecified", "ipv4", "trill", "arp", "ipv6", "mpls_ucast", "mac_security", "fcoe" and "ip". Default is "unspecified".
* `icmpv4_t` - (Optional) ICMPv4 message type; used when ip_protocol is icmp. Allowed values are "echo-rep", "dst-unreach", "src-quench", "echo", "time-exceeded" and "unspecified". Default is "unspecified".
* `icmpv6_t` - (Optional) ICMPv6 message type; used when ip_protocol is icmpv6. Allowed values are "unspecified", "dst-unreach", "time-exceeded", "echo-req", "echo-rep", "nbr-solicit", "nbr-advert" and "redirect". Default is "unspecified".
* `match_dscp` - (Optional) The matching differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".
* `name_alias` - (Optional) Name alias for object filter entry.
* `prot` - (Optional) Level 3 IP protocol. Allowed values are "unspecified", "icmp", "igmp", "tcp", "egp", "igp", "udp", "icmpv6", "eigrp", "ospfigp", "pim", "l2tp" and a number between "0" and "255". Default is "unspecified".
* `s_from_port` - (Optional) Source From Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `s_to_port` - (Optional) Source To Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `stateful` - (Optional) Determines if entry is stateful or not. Allowed values are "yes" and "no". Default is "no".
* `tcp_rules` - (Optional) TCP Rules. TCP Session Rules Allowed values are "ack", "est", "fin", "rst", "syn", "unspecified", and default value is "unspecified". Type: List.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Filter Entry.

## Importing ##

An existing Filter Entry can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_filter_entry.example <Dn>
```