---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_ipsec_tunnel_subnet_pool"
sidebar_current: "docs-aci-data-source-aci_cloud_ipsec_tunnel_subnet_pool"
description: |-
  Data source for the ACI Cloud Subnet Pool for IPsec Tunnels
---

# aci_cloud_ipsec_tunnel_subnet_pool #

Data source for the ACI Cloud Subnet Pool for IPsec Tunnels
<b>Note: This resource is supported in Cloud APIC version > 25.0 only.</b>


## API Information ##

* `Class` - cloudtemplateIpSecTunnelSubnetPool
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_network_name}/ipsecsubnetpool-[{subnet_pool}]

## GUI Information ##

* `Location` - Cloud APIC -> Infrastructure -> Inter-Site Connectivity -> Region Management



## Example Usage ##

```hcl
data "aci_cloud_ipsec_tunnel_subnet_pool" "example" {
  name        = "example"
  subnet_pool = "160.254.10.0/16"
}
```

## Argument Reference ##

* `subnet_pool` - (Required) Subnet of the Subnet Pool for IPsec Tunnels object.

## Attribute Reference ##
* `annotation` - (Optional) Annotation of the Subnet Pool for IPsec Tunnels object.
* `name` - (Required) Subnet Pool Name of the Subnet Pool for IPsec Tunnels object.
