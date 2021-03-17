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

* `tenant_dn` - (Required) Distinguished name of parent tenant object.
* `name` - (Required) Name of BGP peer prefix object.



## Attribute Reference

* `id` - Attribute id set to the Dn of BGP peer prefix object.
* `action` - Action when the maximum prefix limit is reached for BGP peer prefix object.
* `description` - Description for BGP peer prefix object.
* `annotation` - Annotation for BGP peer prefix object.
* `max_pfx` - Maximum number of prefixes allowed from the peer for BGP peer prefix object.
* `name_alias` - Name alias for BGP peer prefix object.
* `restart_time` - Time before restarting peer for BGP peer prefix object.
* `thresh` - Threshold for a maximum number of prefixes for BGP peer prefix object.
