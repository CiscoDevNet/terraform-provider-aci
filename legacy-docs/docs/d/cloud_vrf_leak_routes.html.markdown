---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_vrf_leak_routes"
sidebar_current: "docs-aci-data-source-cloud_vrf_leak_routes"
description: |-
  Data source for Cloud Network Controller Inter-VRF Leaked Routes
---

# aci_cloud_vrf_leak_routes #

Data source for Cloud Network Controller Inter-VRF Leaked Routes


## API Information ##

* `Class` - leakInternalPrefix
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintprefix-[{ip}]


## GUI Information ##

* `Location` - Application Management -> VRFs -> Leak Routes


## Example Usage ##

```hcl
data "aci_cloud_vrf_leak_routes" "cloud_internal_leak_routes" {
  vrf_dn    = aci_vrf.vrf1.id
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Inter-VRF Leaked Route object.
* `annotation` - (Optional) Annotation of the Inter-VRF Leaked Route object.
* `name_alias` - (Optional) Name Alias of the Inter-VRF Leaked Route object.
* `leak_to` - (Optional) A block representing the attributes of `Leak Routes` for the Inter-VRF Leaked Route object. Type: Block.
  * `vrf_dn` - Distinguished name of the destination VRF object, which is mapped to the Inter-VRF Leaked Route object.
