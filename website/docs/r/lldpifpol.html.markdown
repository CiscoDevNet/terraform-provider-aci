---
layout: "aci"
page_title: "ACI: aci_lldp_interface_policy"
sidebar_current: "docs-aci-resource-lldp_interface_policy"
description: |-
  Manages ACI LLDP Interface Policy
---

# aci_lldp_interface_policy #
Manages ACI LLDP Interface Policy

## Example Usage ##

```hcl
	resource "aci_lldp_interface_policy" "foolldp_interface_policy" {
		description = "%s"
		name        = "demo_lldp_pol"
		admin_rx_st = "%s"
		admin_tx_st = "enabled"
		annotation  = "tag_lldp"
		name_alias  = "alias_lldp"
	} 
```
## Argument Reference ##
* `name` - (Required) name of Object lldp_interface_policy.
* `admin_rx_st` - (Optional) admin receive state. Allowed values are "enabled" and "disabled". Default value is "enabled".
* `admin_tx_st` - (Optional) admin transmit state. Allowed values are "enabled" and "disabled". Default value is "enabled".
* `annotation` - (Optional) annotation for object lldp_interface_policy.
* `name_alias` - (Optional) name_alias for object lldp_interface_policy.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the LLDP Interface Policy.

## Importing ##

An existing LLDP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_lldp_interface_policy.example <Dn>
```