terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

data "aci_tenant" "infra_tenant" {
  name = "infra"
}

resource "aci_vrf" "vrf" {
  tenant_dn = data.aci_tenant.infra_tenant.id # Create vrf only in infra tenant.
  name      = "cloudVrf"
}

resource "aci_cloud_ipsec_tunnel_subnet_pool" "ipsec_tunnel_subnet_pool" {
  name        = "cloud_pool"
  subnet_pool = "169.254.0.0/16"
}

# AWS Cloud
# all_regions is set to "yes" only in AWS Cloud
# 
# router_type = "tgw"
resource "aci_cloud_external_network" "external_network_aws_tgw" {
  name             = "cloud_external_network"
  vrf_dn           = aci_vrf.vrf.id
  all_regions      = "yes"
  cloud_vendor     = "aws"
  router_type      = "tgw"
  hub_network_name = "Hub" #hub_network_name should be set when router_type is "tgw" in AWS cloud
}

resource "aci_cloud_external_network_vpn_network" "vpn_network" {
  aci_cloud_external_network_dn = aci_cloud_external_network.external_network_aws_tgw.id
  name                          = "cloud_vpn_network"
  ipsec_tunnel {
    ike_version       = "ikev1"
    public_ip_address = "10.10.10.2"
    subnet_pool_name  = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool_name
    bgp_peer_asn      = "1002"
  }
  ipsec_tunnel {
    ike_version       = "ikev2"
    public_ip_address = "10.10.10.7"
    subnet_pool_name  = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool_name
    bgp_peer_asn      = "1005"
  }
}

data "aci_cloud_external_network" "external_network_example" {
  name = aci_cloud_external_network.external_network_aws_tgw.name
}

output "external_network_output" {
  value = data.aci_cloud_external_network.external_network_example
}

data "aci_cloud_external_network_vpn_network" "example" {
  aci_cloud_external_network_dn = aci_cloud_external_network_vpn_network.vpn_network.aci_cloud_external_network_dn
  name                          = aci_cloud_external_network_vpn_network.vpn_network.name
}

output "vpn_network_output" {
  value = data.aci_cloud_external_network_vpn_network.example
}

# # AWS Cloud
# # all_regions is set to "yes" only in AWS Cloud
# # 
# # # router_type = "c8kv"
# resource "aci_cloud_external_network" "external_network_aws_c8kv" {
#   name         = "cloud_external_network"
#   vrf_dn       = aci_vrf.vrf.id
#   all_regions  = "yes"
#   cloud_vendor = "aws"
#   router_type  = "c8kv"
# }

# resource "aci_cloud_external_network_vpn_network" "vpn_network_c8kv" {
#   aci_cloud_external_network_dn = aci_cloud_external_network.external_network_aws_c8kv.id
#   name                          = "cloud_vpn_network"
#   ipsec_tunnel {
#     ike_version       = "ikev1"
#     public_ip_address = "10.10.10.2"
#     subnet_pool_name  = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool_name
#     bgp_peer_asn      = "1002"
#     source_interfaces = ["gig2", "gig3", "gig4"] #source_interfaces available when router_type is "c8kv" in AWS cloud
#   }
#   ipsec_tunnel {
#     ike_version       = "ikev2"
#     public_ip_address = "10.10.10.7"
#     subnet_pool_name  = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool_name
#     bgp_peer_asn      = "1005"
#     source_interfaces = ["gig2"] #source_interfaces available when router_type is "c8kv" in AWS cloud
#   }
# }