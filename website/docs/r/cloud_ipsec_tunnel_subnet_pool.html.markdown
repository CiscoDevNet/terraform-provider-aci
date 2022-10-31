---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_ipsec_tunnel_subnet_pool
sidebar_current: "docs-aci-resource-aci_cloud_ipsec_tunnel_subnet_pool"
description: |-
  Manages ACI Cloud Subnet Pool for IpSec Tunnels
---

# aci_cloud_ipsec_tunnel_subnet_pool #

Manages ACI Cloud Subnet Pool for IpSec Tunnels

## API Information ##

* `Class` - cloudtemplateIpSecTunnelSubnetPool
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/ipsecsubnetpool-[{subnetpool}]

## GUI Information ##

* `Location` - Cloud APIC -> Infrastructure -> Inter-Site Connectivity -> Region Management


## Example Usage ##

```hcl
resource "aci_cloud_ipsec_tunnel_subnet_pool" "example" {
  subnet_pool  = "160.254.10.0/16"
}
```

## Argument Reference ##

* `subnet_pool` - (Required) Subnetpool address of the Subnet Pool for IpSec Tunnels object.
* `annotation` - (Optional) Annotation of the Subnet Pool for IpSec Tunnels object.
* `subnet_pool_name` - (Required) Subnet Pool Name of the Subnet Pool for IpSec Tunnels object.


## Importing ##

An existing Cloud Subnet Pool for IpSec Tunnels can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_ipsec_tunnel_subnet_pool.example "<Dn>"
```