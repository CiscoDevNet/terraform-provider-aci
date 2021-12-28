---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract"
sidebar_current: "docs-aci-resource-contract"
description: |-
  Manages ACI Contract
---

# aci_contract #
Manages ACI Contract

## Example Usage ##

```hcl
	resource "aci_contract" "foocontract" {
		tenant_dn   =  aci_tenant.dev_tenant.id
		description = "From Terraform"
		name        = "demo_contract"
		annotation  = "tag_contract"
		name_alias  = "alias_contract"
		prio        = "level1"
		scope       = "tenant"
		target_dscp = "unspecified"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object contract.
* `description` - (Optional) Description for object contract.
* `annotation` - (Optional) Annotation for object contract.
* `name_alias` - (Optional) Name alias for object contract.
* `prio` - (Optional) Priority level of the service contract.  Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified".
* `scope` - (Optional)  Represents the scope of this contract. If the scope is set as application-profile, the epg can only communicate with epgs in the same application-profile. Allowed values are "global", "tenant", "application-profile" and "context". Default is "context".

* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".

* `filter` - (Optional) To manage filters from the contract resource. It has the attributes like filter_name, annotation, description and name_alias.
* `filter.filter_name` - (Required) Name of the filter object.
* `filter.description` - (Optional) Description for the filter object.
* `filter.annotation` - (Optional) Annotation for filter object.
* `filter.name_alias` - (Optional) Name alias for filter object.

* `filter.filter_entry` - (Optional) To manage filter entries for particular filter from the contract resource. It has following attributes.
* `filter.filter_entry.filter_entry_name` - (Required) Name of Object filter entry.
* `filter.filter_entry.entry_annotation` - (Optional) Annotation for object filter entry.
* `filter.filter_entry.entry_description` - (Optional) Description for the filter entry.
* `filter.filter_entry.apply_to_frag` - (Optional) Flag to determine whether to apply changes to fragment. Allowed values are "yes" and "no". Default is "no".
* `filter.filter_entry.arp_opc` - (Optional) Open peripheral codes. Allowed values are "unspecified", "req" and "reply". Default is "unspecified".
* `filter.filter_entry.d_from_port` - (Optional) Destination From Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `filter.filter_entry.d_to_port` - (Optional) Destination To Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `filter.filter_entry.ether_t` - (Optional) Ether type for the entry. Allowed values are "unspecified", "ipv4", "trill", "arp", "ipv6", "mpls_ucast", "mac_security", "fcoe" and "ip". Default is "unspecified".
* `filter.filter_entry.icmpv4_t` - (Optional) ICMPv4 message type; used when ip_protocol is icmp. Allowed values are "echo-rep", "dst-unreach", "src-quench", "echo", "time-exceeded" and "unspecified". Default is "unspecified".
* `filter.filter_entry.icmpv6_t` - (Optional) ICMPv6 message type; used when ip_protocol is icmpv6. Allowed values are "unspecified", "dst-unreach", "time-exceeded", "echo-req", "echo-rep", "nbr-solicit", "nbr-advert" and "redirect". Default is "unspecified".
* `filter.filter_entry.match_dscp` - (Optional) The matching differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11",	"AF12",	"AF13",	"CS2",	"AF21",	"AF22",	"AF23",	"CS3",	"AF31",	"AF32",	"AF33",	"CS4",	"AF41",	"AF42",	"AF43",	"CS5",	"VA",	"EF",	"CS6",	"CS7"	and "unspecified". Default is "unspecified".
* `filter.filter_entry.entry_name_alias` - (Optional) Name alias for object filter entry.
* `filter.filter_entry.prot` - (Optional) Level 3 ip protocol. Allowed values are "unspecified", "icmp", "igmp", "tcp", "egp", "igp", "udp", "icmpv6", "eigrp", "ospfigp", "pim" and "l2tp". Default is "unspecified".
* `filter.filter_entry.s_from_port` - (Optional) Source From Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `filter.filter_entry.s_to_port` - (Optional) Source To Port. Accepted values are any valid TCP/UDP port range. Default is "unspecified".
Allowed values: "unspecified", "ftpData", "smtp", "dns", "http","pop3", "https", "rtsp"
* `filter.filter_entry.stateful` - (Optional) Determines if entry is stateful or not. Allowed values are "yes" and "no". Default is "no".
* `filter.filter_entry.tcp_rules` - (Optional) TCP Session Rules. Allowed values are "unspecified", "est", "syn", "ack", "fin" and "rst". Default is "unspecified".
* `relation_vz_rs_graph_att` - (Optional) Relation to class vzRsGraphAtt. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Contract.
* `filter.id` - Exports this attribute for filter object. Set to the Dn for the filter managed by the contract.
* `filter.filter_entry.id` - Exports this attribute for filter entry object of filter object. Set to the Dn for the filter entry managed by the contract.

## Importing ##

An existing Contract can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_contract.example <Dn>
```
