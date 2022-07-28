---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_leak_epg_bd_subnet"
sidebar_current: "docs-aci-resource-vrf_leak_epg_bd_subnet"
description: |-
  Manages ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_vrf_leak_epg_bd_subnet #

Manages ACI Inter-VRF Leaked EPG/BD Subnet

## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintsubnet-[{subnet_ip}]

## GUI Information ##

* `Location` - Tenants -> Networking -> VRFs -> Inter-VRF Leaked Routes for EPG -> EPG/BD Subnets


## Example Usage ##

```hcl
resource "aci_vrf_leak_epg_bd_subnet" "vrf_leak_epg_bd_subnet" {
  vrf_dn                    = aci_vrf.vrf1.id
  ip                        = "1.1.20.2/24"
  allow_l3out_advertisement = true # true -> public, false -> private, default -> false(private)
  leak_to {
    vrf_dn                    = data.aci_vrf.default.id
  }
  leak_to {
    vrf_dn                    = aci_vrf.vrf2.id
    allow_l3out_advertisement = true # true -> public, false -> private, default -> "inherit"
  }
}
```
Tenant and VRF destination for Inter-VRF Leaked Routes
## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `ip` - (Required) IP of the Inter-VRF Leaked EPG/BD Subnet object.
* `annotation` - (Optional) Annotation of the object Inter-VRF Leaked EPG/BD Subnet.
* `allow_l3out_advertisement` - (Optional) Visibility of the Inter-VRF Leaked EPG/BD Subnet object. Allowed values are "true", "false" and default value is "false". Type: String.
* `leak_to` - (Optional) A block representing the attributes of `Tenant and VRF Destinations` for Inter-VRF Leaked Routes object. Type: Block.
  * `vrf_dn` - (Required) Distinguished name of the destination VRF object, which is mapped with `Tenant and VRF Destinations` object.
  * `allow_l3out_advertisement` - (Optional) Scope of the destination VRF object, which is mapped with `Tenant and VRF Destinations` object. Allowed values are "inherit", "true", "false" and default value is "inherit". Type: String.

## Importing ##

An existing Inter-VRFLeakedEPG/BDSubnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vrf_leak_epg_bd_subnet.example <Dn>
```