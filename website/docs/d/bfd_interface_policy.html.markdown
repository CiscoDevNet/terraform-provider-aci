---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bfd_interface_policy"
sidebar_current: "docs-aci-data-source-bfd_interface_policy"
description: |-
  Data source for ACI BFD Interface Policy
---

# aci_bfd_interface_policy #
Data source for ACI BFD Interface Policy

## API Information ##
* `Class` - bfdIfPol
* `Distinguished Name` - uni/tn-{name}/bfdIfPol-{name}

## GUI Information ##
* `Location` - Tenant -> Policies -> Protocol -> BFD

## Example Usage ##

```hcl
data "aci_bfd_interface_policy" "example" {
  tenant_dn  = aci_tenant.tenant_for_bfdIfPol.id
  name  = "example"
}
```

## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent tenant object.
* `name` - (Required) name of object BFD Interface Policy.

## Attribute Reference
* `id` - Attribute id set to the Dn of the BFD Interface Policy.
* `admin_st` - (Optional) Administrative state of the BFD Interface Policy.
* `annotation` - (Optional) Annotation for object BFD Interface Policy.
* `ctrl` - (Optional) Control state of object BFD Interface Policy
* `detect_mult` - (Optional) Detection multiplier for object BFD Interface Policy.
* `echo_admin_st` - (Optional) Echo mode indicator for object BFD Interface Policy.
* `echo_rx_intvl` - (Optional) Echo rx interval for object BFD Interface Policy.
* `min_rx_intvl` - (Optional) Required minimum rx interval for object BFD Interface Policy.
* `min_tx_intvl` - (Optional) Desired minimum tx interval for object BFD Interface Policy.
* `name_alias` - (Optional) Name alias for object BFD Interface Policy.
* `description` - (Optional) Description for object BFD Interface Policy.
