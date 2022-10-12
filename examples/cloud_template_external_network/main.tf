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

resource "aci_cloud_template_external_network" "external_network" {
	name = "cloud_external_network"
	vrf_dn = aci_vrf.vrf.id
}

data "aci_cloud_template_external_network" "example" {
  name  = aci_cloud_template_external_network.external_network.name
}

output "name" {
  value = data.aci_cloud_template_external_network.example
}