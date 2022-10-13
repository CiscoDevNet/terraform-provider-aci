---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_leak_epg_bd_subnet"
sidebar_current: "docs-aci-data-source-vrf_leak_epg_bd_subnet"
description: |-
  Data source for ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_vrf_leak_epg_bd_subnet #

Data source for ACI Inter-VRF Leaked EPG/BD Subnet


## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintsubnet-[{subnet_ip}]

## GUI Information ##

* `Location` - Tenants -> Networking -> VRFs -> Inter-VRF Leaked Routes for EPG -> EPG/BD Subnets
* `Location` - Cloud APIC -> Application Management -> VRFs -> Leak Routes

## Example Usage ##

```hcl
data "aci_vrf_leak_epg_bd_subnet" "internal_subnet" {
  vrf_dn    = aci_vrf.vrf1.id
  ip        = "1.1.20.2/24"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `ip` - (Required) IP of the Inter-VRF Leaked EPG/BD Subnet object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Inter-VRF Leaked EPG/BD Subnet object.
* `annotation` - (Optional) Annotation of the Inter-VRF Leaked EPG/BD Subnet object.
* `name_alias` - (Optional) Name Alias of the Inter-VRF Leaked EPG/BD Subnet object.
* `allow_l3out_advertisement` - (Optional) Visibility of the Inter-VRF Leaked EPG/BD Subnet object.
* `leak_to` - (Optional) A block representing the attributes of `Tenant and VRF Destinations` for the Inter-VRF Leaked Routes object. Type: Block.
  * `vrf_dn` - Distinguished name of the destination VRF object, which is mapped to the `Tenant and VRF Destinations` object.
  * `allow_l3out_advertisement` - Scope of the destination VRF object, which is mapped to the `Tenant and VRF Destinations` object.
