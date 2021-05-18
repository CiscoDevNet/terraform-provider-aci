---
layout: "aci"
page_title: "ACI: aci_vlan_pool"
sidebar_current: "docs-aci-resource-vlan_pool"
description: |-
  Manages ACI VLAN Pool
---

# aci_vlan_pool #

Manages ACI VLAN Pool

## Example Usage ##

```hcl
resource "aci_vlan_pool" "example" {


  name  = "example"

  alloc_mode  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```

## Argument Reference ##

* `name` - (Required) name of Object vlan_pool.
* `alloc_mode` - (Required) allocation mode.  Allowed values: "dynamic", "static"
* `annotation` - (Optional) annotation for object vlan_pool.
* `name_alias` - (Optional) name_alias for object vlan_pool.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VLAN Pool.

## Importing ##

An existing VLAN Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_vlan_pool.example <Dn>
```
