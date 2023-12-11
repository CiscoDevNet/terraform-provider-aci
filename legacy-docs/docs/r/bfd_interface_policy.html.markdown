---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bfd_interface_policy"
sidebar_current: "docs-aci-resource-aci_bfd_interface_policy"
description: |-
  Manages ACI BFD Interface Policy
---

# aci_bfd_interface_policy #
Manages ACI BFD Interface Policy

## API Information ##

* `Class` - bfdIfPol
* `Distinguished Name` - uni/tn-{name}/bfdIfPol-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> BFD

## Example Usage ##

```hcl
resource "aci_bfd_interface_policy" "example" {
  tenant_dn = aci_tenant.tenant_for_bdfIfPol.id
  name = "example"
  admin_st = "enabled"
  annotation  = "example"
  ctrl = "opt-subif"
  detect_mult  = "3"
  echo_admin_st = "disabled"
  echo_rx_intvl  = "50"
  min_rx_intvl  = "50"
  min_tx_intvl  = "50"
  name_alias  = "example"
  description = "example"
}
```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent tenant object.
* `name` - (Required) Name of object BFD Interface Policy.
* `admin_st` - (Optional) Administrative state of the object BFD Interface Policy. Allowed values are "disabled" and "enabled". Default is "enabled".
* `annotation` - (Optional) Annotation for object BFD Interface Policy.
* `ctrl` - (Optional) Control state for object BFD Interface Policy. Allowed values are "opt-subif" and "none". Default is "none".
* `detect_mult` - (Optional) Detection multiplier for object BFD Interface Policy. Range: "1" - "50". Default value is "3".
* `echo_admin_st` - (Optional) Echo mode indicator for object BFD Interface Policy. Allowed values are "disabled" and "enabled". Default is "enabled".  
* `echo_rx_intvl` - (Optional) Echo rx interval for object BFD Interface Policy. Range: "50" - "999". Default value is "50".
* `min_rx_intvl` - (Optional) Required minimum rx interval for boject BFD Interface Policy. Range: "50" - "999". Default value is "50".
* `min_tx_intvl` - (Optional) Desired minimum tx interval for object BFD Interface Policy. Range: "50" - "999". Default value is "50".
* `name_alias` - (Optional) Name alias for object BFD Interface Policy. 
* `description` - (Optional) Description for object BFD Interface Policy.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BFD Interface Policy.

## Importing ##

An existing BFD Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bfd_interface_policy.example <Dn>
```