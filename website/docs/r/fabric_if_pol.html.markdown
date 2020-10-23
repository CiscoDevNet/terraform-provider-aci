---
layout: "aci"
page_title: "ACI: aci_fabric_if_pol"
sidebar_current: "docs-aci-resource-fabric_if_pol"
description: |-
  Manages ACI fabric if pol
---

# aci_fabric_if_pol #
Manages ACI fabric if pol

## Example Usage ##

```hcl

resource "aci_fabric_if_pol" "example" {
  name        = "example"
  description = "hello"
  annotation  = "annotation"
  auto_neg    = "on"
  fec_mode    = "inherit"
  name_alias  = "example"
  speed       = "inherit"
}

```


## Argument Reference ##
* `name` - (Required) name of Object fabric if pol.
* `annotation` - (Optional) annotation for object fabric if pol.
* `auto_neg` - (Optional) policy auto-negotiation. Allowed values: "on", "off"
* `fec_mode` - (Optional) forwarding error correction. Allowed values: "inherit", "cl91-rs-fec", "cl74-fc-fec", "ieee-rs-fec", "cons16-rs-fec", "kp-fec", "disable-fec".
* `link_debounce` - (Optional) link debounce interval
* `name_alias` - (Optional) name_alias for object fabric if pol.
* `speed` - (Optional) port speed. Allowed values: "unknown", "100M", "1G", "10G", "25G", "40G", "50G", "100G","200G", "400G", "inherit".



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Link Level Policy.

## Importing ##

An existing Link Level Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_fabric_if_pol.example <Dn>
```