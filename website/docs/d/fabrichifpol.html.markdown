---
layout: "aci"
page_title: "ACI: aci_fabric_if_pol"
sidebar_current: "docs-aci-data-source-fabric_if_pol"
description: |-
  Data source for ACI fabric if pol
---

# aci_fabric_if_pol #
Data source for ACI fabric if pol

## Example Usage ##

```hcl

data "aci_fabric_if_pol" "example" {
  name  = "example"
}

```


## Argument Reference ##
* `name` - (Required) name of Object fabric if pol.



## Attribute Reference

* `id` - Attribute id set to the Dn of the fabric if pol.
* `annotation` - (Optional) annotation for object fabric if pol.
* `auto_neg` - (Optional) policy auto-negotiation
* `fec_mode` - (Optional) forwarding error correction
* `link_debounce` - (Optional) link debounce interval
* `name_alias` - (Optional) name alias for object fabric if pol.
* `speed` - (Optional) port speed
