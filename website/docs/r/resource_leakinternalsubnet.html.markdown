---
layout: "aci"
page_title: "ACI: aci_inter-vrf_leaked_epg/bd_subnet"
sidebar_current: "docs-aci-resource-inter-vrf_leaked_epg/bd_subnet"
description: |-
  Manages ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_inter-vrf_leaked_epg/bd_subnet #

Manages ACI Inter-VRF Leaked EPG/BD Subnet

## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/leakroutes/leakintsubnet-[{ip}]

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_inter-vrf_leaked_epg/bd_subnet" "example" {
  vrf_dn  = aci_vrf.example.id
  ip  = "example"
  annotation = "orchestrator:terraform"
  ip = 

  scope = "private"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of the parent VRF object.
* `ip` - (Required) Ip of the object Inter-VRF Leaked EPG/BD Subnet.
* `annotation` - (Optional) Annotation of the object Inter-VRF Leaked EPG/BD Subnet.

* `ip` - (Optional) Subnet.The IP address.
* `scope` - (Optional) Visibility of the Subnet.The domain applicable to the capability. Allowed values are "private", "public", and default value is "private". Type: String.


## Importing ##

An existing Inter-VRFLeakedEPG/BDSubnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_inter-vrf_leaked_epg/bd_subnet.example <Dn>
```