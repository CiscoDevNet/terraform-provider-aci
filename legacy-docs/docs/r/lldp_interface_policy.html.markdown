---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_lldp_interface_policy"
sidebar_current: "docs-aci-resource-aci_lldp_interface_policy"
description: |-
  Manages ACI LLDP Interface Policy
---

# aci_lldp_interface_policy #

Manages ACI LLDP Interface Policy

## Example Usage ##

```hcl
resource "aci_lldp_interface_policy" "example" {
  description = "example description"
  name        = "demo_lldp_pol"
  admin_rx_st = "enabled"
  admin_tx_st = "enabled"
  annotation  = "tag_lldp"
  name_alias  = "alias_lldp"
} 
```

## Argument Reference ##

* `name` - (Required) Name of Object LLDP Interface Policy.
* `admin_rx_st` - (Optional) Admin receive state. Allowed values are "enabled" and "disabled". Default value is "enabled".
* `admin_tx_st` - (Optional) Admin transmit state. Allowed values are "enabled" and "disabled". Default value is "enabled".
* `description` - (Optional) Description for object LLDP Interface Policy.
* `annotation` - (Optional) Annotation for object LLDP Interface Policy.
* `name_alias` - (Optional) Name alias for object LLDP Interface Policy.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the LLDP Interface Policy.

## Importing ##

An existing LLDP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_lldp_interface_policy.example <Dn>
```
