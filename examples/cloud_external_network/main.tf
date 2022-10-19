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
  tenant_dn = data.aci_tenant.infra_tenant.id # Create vrf in infra tenant.
  name      = "cloudVrf"
}

resource "aci_cloud_external_network" "external_network" {
	name = "cloud_external_network"
	vrf_dn = aci_vrf.vrf.id
  all_regions = "yes"
}

data "aci_cloud_external_network" "example" {
  name  = aci_cloud_external_network.external_network.name
}

output "name" {
  value = data.aci_cloud_external_network.example
}

resource "aci_cloud_external_network_vpn_network" "vpn_network" {
  aci_cloud_external_network_dn = aci_cloud_external_network.external_network.id
	name = "cloud_vpn_network"
  cloud_ipsec_tunnel {
    ike_version = "ikev2"
    public_ip_address = "10.10.10.2"
    subnet_pool_name = "azure_pool"
    # peer_asn = "1002"
  }
}