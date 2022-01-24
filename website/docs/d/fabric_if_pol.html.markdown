---
subcategory: "Access Policies"
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
* `name` - (Required) Name of object fabric if pol.

## Attribute Reference

* `id` - Attribute id set to the Dn of the fabric if pol.
* `annotation` - (Optional) Annotation for object fabric if pol.
* `description` - (Optional) Description for object fabric if pol.
* `auto_neg` - (Optional) Policy auto-negotiation for object fabric if pol.
* `fec_mode` - (Optional) Forwarding error correction for object fabric if pol.
* `link_debounce` - (Optional) Link debounce interval for object fabric if pol.
* `name_alias` - (Optional) Name alias for object fabric if pol.
* `speed` - (Optional) Port speed for object fabric if pol.
