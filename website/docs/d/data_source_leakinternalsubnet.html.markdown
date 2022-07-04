---
layout: "aci"
page_title: "ACI: aci_inter-vrf_leaked_epg/bd_subnet"
sidebar_current: "docs-aci-data-source-inter-vrf_leaked_epg/bd_subnet"
description: |-
  Data source for ACI Inter-VRF Leaked EPG/BD Subnet
---

# aci_inter-vrf_leaked_epg/bd_subnet #

Data source for ACI Inter-VRF Leaked EPG/BD Subnet


## API Information ##

* `Class` - leakInternalSubnet
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/leakroutes/leakintsubnet-[{ip}]

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_inter-vrf_leaked_epg/bd_subnet" "example" {
  vrf_dn  = aci_vrf.example.id
  ip  = "example"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name of parent VRF object.
* `ip` - (Required) Ip of object Inter-VRF Leaked EPG/BD Subnet.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Inter-VRF Leaked EPG/BD Subnet.
* `annotation` - (Optional) Annotation of object Inter-VRF Leaked EPG/BD Subnet.
* `name_alias` - (Optional) Name Alias of object Inter-VRF Leaked EPG/BD Subnet.
* `ip` - (Optional) Subnet. The IP address.
* `scope` - (Optional) Visibility of the Subnet. The domain applicable to the capability.
