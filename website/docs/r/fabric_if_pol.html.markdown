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
  name        = "fabric_if_pol_1"
  description = "Link Level description"
  annotation  = "fabric_if_pol_tag"
  auto_neg    = "on"
  fec_mode    = "inherit"
  name_alias  = "alias_ifpol"
  link_debounce  = "100"
  speed       = "inherit"
}
```

## Argument Reference ##

* `name` - (Required) Name of object fabric if pol.
* `annotation` - (Optional) Annotation for object fabric if pol.
* `description` - (Optional) Description for object fabric if pol.
* `auto_neg` - (Optional) Policy auto negotiation for object fabric if pol. Allowed values: "on", "off". Default value is "on".
* `fec_mode` - (Optional) Forwarding error correction for object fabric if pol. Allowed values: "inherit", "cl91-rs-fec", "cl74-fc-fec", "ieee-rs-fec", "cons16-rs-fec", "kp-fec", "disable-fec". Default value is "inherit".
* `link_debounce` - (Optional) Link debounce interval for object fabric if pol. Range of allowed values: "0" to "5000". Default value is "100".
* `name_alias` - (Optional) Name alias for object fabric if pol.
* `speed` - (Optional) Port speed for object fabric if pol. Allowed values: "unknown", "100M", "1G", "10G", "25G", "40G", "50G", "100G","200G", "400G", "inherit". Default value is "inherit".

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the Dn of the Link Level Policy.

## Importing ##

An existing Link Level Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_fabric_if_pol.example <Dn>
```
