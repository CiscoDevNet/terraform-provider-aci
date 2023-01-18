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
  name = "cloud_pool"
  subnet_pool      = "169.254.0.0/16"
}

# GCP cloud
# all_regions is set to "no" and regions can be set only in GCP Cloud
resource "aci_cloud_external_network" "external_network_gcp" {
  name         = "cloud_external_network"
  vrf_dn       = aci_vrf.vrf.id
  cloud_vendor = "gcp"
  all_regions  = "yes"
  regions      = ["europe-west4", "europe-west3"]
}

# GCP cloud
resource "aci_cloud_external_network_vpn_network" "vpn_network" {
  aci_cloud_external_network_dn = aci_cloud_external_network.external_network_gcp.id
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
  name = aci_cloud_external_network.external_network_gcp.name
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