---
layout: "aci"
page_title: "ACI: aci_interface_fc_policy"
sidebar_current: "docs-aci-resource-interface_fc_policy"
description: |-
  Manages ACI Interface FC Policy
---

# aci_interface_fc_policy #

Manages ACI Interface FC Policy

## Example Usage ##

```hcl
resource "aci_interface_fc_policy" "example" {
  name         = "demo_policy"
  annotation   = "tag_if_policy"
  automaxspeed = "32G"
  fill_pattern = "default"
  name_alias   = "demo_alias"
  port_mode    = "f"
  rx_bb_credit = "64"
  speed        = "auto"
  trunk_mode   = "trunk-off"
}

```

## Argument Reference ##

* `name` - (Required) name of Object interface_fc_policy.
* `annotation` - (Optional) annotation for object interface_fc_policy.
* `automaxspeed` - (Optional) automaxspeed for object interface_fc_policy.
* `fill_pattern` - (Optional) Fill Pattern for native FC ports. Allowed values are "ARBFF" and "IDLE". Default is "IDLE".
* `name_alias` - (Optional) name_alias for object interface_fc_policy.
* `port_mode` - (Optional) In which mode Ports should be used. Allowed values are "f" and "np". Default is "f".
* `rx_bb_credit` - (Optional) Receive buffer credits for native FC ports Range:(16 - 64). Default value is 64.
* `speed` - (Optional) cpu or port speed. All the supported values are unknown, auto, 4G, 8G, 16G, 32G. Default value is auto.  
* `trunk_mode` - (Optional) Trunking on/off for native FC ports. Allowed values are "un-init", "trunk-off", "trunk-on" and "auto".Default value is "trunk-off".

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Interface FC Policy.

## Importing ##

An existing Interface FC Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_interface_fc_policy.example <Dn>
```
