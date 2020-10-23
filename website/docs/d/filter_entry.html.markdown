---
layout: "aci"
page_title: "ACI: aci_filter_entry"
sidebar_current: "docs-aci-data-source-filter_entry"
description: |-
  Data source for ACI Filter Entry
---

# aci_filter_entry #
Data source for ACI Filter Entry

## Example Usage ##

```hcl
data "aci_filter_entry" "http" {
  filter_dn  = "${aci_filter.http_flt.id}"
  name       = "http"
}
```
## Argument Reference ##
* `filter_dn` - (Required) Distinguished name of parent Filter object.
* `name` - (Required) name of Object filter_entry.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Filter Entry.
* `annotation` - (Optional) annotation for object filter_entry.
* `apply_to_frag` - (Optional) Flag to determine whether to apply changes to fragment.
* `arp_opc` - (Optional) open peripheral codes.
* `d_from_port` - (Optional) Destination From Port.
* `d_to_port` - (Optional) Destination To Port.
* `ether_t` - (Optional) ether type for the entry.
* `icmpv4_t` - (Optional) ICMPv4 message type; used when ip_protocol is icmp.
* `icmpv6_t` - (Optional) ICMPv6 message type; used when ip_protocol is icmpv6.
* `match_dscp` - (Optional) The matching differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
* `name_alias` - (Optional) name_alias for object filter_entry.
* `prot` - (Optional) level 3 ip protocol.
* `s_from_port` - (Optional) Source From Port.
* `s_to_port` - (Optional) Source To Port.
* `stateful` - (Optional) Determines if entry is stateful or not.
* `tcp_rules` - (Optional) TCP Session Rules.
