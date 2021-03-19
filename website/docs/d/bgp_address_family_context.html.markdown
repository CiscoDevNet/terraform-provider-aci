---
layout: "aci"
page_title: "ACI: aci_bgp_address_family_context"
sidebar_current: "docs-aci-data-source-bgp_address_family_context"
description: |-
  Data source for ACI BGP Address Family Context
---

# aci_bgp_address_family_context #
Data source for ACI BGP Address Family Context

## Example Usage ##

```hcl
data "aci_bgp_address_family_context" "check" {
  tenant_dn = "${aci_tenant.tenentcheck.id}"
  name      = "one"
}
```


## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent tenant object.
* `name` - (Required) name of BGP address family context object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the BGP address family context object.
* `annotation` - (Optional) Annotation for BGP address family context object.
* `ctrl` - (Optional) Control state for BGP address family context object.
* `e_dist` - (Optional) Administrative distance of EBGP routes for BGP address family context object.
* `i_dist` - (Optional) Administrative distance of IBGP routes for BGP address family context object.
* `local_dist` - (Optional) Administrative distance of local routes for BGP address family context object.
* `max_ecmp` - (Optional) Maximum number of equal-cost paths for BGP address family context object.
* `max_ecmp_ibgp` - (Optional) Maximum ECMP IBGP for BGP address family context object.
* `name_alias` - (Optional) Name alias for BGP address family context object.
