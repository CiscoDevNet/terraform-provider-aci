---
layout: "aci"
page_title: "ACI: aci_leak_internal_subnet"
sidebar_current: "docs-aci-resource-leak_internal_subnet"
description: |-
  Manages ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_leak_internal_subnet #

Manages ACI Inter-VRF Leaked EPG/BD Subnet

## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{tenant_name}/ctx-{vrf_name}/leakroutes/leakintsubnet-[{subnet_ip}]

## GUI Information ##

* `Location` - Tenants -> Networking -> VRFs -> Inter-VRF Leaked Routes for EPG -> EPG/BD Subnets


## Example Usage ##

```hcl
resource "aci_leak_internal_subnet" "internal_subnet" {
  vrf_dn    = aci_vrf.vrf1.id # Source VRF DN
  ip        = "1.1.20.2/24"
  vrf_scope = "public"
  leak_to {
    destination_vrf_name    = "default"
    destination_tenant_name = "common"
  }
  leak_to {
    destination_vrf_name    = aci_vrf.vrf2.name
    destination_tenant_name = aci_tenant.terraform_vrf.name
    destination_vrf_scope   = "private"
  }
}
```
Tenant and VRF destination for Inter-VRF Leaked Routes
## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `ip` - (Required) IP of the Inter-VRF Leaked EPG/BD Subnet object.
* `annotation` - (Optional) Annotation of the object Inter-VRF Leaked EPG/BD Subnet.
* `vrf_scope` - (Optional) Visibility of the Inter-VRF Leaked EPG/BD Subnet object. Allowed values are "private", "public", and default value is "private".
* `leak_to` - (Optional) A block representing the attributes of `Tenant and VRF destination` for Inter-VRF Leaked Routes object. Type: Block.
  * `destination_vrf_name` - Name of the destination VRF object, which is mapped with `Tenant and VRF Destinations` object.
  * `destination_vrf_scope` - Scope of the `Tenant and VRF Destinations` object. Allowed values are "inherit", "private", "public", and default value is "inherit".
  * `destination_tenant_name` - Name of the destination Tenant object, which is mapped with `Tenant and VRF Destinations` object.

## Importing ##

An existing Inter-VRFLeakedEPG/BDSubnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leak_internal_subnet.example <Dn>
```