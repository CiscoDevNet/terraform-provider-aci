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


resource "aci_vrf" "vrf" {
  tenant_dn = "uni/tn-infra"
  name      = "cloudVrf"
}
# Works only with infra tenant and name set to default
# resource "aci_cloud_template_infra_network" "infra_network" {
#   tenant_dn = "uni/tn-infra"
# 	name = "default"
# }

resource "aci_cloud_template_external_network" "external_network" {
  infra_network_template_dn = "uni/tn-infra/infranetwork-default"
	# hub_network_name = "cloud_external_network"
	name = "cloud_external_network"
	vrf_name = aci_vrf.vrf.name
}