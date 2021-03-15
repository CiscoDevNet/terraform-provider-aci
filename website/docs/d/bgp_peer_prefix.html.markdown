---
layout: "aci"
page_title: "ACI: aci_bgp_peer_prefix"
sidebar_current: "docs-aci-data-source-bgp_peer_prefix"
description: |-
  Data source for ACI BGP Peer Prefix
---

# aci_bgp_peer_prefix #
Data source for ACI BGP Peer Prefix

## Example Usage ##

```hcl
data "aci_bgp_peer_prefix" "example" {
  tenant_dn = "${aci_tenant.tenentcheck.id}"
  name      = "one"
}
```


## Argument Reference ##

* `tenant_dn` - (Required) distinguished name of parent tenant object.
* `name` - (Required) name of BGP peer prefix object.



## Attribute Reference

* `id` - attribute id set to the Dn of BGP peer prefix object.
* `action` - action when the maximum prefix limit is reached for BGP peer prefix object.
* `description` - description for BGP peer prefix object.
* `annotation` - annotation for BGP peer prefix object.
* `max_pfx` - maximum number of prefixes allowed from the peer for BGP peer prefix object.
* `name_alias` - name_alias for BGP peer prefix object.
* `restart_time` - time before restarting peer for BGP peer prefix object.
* `thresh` - threshold for a maximum number of prefixes for BGP peer prefix object.
