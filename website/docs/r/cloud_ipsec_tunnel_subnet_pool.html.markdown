---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_ipsec_tunnel_subnet_pool
sidebar_current: "docs-aci-resource-aci_cloud_ipsec_tunnel_subnet_pool"
description: |-
  Manages Cloud Network Controller Cloud Subnet Pool for IPsec Tunnels
---

# aci_cloud_ipsec_tunnel_subnet_pool #

Manages Cloud Network Controller Cloud Subnet Pool for IPsec Tunnels
<b>Note: This resource is supported in Cloud Network Controller version > 25.0 only.</b>

## API Information ##

* `Class` - cloudtemplateIpSecTunnelSubnetPool
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_network_name}/ipsecsubnetpool-[{subnetpool}]

## GUI Information ##

* `Location` - Cloud Network Controller -> Infrastructure -> Inter-Site Connectivity -> Region Management


## Example Usage ##

```hcl
resource "aci_cloud_ipsec_tunnel_subnet_pool" "example" {
  name        = "subent_pool_1"
  subnet_pool = "160.254.10.0/16"
}
``` 

## Argument Reference ##

* `name` - (Required) Subnet Pool Name of the Subnet Pool for IPsec Tunnels object.
* `subnet_pool` - (Required) Subnetpool address of the Subnet Pool for IPsec Tunnels object.
* `annotation` - (Optional) Annotation of the Subnet Pool for IPsec Tunnels object.


## Importing ##

An existing Cloud Subnet Pool for IPsec Tunnels can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_ipsec_tunnel_subnet_pool.example "<Dn>"
```