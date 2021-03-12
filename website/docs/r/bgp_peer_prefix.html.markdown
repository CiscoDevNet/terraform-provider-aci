---
layout: "aci"
page_title: "ACI: aci_bgp_peer_prefix"
sidebar_current: "docs-aci-resource-bgp_peer_prefix"
description: |-
  Manages ACI BGP Peer Prefix
---

# aci_bgp_peer_prefix #
Manages ACI BGP Peer Prefix

## Example Usage ##

```hcl
resource "aci_bgp_peer_prefix" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  action  = "example"
  annotation  = "example"
  max_pfx  = "example"
  name_alias  = "example"
  restart_time  = "example"
  thresh  = "example"
}
```


## Argument Reference ##

* `tenant_dn` - (Required) distinguished name of parent tenant object.
* `name` - (Required) name of BGP peer prefix object.
* `action` - (Optional) action when the maximum prefix limit is reached for BGP peer prefix object.
* `annotation` - (Optional) annotation for BGP peer prefix object.
* `max_pfx` - (Optional) maximum number of prefixes allowed from the peer for BGP peer prefix object.
* `name_alias` - (Optional) name_alias for BGP peer prefix object.
* `restart_time` - (Optional) time before restarting peer for BGP peer prefix object.
* `thresh` - (Optional) threshold for a maximum number of prefixes for BGP peer prefix object.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Peer Prefix.

## Importing ##

An existing BGP Peer Prefix can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bgp_peer_prefix.example <Dn>
```