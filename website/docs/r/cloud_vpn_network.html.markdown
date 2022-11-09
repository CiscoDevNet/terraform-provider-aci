---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network_vpn_network"
sidebar_current: "docs-aci-resource-aci_cloud_external_network_vpn_network"
description: |-
  Manages ACI Template for VPN Network
---

# aci_cloud_external_network_vpn_network #

Manages ACI Template for VPN Network

## API Information ##

* `Class` - cloudtemplateVpnNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}/vpnnetwork-{name}

## GUI Information ##

* `Location` -  Tenants -> {tenant_name} -> Application Management -> External Networks -> VPN Networks


## Example Usage ##

```hcl
resource "aci_cloud_external_network_vpn_network" "example" {
  aci_cloud_external_network_dn  = aci_cloud_external_network.example.id
  name  = "example"
  remote_site_id = "0"
  remote_site_name = "remote_site_1"
  ipsec_tunnel {
    ike_version = "ikev2"
    public_ip_address = "10.10.10.2"
    subnet_pool_name = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool_name
    bgp_peer_asn = "1000"
    source_interfaces = ["gig2", "gig3", "gig4"]
  }
}
```

## Argument Reference ##

* `aci_cloud_external_network_dn` - (Required) Distinguished name of the parent TemplateforExternalNetwork object.
* `name` - (Required) Name of the Cloud VPN Network object.
* `remote_site_id` - (Optional) Remote Site ID. Allowed range is 0-1000 and default value is "0".
* `remote_site_name` - (Optional) Remote Site Name.
* `ipsec_tunnel` - (Optional) IPSec tunnel destination. Type: Block.
    * `ike_version` - (Required) IKE version. Allowed values are "ikev1", "ikev2", and default value is "ikev2".
    * `public_ip_address` - (Required) Peer address of the Cloud IpSec tunnel object.
    * `subnet_pool_name` - (Required) Subnet Pool Name.
    * `pre_shared_key` - (Optional) Pre Shared Key for all tunnels to this peeraddr.
    * `bgp_peer_asn` - (Required) BGP ASN Number. A number that uniquely identifies an autonomous system.
    * `source_interfaces` - (Optional) Source Interface Ids of the object for IpSec tunnel Source Interface. It is available only on Azure cAPIC.


## Importing ##

An existing Cloud VPN Network can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_external_network_vpn_network.example "<Dn>"
```