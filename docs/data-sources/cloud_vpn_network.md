---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network_vpn_network"
sidebar_current: "docs-aci-data-source-aci_cloud_external_network_vpn_network"
description: |-
  Data source for Cloud Network Controller Cloud Template for VPN Network
---

# aci_cloud_external_network_vpn_network #

Data source for Cloud Network Controller Cloud Template for VPN Network
<b>Note: This resource is supported in Cloud Network Controller version > 25.0 only.</b>


## API Information ##

* `Class` - cloudtemplateVpnNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_network_name}/extnetwork-{external_network_name}/vpnnetwork-{vpn_network_name}

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> External Networks -> VPN Networks



## Example Usage ##

```hcl
data "aci_cloud_external_network_vpn_network" "example" {
  aci_cloud_external_network_dn = aci_cloud_external_network.example.id
  name                          = "example"
}
```

## Argument Reference ##

* `aci_cloud_external_network_dn` - (Required) Distinguished name of parent TemplateforExternalNetwork object.
* `name` - (Required) Name of the Cloud VPN Network object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Cloud VPN Network.
* `remote_site_id` - (Optional) Remote Site ID. 
* `remote_site_name` - (Optional) Name of the Remote Site. 
* `ipsec_tunnel` - (Optional) IPsec tunnel destination (cloudtemplateIpSecTunnelSourceInterface class). Type: Block.
    * `ike_version` - (Required) IKE version. Allowed values are "ikev1", "ikev2", and default value is "ikev2".
    * `public_ip_address` - (Required) Peer address of the Cloud IPsec tunnel object.
    * `subnet_pool_name` - (Required) Subnet Pool Name.
    * `pre_shared_key` - (Optional) Pre Shared Key for all tunnels to this peer address.
    * `bgp_peer_asn` - (Required) BGP ASN Number. A number that uniquely identifies an autonomous system.
    * `source_interfaces` - (Optional) Source Interface Ids of the object for IPsec tunnel Source Interface. It is available only on Azure cAPIC.


