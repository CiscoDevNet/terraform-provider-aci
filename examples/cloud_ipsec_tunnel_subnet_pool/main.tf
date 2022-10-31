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

resource "aci_cloud_ipsec_tunnel_subnet_pool" "ipsec_tunnel_subnet_pool" {
  subnet_pool_name = "test"
	subnet_pool = "160.254.10.0/16"
}

data "aci_cloud_ipsec_tunnel_subnet_pool" "example" {
  subnet_pool  = aci_cloud_ipsec_tunnel_subnet_pool.ipsec_tunnel_subnet_pool.subnet_pool
}

output "ipsec_tunnel_subnet_pool_output" {
  value = data.aci_cloud_ipsec_tunnel_subnet_pool.example
}
