---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_cloud_vrf_leak_routes"
sidebar_current: "docs-aci-resource-cloud_vrf_leak_routes"
description: |-
  Manages Cloud ACI Inter-VRF Leaked Internal Prefix
---

# aci_cloud_vrf_leak_routes #

Manages Cloud ACI Inter-VRF Leaked Internal Prefix

## API Information ##

* `Class` - leakInternalPrefix
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintprefix-[{ip}]

## GUI Information ##

* `Location` - Application Management -> VRFs -> Leak Routes


## Example Usage ##

```hcl
# Only for the Cloud APIC Version >= 25.0
resource "aci_cloud_vrf_leak_routes" "cloud_internal_leak_routes" {
  vrf_dn = aci_vrf.src_vrf.id
  leak_to {
    vrf_dn = aci_vrf.dst_vrf1.id
  }
  leak_to {
    vrf_dn = aci_vrf.dst_vrf2.id
  }
}
```
Leak Routes for the Inter-VRF Leaked Internal Prefix
## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `annotation` - (Optional) Annotation of the Inter-VRF Leaked Internal Prefix object.
* `name_alias` - (Optional) Name Alias of the Inter-VRF Leaked Internal Prefix object.
* `leak_to` - (Optional) A block representing the attributes of `Leak Routes` for the Inter-VRF Leaked Internal Prefix object. Type: Block.
  * `vrf_dn` - (Required) Distinguished name of the destination VRF object, which is mapped to the Inter-VRF Leaked Internal Prefix object.

## Importing ##

An existing Inter-VRF Leaked Internal Prefix can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_vrf_leak_routes.example <Dn>
```