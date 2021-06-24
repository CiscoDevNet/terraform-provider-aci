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
  tenant_dn    = aci_tenant.tenentcheck.id
  name         = "one"
  description  = "from terraform"
  action       = "shut"
  annotation   = "example"
  max_pfx      = "200"
  name_alias   = "example"
  restart_time = "200"
  thresh       = "85"
}
```


## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of BGP peer prefix object.
* `action` - (Optional) Action when the maximum prefix limit is reached for BGP peer prefix object. Allowed values are "log", "reject", "restart" and "shut". Default value is "reject".
* `description` - (Optional) Description for BGP peer prefix object.
* `annotation` - (Optional) Annotation for BGP peer prefix object.
* `max_pfx` - (Optional) Maximum number of prefixes allowed from the peer for BGP peer prefix object. Default value is "20000".
* `name_alias` - (Optional) Name alias for BGP peer prefix object.
* `restart_time` - (Optional) The period of time in minutes before restarting the peer when the prefix limit is reached for BGP peer prefix object. Default value is "infinite".
* `thresh` - (Optional) Threshold percentage of the maximum number of prefixes before a warning is issued for BGP peer prefix object. Default value is "75".
 


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Peer Prefix.

## Importing ##

An existing BGP Peer Prefix can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bgp_peer_prefix.example <Dn>
```