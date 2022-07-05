---
layout: "aci"
page_title: "ACI: aci_leak_internal_subnet"
sidebar_current: "docs-aci-data-source-leak_internal_subnet"
description: |-
  Data source for ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_leak_internal_subnet #

Data source for ACI Inter-VRF Leaked EPG/BD Subnet


## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintsubnet-[{subnet_ip}]

## GUI Information ##

* `Location` - Tenants -> Networking -> VRFs -> Inter-VRF Leaked Routes for EPG -> EPG/BD Subnets


## Example Usage ##

```hcl
data "aci_leak_internal_subnet" "internal_subnet" {
  vrf_dn    = aci_vrf.vrf1.id # Source VRF DN
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
* `vrf_scope` - (Optional) Visibility of the Inter-VRF Leaked EPG/BD Subnet object.
* `leak_to` - (Optional) A block representing the attributes of `Tenant and VRF destination` for Inter-VRF Leaked Routes object. Type: Block.
  * `destination_vrf_name` - Name of the destination VRF object, which is mapped with `Tenant and VRF Destinations` object.
  * `destination_vrf_scope` - Scope of the `Tenant and VRF Destinations` object.
  * `destination_tenant_name` - Name of the destination Tenant object, which is mapped with `Tenant and VRF Destinations` object.