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
  description = "From Terraform"
  alloc_mode  = "static"
  annotation  = "example"
  name_alias  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of Object vlan pool.
* `alloc_mode` - (Required) Allocation mode for object vlan_pool. Allowed values: "dynamic", "static"
* `description` - (Optional) Description for  object vlan pool.
* `annotation` - (Optional) Annotation for object vlan pool.
* `name_alias` - (Optional) Name alias for  object vlan pool.

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VLAN Pool.

## Importing ##

An existing VLAN Pool can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_vlan_pool.example <Dn>
```
